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

	t.Run("existing package", func(t *testing.T) {
		if got := epub.SelectPackage("pub1"); got != pkgPub1 {
			t.Errorf("expected package pointer to match, got %p", got)
		}
	})

	t.Run("non-existing package", func(t *testing.T) {
		if got := epub.SelectPackage("non-existing"); got != nil {
			t.Errorf("expected nil, got %v", got)
		}
	})
}

func TestEpub_SelectedPackage(t *testing.T) {
	pkgPub := &pkg.Package{}
	epub := &Epub{
		packagePubs: map[string]*pkg.Package{"content": pkgPub},
		rendition:   "content",
	}

	if got := epub.SelectedPackage(); got != pkgPub {
		t.Errorf("expected selected package pointer, got %p", got)
	}
}

func TestEpub_DefaultPackage(t *testing.T) {
	pkgPub := &pkg.Package{}
	epub := &Epub{
		packagePubs: map[string]*pkg.Package{"content": pkgPub},
	}

	if got := epub.DefaultPackage(); got != pkgPub {
		t.Errorf("expected default package pointer, got %p", got)
	}
}
