package localstorage

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var testFolder = "../_test"
var testImageSourceFolder = "../../_test/source"
var testImagesStorageFolder = "../../_test/images"
var testImagesHistoryFolder = "../../_test/_history"

type URLImage struct {
	DestName string
	URL      string
}

var imageURLs = []URLImage{
	{
		"IMG_1001.JPG",
		"https://user-images.githubusercontent.com/7953519/82398953-3ee56800-9a7e-11ea-9800-d88b4b58c9de.JPG",
	},
	{
		"IMG_1002.JPG",
		"https://user-images.githubusercontent.com/7953519/82398985-53296500-9a7e-11ea-99fd-3fd52cf1e44c.JPG",
	},
	{
		"IMG_1003.JPG",
		"https://user-images.githubusercontent.com/7953519/82398985-53296500-9a7e-11ea-99fd-3fd52cf1e44c.JPG",
	},
	{
		"IMG_1004.JPG",
		"https://user-images.githubusercontent.com/7953519/82398985-53296500-9a7e-11ea-99fd-3fd52cf1e44c.JPG",
	},
	{
		"IMG_1005.JPG",
		"https://user-images.githubusercontent.com/7953519/82398991-59b7dc80-9a7e-11ea-83f5-c3a4a4bd575d.JPG",
	},
	{
		"IMG_1006.JPG",
		"https://user-images.githubusercontent.com/7953519/82399002-5f152700-9a7e-11ea-9717-8e248f6cd2c3.JPG",
	},
	{
		"IMG_1007.JPG",
		"https://user-images.githubusercontent.com/7953519/82399017-65a39e80-9a7e-11ea-9900-7f9b5b962e56.JPG",
	},
	{
		"IMG_1008.JPG",
		"https://user-images.githubusercontent.com/7953519/82399023-6a685280-9a7e-11ea-8f46-97219b33232e.JPG",
	},
	{
		"IMG_1009.JPG",
		"https://user-images.githubusercontent.com/7953519/82399027-6e947000-9a7e-11ea-8176-cc3f33336371.JPG",
	},
	{
		"IMG_1010.JPG",
		"https://user-images.githubusercontent.com/7953519/82399042-748a5100-9a7e-11ea-9e4b-353b88209544.JPG",
	},
	{
		"IMG_1011.JPG",
		"https://user-images.githubusercontent.com/7953519/82399059-7c49f580-9a7e-11ea-848b-8128665ea6c2.JPG",
	},
	{
		"test_resizing_image.png",
		"https://user-images.githubusercontent.com/7953519/82418631-6d2b6d80-9aa7-11ea-9803-c0a0e7379276.png",
	},
}

func downloadTestFiles() error {
	//Create test folder
	var err error
	err = os.MkdirAll(testImageSourceFolder, os.ModePerm)
	err = os.MkdirAll(testImagesStorageFolder, os.ModePerm)
	err = os.MkdirAll(testImagesHistoryFolder, os.ModePerm)
	RemoveContents(testImagesStorageFolder)
	RemoveContents(testImagesHistoryFolder)
	if err != nil {
		fmt.Println(err)
		return err
	}

	//Download and save files
	for _, img := range imageURLs {
		localPath := filepath.Join(testImageSourceFolder, img.DestName)
		if _, err := os.Stat(localPath); err == nil {
			continue
		}
		resp, err := http.Get(img.URL)
		if err != nil {
			fmt.Println(err)
			return err
		}
		out, err := os.Create(localPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		io.Copy(out, resp.Body)
		resp.Body.Close()

	}
	return nil
}

func deleteTestFiles() {
	return
}
