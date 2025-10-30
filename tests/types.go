package tests

import (
	"path"
	"testing"

	"github.com/raitucarp/epub"
)

type epubData struct {
	name        string
	title       string
	description string
	author      string
	cover       string
	xhtmlCount  int

	epub epub.Reader
}

const epubPath = "./data"
const coverPath = "./cover"

func (data *epubData) attachEpub(reader epub.Reader) {
	data.epub = reader
}

func (data *epubData) testTitle() func(t *testing.T) {
	return func(t *testing.T) {
		actual := data.epub.Title()
		expected := data.title

		if actual != expected {
			t.Errorf("Title is not equal, actual = %s, expected = %s", actual, expected)
		}

	}
}

func (data *epubData) testAuthor() func(t *testing.T) {
	return func(t *testing.T) {
		actual := data.epub.Author()
		expected := data.author

		if actual != expected {
			t.Errorf("Author is not equal, actual = %s, expected = %s", actual, expected)
		}
	}
}

func (data *epubData) testDescription() func(t *testing.T) {
	return func(t *testing.T) {
		actual := data.epub.Description()
		expected := data.description

		if actual != expected {
			t.Errorf("Description is not equal, actual = %s, expected = %s", actual, expected)
		}
	}
}

func (data *epubData) testCover() func(t *testing.T) {
	return func(t *testing.T) {
		actual := data.epub.Cover()
		expected, err := readCover(data.cover)
		if err != nil {
			t.Errorf("Something error when reading cover %s", err)
		}

		if !isCoverEqual(*actual, expected) {
			t.Errorf("Cover is not equal, actual = %s, expected %s", (*actual).Bounds().String(), expected.Bounds().String())
		}
	}
}

func (data *epubData) test() func(t *testing.T) {
	return func(t *testing.T) {
		epub, err := epub.OpenReader(path.Join(epubPath, data.name))
		if err != nil {
			t.Errorf("Something error %s", err)
			return
		}

		data.attachEpub(epub)

		t.Run("title", data.testTitle())
		t.Run("author", data.testAuthor())
		t.Run("description", data.testDescription())
		t.Run("cover", data.testCover())

		t.Run("xhtml count", func(t *testing.T) {
			actual := len(epub.ContentDocumentXHTML())
			expected := data.xhtmlCount

			if actual != expected {
				t.Errorf("XHTML Contents mismatch count actual %d, expected %d", actual, expected)
			}
		})
	}
}
