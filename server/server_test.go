package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/ptcoffee/image-server/database"
	localstorage "github.com/ptcoffee/image-server/storages/local"
	"github.com/stretchr/testify/assert"
)

var server *Server
var addedFilePath string = "ERROR_FIFIFIF.JPG"
var code = 1
var dbURL = "postgres://image-server:@:54321/image-server?sslmode=disable"

func setup() {
	db := database.NewClean(dbURL)
	server = NewServer(db)
	reset()
	server.SetupRouter()
}

func reset() {
	server.db = database.NewClean(dbURL)
	server.storage.WorkingDir = testImagesStorageFolder
	server.storage.HistoryDir = testImagesHistoryFolder
	localstorage.RemoveContents(server.storage.WorkingDir)
	localstorage.RemoveContents(server.storage.HistoryDir)
	addedFilePath = filepath.Join(testImageSourceFolder, imageURLs[0].DestName)
}

func TestMain(m *testing.M) {
	downloadTestFiles()
	setup()
	deleteTestFiles()
	code = m.Run()
	os.Exit(code)
}

func performRequest(r http.Handler, method, path string, reader io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, reader)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func performJSONRequest(r http.Handler, method, path string, j gin.H) *httptest.ResponseRecorder {
	data, _ := json.Marshal(j)
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)
	return recorder
}

func requestAddFile(method, url string) (*httptest.ResponseRecorder, error) {
	filename := filepath.Base(addedFilePath)
	b := &bytes.Buffer{}
	var fw io.Writer
	w := multipart.NewWriter(b)
	file, err := os.Open(addedFilePath)
	fw, err = w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return nil, err
	}
	file.Close()
	err = w.WriteField("name", filename)
	if err != nil {
		return nil, err
	}
	w.Close()

	// Perform a GET request with that handler.
	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, err
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())
	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, req)
	return recorder, nil
}

func TestAddFile(t *testing.T) {
	setup()
	recorder, err := requestAddFile("POST", "/admin/image")
	if err != nil {
		assert.Error(t, err)
		return
	}

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAddDuplicatedFile(t *testing.T) {
	reset()
	t.Run("Add new file to delete", TestAddFile)
	recorder, err := requestAddFile("POST", "/admin/image")
	if err != nil {
		assert.Error(t, err)
		return
	}
	fmt.Printf("Response: %v %v\n", recorder.Body.String(), recorder.Header())
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestRenameFile(t *testing.T) {
	t.Run("Add new file to rename", TestAddFile)
	data := gin.H{
		"name": "hehe.jpg",
	}

	recorder := performJSONRequest(server.router,
		"POST",
		fmt.Sprintf("/admin/image/%v/rename", 1),
		data,
	)
	assert.Equal(t, 200, recorder.Code)
}

func TestRenameFileShouldFail(t *testing.T) {
	reset()
	recorder := performJSONRequest(server.router,
		"POST",
		fmt.Sprintf("/admin/image/%v/rename", 1),
		nil,
	)
	assert.Equal(t, 400, recorder.Code)
	recorder = performJSONRequest(server.router,
		"POST",
		fmt.Sprintf("/admin/image/%v/rename", "z"),
		nil,
	)
	assert.Equal(t, 400, recorder.Code)
	t.Run("Add new file to rename", TestAddFile)
	recorder = performJSONRequest(server.router,
		"POST",
		fmt.Sprintf("/admin/image/%v/rename", 1),
		nil,
	)
	assert.Equal(t, 400, recorder.Code)
}

func TestDeleteFile(t *testing.T) {
	reset()
	t.Run("Add new file to delete", TestAddFile)
	recorder := performRequest(server.router, "DELETE", fmt.Sprintf("/admin/image/%v", 1), nil)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestDeleteFileShouldFail(t *testing.T) {
	reset()
	recorder := performRequest(server.router, "DELETE", fmt.Sprintf("/admin/image/%v", 1), nil)
	assert.Equal(t, 400, recorder.Code)
	recorder = performRequest(server.router, "DELETE", fmt.Sprintf("/admin/image/%v", "z"), nil)
	assert.Equal(t, 400, recorder.Code)
}

func TestReplaceFile(t *testing.T) {
	t.Run("Add new file to replace", TestAddFile)
	addedFilePath = filepath.Join(testImageSourceFolder, imageURLs[4].DestName)
	recorder, err := requestAddFile("POST", "/admin/image/1/replace")
	if err != nil {
		assert.Fail(t, err.Error())
		return
	}

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetImagesInfo(t *testing.T) {
	t.Run("Add new file to replace", TestAddFile)
	myURL := "/admin/images?orderBy%5B%5D=id&orderBy%5B%5D=fullname&pageSize=10"
	server.db.LogMode(true)
	recorder := performRequest(server.router, "GET", myURL, nil)
	t.Logf("Response: %v", recorder.Body.String())
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAddTag(t *testing.T) {
	t.Run("Add new file to replace", TestAddFile)
	recorder := performRequest(server.router, "PUT", "/admin/image/1/tag/test_tag", nil)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestAdd2Tags(t *testing.T) {
	t.Run("Add new tag to delete", TestAddTag)
	recorder := performRequest(server.router, "PUT", "/admin/image/1/tag/test_tag", nil)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestRemoveTag(t *testing.T) {
	t.Run("Add new tag to delete", TestAddTag)
	recorder := performRequest(server.router, "DELETE", "/admin/image/1/tag/test_tag", nil)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetImagesHavingTags(t *testing.T) {
	t.Run("Add new tag to get", TestAddTag)
	recorder := performRequest(server.router, "GET", fmt.Sprintf("/admin/images?pageSize=10"), nil)
	t.Logf("Response: %v", recorder.Body.String())
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetImageByID(t *testing.T) {
	t.Run("Add image to get", TestAddFile)
	recorder := performRequest(server.router, "GET", "/admin/image/1", nil)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetImageByIDShouldFail(t *testing.T) {
	reset()
	recorder := performRequest(server.router, "GET", "/admin/image/1", nil)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}

func TestGetResizedImage(t *testing.T) {
	reset()
	addedFilePath = filepath.Join(testImageSourceFolder, "test_resizing_image.png")
	recorder, err := requestAddFile("POST", "/admin/image")
	if err != nil {
		assert.Error(t, err)
		return
	}

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, recorder.Code)

	recorder2 := performRequest(server.router, "GET", "/images/size/400/0/test_resizing_image.png", nil)
	assert.Equal(t, http.StatusOK, recorder2.Code)
}

func TestGetResizedImageShouldWork(t *testing.T) {
	os.Setenv("IMAGE_MAX_WIDTH", "100")
	os.Setenv("IMAGE_MAX_HEIGHT", "100")
	setup()
	addedFilePath = filepath.Join(testImageSourceFolder, "test_resizing_image.png")
	recorder, err := requestAddFile("POST", "/admin/image")
	if err != nil {
		assert.Error(t, err)
		return
	}

	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, recorder.Code)

	recorder2 := performRequest(server.router, "GET", "/images/size/400/0/test_resizing_image.png", nil)
	assert.Equal(t, http.StatusOK, recorder2.Code)
}

func TestGetResizedImageShouldFail(t *testing.T) {
	reset()
	// Assert we encoded correctly,
	// the request gives a 200
	recorder := performRequest(server.router, "GET", "/images/size/400/0/test_resizing_image.png", nil)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	t.Run("Add image to get", TestAddFile)
	recorder = performRequest(server.router, "GET", "/images/size/400/0/IMG_1001.JPG", nil)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
}
