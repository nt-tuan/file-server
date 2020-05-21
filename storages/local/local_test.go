package localstorage

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/lib/pq"
	"github.com/thanhtuan260593/file-server/database"
	"github.com/twinj/uuid"
)

var store *Storage
var addingFile string = "test.jpg"
var dbURL = "postgres://file-server:@:54321/file-server?sslmode=disable"

func setup() {
	db := database.NewClean(dbURL)
	store = NewStorage(db)
	store.WorkingDir = testImagesStorageFolder
	store.HistoryDir = testImagesHistoryFolder
	RemoveContents(store.WorkingDir)
	RemoveContents(store.HistoryDir)
	addingFile = filepath.Join(testImageSourceFolder, imageURLs[0].DestName)
}
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestAddFile(t *testing.T) {
	addingFile = uuid.NewV4().String() + ".jpg"
	reader, err := os.Open("../../files/_source/IMG_1001.JPG")
	if err != nil {
		t.Error(err)
		return
	}
	defer reader.Close()
	// var fname string
	_, err = store.AddFile(reader, addingFile)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = store.AddFile(reader, addingFile)
	if err == nil {
		t.FailNow()
	}

	if !errors.Is(err, ErrFileExisted) {
		t.Error(err)
		return
	}
}

func TestRemoveFile(t *testing.T) {
	t.Run("Create file to remove file", TestAddFile)
	historyFile, err := store.DeleteFile(addingFile)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(historyFile)
}

func TestRenameFile(t *testing.T) {
	t.Run("Create file to rename file", TestAddFile)
	newName := uuid.NewV4().String() + ".jpg"
	if _, err := store.RenameFile(addingFile, newName); err != nil {
		t.Error(err)
		return
	}
	// t.Logf("Rename %v to %v", addingFile, newName)
}

func TestReplaceFile(t *testing.T) {
	reader, err := os.Open("../../files/_source/IMG_1004.JPG")
	if err != nil {
		t.Error(err)
		return
	}
	t.Run("Create file to rename file", TestAddFile)
	if _, err := store.ReplaceFile(addingFile, reader); err != nil {
		t.Error(err)
		return
	}
}

func TestCreateMissingFiles(t *testing.T) {
	setup()
	store.WorkingDir = testImageSourceFolder
	store.CreateMissingFiles()
}
