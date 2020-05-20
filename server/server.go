package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/thanhtuan260593/file-server/database"
	localstorage "github.com/thanhtuan260593/file-server/storages/local"
)

//Server struct
type Server struct {
	db      *database.DB
	config  *Config
	storage *localstorage.Storage
	router  *gin.Engine
}

//NewServer will instantiate a new server
func NewServer(db *database.DB) *Server {
	var sv = Server{}
	sv.db = db
	sv.config = NewConfig()
	sv.storage = localstorage.NewStorage(sv.db)
	return &sv
}

//SetupRouter of server
func (s *Server) SetupRouter() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*", "http://localhost:3000"},
		AllowMethods: []string{"POST", "GET", "DELETE", "PUT", "PATCH"},
		AllowHeaders: []string{"*"},
	}))
	imageGroup := router.Group("images")

	// Register public route
	imageGroup.Static("/static", s.storage.WorkingDir)
	imageGroup.GET("/size/:width/:height/:name", s.HandleResize)

	// Register private route
	adminGroup := router.Group("admin")
	adminGroup.GET("/images", s.HandleGetImages)
	adminGroup.GET("/image/:id", s.HandleGetImageByID)
	adminGroup.DELETE("/image/:id", s.HandleDeleteImage)
	adminGroup.PUT("/image", s.HandleUploadImage)
	adminGroup.POST("/image/:id/rename", s.HandleRenameImage)
	adminGroup.POST("/image/:id/replace", s.HandleReplaceImage)
	adminGroup.PUT("/image/:id/tag/:tag", s.HandleAddImageTag)
	adminGroup.DELETE("/image/:id/tag/:tag", s.HandleRemoveImageTag)
	s.router = router
}

//Start server
func (s *Server) Start() *http.Server {
	port := os.Getenv("PORT")
	// Listen and serve on port
	srv := &http.Server{
		Addr:    port,
		Handler: s.router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return srv
}
