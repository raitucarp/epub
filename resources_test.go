package epub

import (
	"reflect"
	"testing"
)

func TestSelectResourceById(t *testing.T) {
	r := &Reader{
		epub: &Epub{
			resources: []PublicationResource{
				{ID: "res1", Href: "test1.xhtml"},
				{ID: "res2", Href: "test2.xhtml"},
			},
		},
	}

	t.Run("existing resource", func(t *testing.T) {
		res := r.SelectResourceById("res1")
		if res == nil {
			t.Fatalf("expected resource, got nil")
		}
		if res.ID != "res1" {
			t.Errorf("expected ID 'res1', got %q", res.ID)
		}
	})

	t.Run("non-existing resource", func(t *testing.T) {
		res := r.SelectResourceById("nonexistent")
		if res != nil {
			t.Fatalf("expected nil, got %v", res)
		}
	})
}

func TestSelectResourceByHref(t *testing.T) {
	r := &Reader{
		epub: &Epub{
			resources: []PublicationResource{
				{ID: "res1", Href: "test1.xhtml"},
				{ID: "res2", Href: "test2.xhtml"},
			},
		},
	}

	t.Run("existing resource", func(t *testing.T) {
		res := r.SelectResourceByHref("test2.xhtml")
		if res == nil {
			t.Fatalf("expected resource, got nil")
		}
		if res.Href != "test2.xhtml" {
			t.Errorf("expected Href 'test2.xhtml', got %q", res.Href)
		}
	})

	t.Run("non-existing resource", func(t *testing.T) {
		res := r.SelectResourceByHref("nonexistent.xhtml")
		if res != nil {
			t.Fatalf("expected nil, got %v", res)
		}
	})
}

func TestResources(t *testing.T) {
	resources := []PublicationResource{
		{ID: "res1", Href: "test1.xhtml"},
		{ID: "res2", Href: "test2.xhtml"},
	}
	r := &Reader{
		epub: &Epub{
			resources: resources,
		},
	}

	res := r.Resources()
	if len(res) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(res))
	}

	if !reflect.DeepEqual(res, resources) {
		t.Errorf("expected resources %v, got %v", resources, res)
	}
}
