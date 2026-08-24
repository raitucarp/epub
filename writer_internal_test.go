package epub

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/pkg"
)

func TestWriter_Modified(t *testing.T) {
	w := New("test-pub-id")
	w.Modified(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC))

	metas := w.epub.SelectedPackage().Metadata.Meta
	if len(metas) != 1 {
		t.Fatalf("expected 1 meta element, got %d", len(metas))
	}

	if metas[0].Property != "dcterms:modified" {
		t.Errorf("expected property dcterms:modified, got %q", metas[0].Property)
	}
	if metas[0].Value != "2024-01-02T03:04:05Z" {
		t.Errorf("expected value 2024-01-02T03:04:05Z, got %q", metas[0].Value)
	}
}

func TestWriter_EnsureModifiedDate(t *testing.T) {
	t.Run("adds when missing", func(t *testing.T) {
		w := New("test-pub-id")
		w.ensureModifiedDate()

		var found bool
		for _, meta := range w.epub.SelectedPackage().Metadata.Meta {
			if meta.Property == "dcterms:modified" && meta.Refines == "" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected dcterms:modified to be added")
		}
	})

	t.Run("does not duplicate", func(t *testing.T) {
		w := New("test-pub-id")
		w.Modified(time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC))
		w.ensureModifiedDate()

		var count int
		for _, meta := range w.epub.SelectedPackage().Metadata.Meta {
			if meta.Property == "dcterms:modified" && meta.Refines == "" {
				count++
			}
		}
		if count != 1 {
			t.Errorf("expected exactly one dcterms:modified, got %d", count)
		}
	})
}

func TestWriter_GuardCheck(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*pkg.Package, *Epub)
		expectedErr string
	}{
		{
			name: "Missing Identifiers",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = nil
			},
			expectedErr: "Package should have identifiers.",
		},
		{
			name: "Missing Titles",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = nil
			},
			expectedErr: "Package should have titles.",
		},
		{
			name: "Missing Languages",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = []pkg.DCTitle{{Value: "test-title"}}
				p.Metadata.Languages = nil
			},
			expectedErr: "Package should have languages.",
		},
		{
			name: "Missing Manifest Items",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = []pkg.DCTitle{{Value: "test-title"}}
				p.Metadata.Languages = []pkg.DCLanguage{{Value: "en"}}
				p.Manifest.Items = nil
			},
			expectedErr: "No content insides.",
		},
		{
			name: "Missing Text Content",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = []pkg.DCTitle{{Value: "test-title"}}
				p.Metadata.Languages = []pkg.DCLanguage{{Value: "en"}}
				p.Manifest.Items = []pkg.Item{
					{MediaType: "image/jpeg", Properties: pkg.CoverImageProperty},
				}
			},
			expectedErr: "No text content insides.",
		},
		{
			name: "Missing Cover Image",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = []pkg.DCTitle{{Value: "test-title"}}
				p.Metadata.Languages = []pkg.DCLanguage{{Value: "en"}}
				p.Manifest.Items = []pkg.Item{
					{MediaType: pkg.MediaTypeXHTML},
				}
			},
			expectedErr: "No cover images.",
		},
		{
			name: "Missing Table of Contents",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = []pkg.DCTitle{{Value: "test-title"}}
				p.Metadata.Languages = []pkg.DCLanguage{{Value: "en"}}
				p.Manifest.Items = []pkg.Item{
					{MediaType: pkg.MediaTypeXHTML},
					{MediaType: "image/jpeg", Properties: pkg.CoverImageProperty},
				}
				e.navigationCenterEXtended = nil
			},
			expectedErr: "No table of contents.",
		},
		{
			name: "Valid Package",
			setup: func(p *pkg.Package, e *Epub) {
				p.Metadata.Identifiers = []pkg.DCIdentifier{{Value: "test-id"}}
				p.Metadata.Titles = []pkg.DCTitle{{Value: "test-title"}}
				p.Metadata.Languages = []pkg.DCLanguage{{Value: "en"}}
				p.Manifest.Items = []pkg.Item{
					{MediaType: pkg.MediaTypeXHTML},
					{MediaType: "image/jpeg", Properties: pkg.CoverImageProperty},
				}
				e.navigationCenterEXtended = &ncx.NCX{}
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{
				epub: &Epub{
					packagePubs: map[string]*pkg.Package{
						"content": {
							Metadata: pkg.Metadata{},
							Manifest: pkg.Manifest{},
						},
					},
				},
			}
			tt.setup(w.epub.packagePubs["content"], w.epub)
			err := w.guardCheck()
			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("Expected error '%s', got nil", tt.expectedErr)
				} else if err.Error() != tt.expectedErr {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedErr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got '%s'", err.Error())
				}
			}
		})
	}
}

func TestWriter_PathTraversal(t *testing.T) {
	w := New("test-pub-id")

	tests := []struct {
		name string
		path string
	}{
		{"Absolute Path", "/etc/passwd"},
		{"Parent Traversal", "../../../etc/passwd"},
		{"Hidden Traversal", "foo/../../etc/passwd"},
		{"Backslash Traversal", "..\\..\\etc\\passwd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := w.AddContentFile(tt.path)
			if err == nil {
				t.Errorf("AddContentFile(%q) expected error, got nil", tt.path)
			}

			err = w.CoverFile(tt.path)
			if err == nil {
				t.Errorf("CoverFile(%q) expected error, got nil", tt.path)
			}

			res := w.AddImageFile(tt.path)
			if res.ID != "" {
				t.Errorf("AddImageFile(%q) expected empty resource on invalid path, got %v", tt.path, res)
			}
		})
	}
}

func TestWriter_SymlinkTraversal(t *testing.T) {
	w := New("test-pub-id")

	// Create a safe directory
	safeDir := t.TempDir()

	// Create a sensitive file outside the safe dir
	sensitiveFile := filepath.Join(t.TempDir(), "secret.txt")
	err := os.WriteFile(sensitiveFile, []byte("super secret"), 0644)
	if err != nil {
		t.Fatalf("Failed to create sensitive file: %v", err)
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Change working directory to safeDir
	err = os.Chdir(safeDir)
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}
	defer os.Chdir(oldWd) // Restore original working directory

	// Create a symlink pointing to the sensitive file
	err = os.Symlink(sensitiveFile, "link_to_secret")
	if err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	// Test AddContentFile
	_, err = w.AddContentFile("link_to_secret")
	if err == nil {
		t.Errorf("AddContentFile: expected error when reading via escaping symlink, got nil")
	}

	// Test CoverFile
	err = w.CoverFile("link_to_secret")
	if err == nil {
		t.Errorf("CoverFile: expected error when reading via escaping symlink, got nil")
	}

	// Test AddImageFile
	res := w.AddImageFile("link_to_secret")
	if res.ID != "" || len(res.Content) > 0 {
		t.Errorf("AddImageFile: expected empty resource when reading via escaping symlink, got %v", res)
	}
}
