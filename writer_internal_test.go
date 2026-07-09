package epub

import (
	"testing"

	"github.com/raitucarp/epub/ncx"
	"github.com/raitucarp/epub/pkg"
)

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
