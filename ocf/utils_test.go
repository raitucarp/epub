package ocf

import "testing"

func TestGetRootDirectory(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"meta inf file", "META-INF/container.xml", "META-INF"},
		{"meta inf nested", "META-INF/signatures/sig.xml", "META-INF"},
		{"no directory", "mimetype", "mimetype"},
		{"regular directory", "EPUB/package.opf", "EPUB"},
		{"nested directory", "EPUB/text/chapter.xhtml", "EPUB"},
		{"absolute path", "/foo/bar", ""},
		{"dot prefix", "./EPUB/package.opf", "EPUB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRootDirectory(tt.path); got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
