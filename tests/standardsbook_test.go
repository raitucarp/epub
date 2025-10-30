package tests

import (
	"testing"
)

var testData = []epubData{
	{
		name:        "arthur-conan-doyle_the-white-company.epub",
		title:       "The White Company",
		author:      "Arthur Conan Doyle",
		description: "A young English novice in the Middle Ages embarks on action-packed adventure as he joins a motley band of travelers seeking to join with a free company of archers.",
		cover:       "the-white-company.jpg",
		xhtmlCount:  45,
	},
}

func TestOpenReader(t *testing.T) {
	for _, data := range testData {
		t.Run(data.name, data.test())
	}
}
