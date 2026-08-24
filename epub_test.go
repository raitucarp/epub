package epub

import (
	"testing"

	"github.com/raitucarp/epub/ocf"
	"github.com/raitucarp/epub/pkg"
)

func TestEpub_SelectPackage(t *testing.T) {
	pkgPub1 := &pkg.Package{}
	pkgPub2 := &pkg.Package{}

	epub := &Epub{
		packagePubs: map[string]*pkg.Package{
			"pub1": pkgPub1,
			"pub2": pkgPub2,
		},
		rendition: "pub2",
	}

	t.Run("SelectPackage existing", func(t *testing.T) {
		selected := epub.SelectPackage("pub1")
		if selected != pkgPub1 {
			t.Errorf("expected package pointer to match, got %p", selected)
		}
	})

	t.Run("SelectPackage non-existing", func(t *testing.T) {
		selected := epub.SelectPackage("non-existing")
		if selected != nil {
			t.Errorf("expected nil, got %v", selected)
		}
	})

	t.Run("SelectedPackage based on rendition", func(t *testing.T) {
		selected := epub.SelectedPackage()
		if selected != pkgPub2 {
			t.Errorf("expected package pointer to match, got %p", selected)
		}
	})
}

func TestRenditionKey(t *testing.T) {
	tests := []struct {
		name     string
		rootFile ocf.RootFile
		index    int
		expected string
	}{
		{"first is always default", ocf.RootFile{Label: "main"}, 0, "default"},
		{"second with no attributes", ocf.RootFile{}, 1, "rendition-1"},
		{"second with layout", ocf.RootFile{Layout: "pre-paginated"}, 1, "pre-paginated"},
		{"second with multiple attributes", ocf.RootFile{Layout: "pre-paginated", Language: "en"}, 2, "pre-paginated_en"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renditionKey(tt.rootFile, tt.index); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestReader_ListRenditions(t *testing.T) {
	r := &Reader{
		epub: &Epub{
			packagePubs: map[string]*pkg.Package{
				"pre-paginated": {},
				"default":       {},
				"rendition-1":   {},
			},
		},
	}

	got := r.ListRenditions()
	if len(got) != 3 {
		t.Fatalf("expected 3 renditions, got %d", len(got))
	}
	if got[0] != "default" {
		t.Errorf("expected default rendition first, got %v", got)
	}
}
