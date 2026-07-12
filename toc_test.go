package epub

import (
	"testing"
)

func TestTOC_JSON(t *testing.T) {
	tests := []struct {
		name     string
		toc      *TOC
		expected string
	}{
		{
			name: "nested items",
			toc: &TOC{
				Title: "Chapter 1",
				Href:  "chapter1.html",
				Items: []TOC{
					{
						Title: "Section 1.1",
						Href:  "chapter1.html#section1",
					},
					{
						Title: "Section 1.2",
						Href:  "chapter1.html#section2",
					},
				},
			},
			expected: `{"title":"Chapter 1","href":"chapter1.html","items":[{"title":"Section 1.1","href":"chapter1.html#section1"},{"title":"Section 1.2","href":"chapter1.html#section2"}]}`,
		},
		{
			name:     "empty TOC",
			toc:      &TOC{},
			expected: `{}`,
		},
		{
			name: "no items",
			toc: &TOC{
				Title: "Chapter 2",
				Href:  "chapter2.html",
			},
			expected: `{"title":"Chapter 2","href":"chapter2.html"}`,
		},
		{
			name: "special characters",
			toc: &TOC{
				Title: "Chapter 3: \"The Awakening\" & <Others>",
				Href:  "chapter3.html?foo=bar&baz=qux",
			},
			expected: `{"title":"Chapter 3: \"The Awakening\" \u0026 \u003cOthers\u003e","href":"chapter3.html?foo=bar\u0026baz=qux"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := tt.toc.JSON()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if string(b) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(b))
			}
		})
	}
}
