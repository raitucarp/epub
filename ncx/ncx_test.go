package ncx

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    *NCX
		wantErr bool
	}{
		{
			name: "basic ncx",
			data: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1">
  <head>
    <meta name="dtb:uid" content="urn:uuid:12345"/>
  </head>
  <docTitle>
    <text>Test Title</text>
  </docTitle>
  <navMap>
    <navPoint id="navPoint-1" playOrder="1">
      <navLabel>
        <text>Chapter 1</text>
      </navLabel>
      <content src="chapter1.html"/>
    </navPoint>
  </navMap>
</ncx>`),
			want: &NCX{
				Version: "2005-1",
				Head: Head{
					Meta: []Meta{
						{Name: "dtb:uid", Content: "urn:uuid:12345"},
					},
				},
				DocTitle: TextElement{Text: "Test Title"},
				NavMap: NavMap{
					NavPoints: []NavPoint{
						{
							ID:        "navPoint-1",
							PlayOrder: "1",
							NavLabel:  NavLabel{Text: "Chapter 1"},
							Content:   Content{Src: "chapter1.html"},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "comprehensive ncx",
			data: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<ncx xmlns="http://www.daisy.org/z3986/2005/ncx/" version="2005-1" lang="en">
  <head>
    <meta name="dtb:uid" content="urn:uuid:67890"/>
    <meta name="dtb:depth" content="2"/>
  </head>
  <docTitle>
    <text>Comprehensive Title</text>
  </docTitle>
  <docAuthor>
    <text>John Doe</text>
  </docAuthor>
  <navMap>
    <navPoint id="np-1" class="h1" playOrder="1">
      <navLabel>
        <text>Part 1</text>
      </navLabel>
      <content src="part1.html"/>
      <navPoint id="np-2" class="h2" playOrder="2">
        <navLabel>
          <text>Chapter 1.1</text>
        </navLabel>
        <content src="part1.html#ch1.1"/>
      </navPoint>
    </navPoint>
  </navMap>
  <pageList>
    <pageTarget id="pt-1" type="normal" value="1">
      <navLabel>
        <text>1</text>
      </navLabel>
      <content src="page1.html"/>
    </pageTarget>
  </pageList>
  <navList>
    <navLabel>
      <text>Figures</text>
    </navLabel>
    <navTarget id="nt-1">
      <navLabel>
        <text>Figure 1</text>
      </navLabel>
      <content src="fig1.html"/>
    </navTarget>
  </navList>
</ncx>`),
			want: &NCX{
				Version: "2005-1",
				Lang:    "en",
				Head: Head{
					Meta: []Meta{
						{Name: "dtb:uid", Content: "urn:uuid:67890"},
						{Name: "dtb:depth", Content: "2"},
					},
				},
				DocTitle:  TextElement{Text: "Comprehensive Title"},
				DocAuthor: TextElement{Text: "John Doe"},
				NavMap: NavMap{
					NavPoints: []NavPoint{
						{
							ID:        "np-1",
							Class:     "h1",
							PlayOrder: "1",
							NavLabel:  NavLabel{Text: "Part 1"},
							Content:   Content{Src: "part1.html"},
							NavPoints: []NavPoint{
								{
									ID:        "np-2",
									Class:     "h2",
									PlayOrder: "2",
									NavLabel:  NavLabel{Text: "Chapter 1.1"},
									Content:   Content{Src: "part1.html#ch1.1"},
								},
							},
						},
					},
				},
				PageList: &PageList{
					PageTargets: []PageTarget{
						{
							ID:       "pt-1",
							Type:     "normal",
							Value:    "1",
							NavLabel: NavLabel{Text: "1"},
							Content:  Content{Src: "page1.html"},
						},
					},
				},
				NavLists: []NavList{
					{
						NavLabel: NavLabel{Text: "Figures"},
						NavTargets: []NavTarget{
							{
								ID:       "nt-1",
								NavLabel: NavLabel{Text: "Figure 1"},
								Content:  Content{Src: "fig1.html"},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid xml",
			data: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<ncx>
  <head>
    <meta> <!-- unclosed tag -->
  </head>
</ncx>`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    []byte(""),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil {
					t.Fatal("expected ncx to be non-nil")
				}
				// We don't check XMLName because it's populated automatically by xml.Unmarshal
				// So we set it on tt.want to match what got has to use DeepEqual cleanly
				tt.want.XMLName = got.XMLName

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parse() got = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}
