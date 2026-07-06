package epub

import (
	"testing"

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
