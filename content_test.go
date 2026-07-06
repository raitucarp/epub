package epub

import (
	"testing"
	"golang.org/x/net/html"
	"github.com/raitucarp/epub/pkg"
)

func TestReadContentHTMLById(t *testing.T) {
	r := &Reader{
		epub: &Epub{
			resources: []PublicationResource{
				{
					ID:       "valid-xhtml",
					Href:     "valid.xhtml",
					MIMEType: pkg.MediaTypeXHTML,
					Content:  []byte("<html><body><p>Valid Content</p></body></html>"),
				},
				{
					ID:       "image-res",
					Href:     "image.jpg",
					MIMEType: pkg.MediaTypeJPEG,
					Content:  []byte("fake-jpeg-content"),
				},
			},
		},
	}

	t.Run("existing xhtml resource", func(t *testing.T) {
		doc := r.ReadContentHTMLById("valid-xhtml")
		if doc == nil {
			t.Fatalf("expected parsed html.Node, got nil")
		}
		if doc.Type != html.DocumentNode {
			t.Errorf("expected DocumentNode, got %v", doc.Type)
		}
	})

	t.Run("non-existent id", func(t *testing.T) {
		doc := r.ReadContentHTMLById("non-existent")
		if doc != nil {
			t.Fatalf("expected nil, got %v", doc)
		}
	})

	t.Run("mismatched mime type", func(t *testing.T) {
		doc := r.ReadContentHTMLById("image-res")
		if doc != nil {
			t.Fatalf("expected nil for non-XHTML mime type, got %v", doc)
		}
	})
}
