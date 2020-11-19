package worker

import (
	"log"

	"github.com/ptcoffee/image-server/database"
	localstorage "github.com/ptcoffee/image-server/storages/local"
)

// Worker struct
type Worker struct {
	db *database.DB
	lc *localstorage.Storage
}

// NewWorker returns a worker instance
func NewWorker(dbURL string) *Worker {
	worker := Worker{
		db: database.New(dbURL),
		lc: localstorage.NewStorage(database.New(dbURL)),
	}
	return &worker
}

// FixImageMetadata will fix image size and file size of image info stored in database
func (worker Worker) FixImageMetadata() error {
	files, err := worker.db.GetFiles([]string{}, nil, nil, []string{})
	if err != nil {
		return err
	}
	count := 0
	for _, file := range files {
		imageSize, _ := worker.lc.GetImageSize(file.Fullname)
		diskSize, _ := worker.lc.GetFileSize(file.Fullname)
		if imageSize.Width != file.Width || imageSize.Height != file.Height || diskSize != file.DiskSize {
			if err := worker.db.Model(&file).Updates(&database.File{
				Width:    imageSize.Width,
				Height:   imageSize.Height,
				DiskSize: diskSize,
			}).Error; err != nil {
				return err
			}
			count++
		}
	}
	log.Printf("Update %d files out of %d", count, len(files))
	return nil
}
