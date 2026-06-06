package ocf

import (
	"testing"
)

func TestOCFZipContainer_MimeType(t *testing.T) {
	container := &OCFZipContainer{
		files: map[string][]byte{
			"mimetype": []byte(MimeType),
		},
	}

	mime := container.MimeType()
	if mime != MimeType {
		t.Errorf("Expected MimeType to be %s, but got %s", MimeType, mime)
	}
}

func TestOCFZipContainer_MimeType_Empty(t *testing.T) {
	container := &OCFZipContainer{
		files: map[string][]byte{},
	}

	mime := container.MimeType()
	if mime != "" {
		t.Errorf("Expected empty MimeType, but got %s", mime)
	}
}

func TestOCFZipContainer_MimeType_Custom(t *testing.T) {
	container := &OCFZipContainer{
		files: map[string][]byte{
			"mimetype": []byte("application/custom+zip"),
		},
	}

	mime := container.MimeType()
	if mime != "application/custom+zip" {
		t.Errorf("Expected MimeType to be application/custom+zip, but got %s", mime)
	}
}

func TestOCFZipContainer_MimeType_NilMap(t *testing.T) {
	container := &OCFZipContainer{
		files: nil,
	}

	mime := container.MimeType()
	if mime != "" {
		t.Errorf("Expected empty MimeType for nil map, but got %s", mime)
	}
}

func TestOCFZipContainer_Getters(t *testing.T) {
	container := &OCFZipContainer{
		metaInf: MetaInf{
			container:  Container{},
			signatures: Signatures{},
			encryption: Encryption{},
			metadata:   Metadata{},
			rights:     Rights{},
			manifest:   Manifest{},
		},
	}

	if container.Container() == nil {
		t.Errorf("Expected Container to not be nil")
	}

	if container.Signatures() == nil {
		t.Errorf("Expected Signatures to not be nil")
	}

	if container.Encryption() == nil {
		t.Errorf("Expected Encryption to not be nil")
	}

	if container.Metadata() == nil {
		t.Errorf("Expected Metadata to not be nil")
	}

	if container.Rights() == nil {
		t.Errorf("Expected Rights to not be nil")
	}

	if container.Manifest() == nil {
		t.Errorf("Expected Manifest to not be nil")
	}
}

func TestOCFZipContainer_AllFiles(t *testing.T) {
	files := map[string][]byte{
		"mimetype": []byte(MimeType),
		"META-INF/container.xml": []byte("container"),
	}
	container := &OCFZipContainer{
		files: files,
	}

	allFiles := container.AllFiles()
	if len(allFiles) != 2 {
		t.Errorf("Expected AllFiles to return 2 files, got %d", len(allFiles))
	}
}

func TestOCFZipContainer_SelectFile(t *testing.T) {
	files := map[string][]byte{
		"mimetype": []byte(MimeType),
	}
	container := &OCFZipContainer{
		files: files,
	}

	data, err := container.SelectFile("mimetype")
	if err != nil {
		t.Errorf("Expected SelectFile to return no error, got %v", err)
	}
	if string(data) != MimeType {
		t.Errorf("Expected SelectFile to return %s, got %s", MimeType, string(data))
	}

	_, err = container.SelectFile("notfound")
	if err == nil {
		t.Errorf("Expected SelectFile to return error for notfound file")
	}
}

func TestOCFZipContainer_NonMetaInfFiles(t *testing.T) {
	files := map[string][]byte{
		"mimetype": []byte(MimeType),
		"META-INF/container.xml": []byte("container"),
		"OEBPS/content.opf": []byte("content"),
	}
	container := &OCFZipContainer{
		files: files,
	}

	nonMetaInfFiles := container.NonMetaInfFiles()
	if len(nonMetaInfFiles) != 2 {
		t.Errorf("Expected NonMetaInfFiles to return 2 files, got %d", len(nonMetaInfFiles))
	}
	if _, ok := nonMetaInfFiles["META-INF/container.xml"]; ok {
		t.Errorf("Expected NonMetaInfFiles to not return META-INF/container.xml")
	}
}
