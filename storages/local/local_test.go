package localstorage

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
	"github.com/ptcoffee/image-server/database"
	"github.com/twinj/uuid"
)

var store *Storage
var addedFile, replacedFile URLImage
var dbURL = "postgres://image-server:@:54321/image-server?sslmode=disable"

func loadExternalFile() {
	addedFile = imageURLs[0]
	replacedFile = imageURLs[1]
}

func setup() {
	downloadTestFiles()
	reset()
}

func reset() {
	db := database.NewClean(dbURL)
	store = NewStorage(db)
	store.WorkingDir = testImagesStorageFolder
	store.HistoryDir = testImagesHistoryFolder
	RemoveContents(store.WorkingDir)
	RemoveContents(store.HistoryDir)
	loadExternalFile()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestAddFile(t *testing.T) {
	reset()
	path := filepath.Join(testImageSourceFolder, addedFile.DestName)
	reader, err := os.Open(path)
	if err != nil {
		t.Error(err)
		return
	}

	// var fname string
	_, err = store.AddFile(reader, addedFile.DestName)
	if err != nil {
		t.Error(err)
		return
	}
	reader.Close()
}

func TestRemoveFile(t *testing.T) {
	reset()
	path := filepath.Join(testImageSourceFolder, addedFile.DestName)
	reader, err := os.Open(path)
	if err != nil {
		t.Error(err)
		return
	}
	// var fname string
	file, err := store.AddFile(reader, addedFile.DestName)
	if err != nil {
		t.Error(err)
		return
	}
	if err := store.DeleteFile(file); err != nil {
		t.Error(err)
		return
	}
}

func TestRenameFile(t *testing.T) {
	t.Run("Create file to rename file", TestAddFile)
	newName := uuid.NewV4().String() + ".jpg"
	if _, err := store.RenameFile(addedFile.DestName, newName); err != nil {
		t.Error(err)
		return
	}
	// t.Logf("Rename %v to %v", addingFile, newName)
}

func TestReplaceFile(t *testing.T) {
	reset()
	path := filepath.Join(testImageSourceFolder, addedFile.DestName)
	reader, err := os.Open(path)
	if err != nil {
		t.Error(err)
		return
	}
	// var fname string
	file, err := store.AddFile(reader, addedFile.DestName)
	if err != nil {
		t.Error(err)
		return
	}
	replacedPath := filepath.Join(testImageSourceFolder, replacedFile.DestName)
	replaceReader, err := os.Open(replacedPath)
	if err != nil {
		t.Error(err)
		return
	}
	if _, err := store.ReplaceFile(file, replaceReader); err != nil {
		t.Error(err)
		return
	}
}

func TestCreateMissingFiles(t *testing.T) {
	reset()
	store.WorkingDir = testImageSourceFolder
	store.CreateMissingFiles()
}
