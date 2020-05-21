package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/thanhtuan260593/file-server/docs"

	// swagger embed files
	"github.com/thanhtuan260593/file-server/database"
	localstorage "github.com/thanhtuan260593/file-server/storages/local"
)

//Server struct
type Server struct {
	db      *database.DB
	config  *Config
	storage *localstorage.Storage
	router  *gin.Engine
	port    string
}

//NewServer will instantiate a new server
func NewServer(db *database.DB) *Server {
	var sv = Server{}
	sv.db = db
	sv.config = NewConfig()
	sv.storage = localstorage.NewStorage(sv.db)
	sv.port = ":5000"
	port := os.Getenv("PORT")
	if port != "" {
		sv.port = port
	}
	if logStr := os.Getenv("DATABASE_LOG_MODE"); logStr != "" {
		if v, err := strconv.ParseBool(logStr); err == nil && v {
			sv.db.LogMode(true)
		}
	}
	return &sv
}

func setupSwaggerInfo() {
	host := os.Getenv("HOST")
	basePath := os.Getenv("BASE_PATH")
	title := os.Getenv("SWAGGER_TITLE")
	version := os.Getenv("SWAGGER_VERSION")
	docs.SwaggerInfo.Version = "1.0"
	if host != "" {
		docs.SwaggerInfo.Host = host
	}
	if basePath != "" {
		docs.SwaggerInfo.BasePath = basePath
	}
	if title != "" {
		docs.SwaggerInfo.Title = title
	}
	if version != "" {
		docs.SwaggerInfo.Version = version
	}
}

//SetupRouter of server
func (s *Server) SetupRouter() {
	// programatically set swagger info
	setupSwaggerInfo()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*", "http://localhost:3000"},
		AllowMethods: []string{"POST", "GET", "DELETE", "PUT", "PATCH"},
		AllowHeaders: []string{"*"},
	}))
	imageGroup := router.Group("images")

	// Register public route
	imageGroup.Static("/static", s.storage.WorkingDir)
	imageGroup.GET("/size/:width/:height/*name", s.HandleResize)

	// Register private route
	adminGroup := router.Group("/admin")
	adminGroup.GET("/images", s.HandleGetImages)
	adminGroup.GET("/image/:id", s.HandleGetImageByID)
	adminGroup.DELETE("/image/:id", s.HandleDeleteImage)
	adminGroup.PUT("/image", s.HandleUploadImage)
	adminGroup.POST("/image/:id/rename", s.HandleRenameImage)
	adminGroup.POST("/image/:id/replace", s.HandleReplaceImage)
	adminGroup.PUT("/image/:id/tag/:tag", s.HandleAddImageTag)
	adminGroup.DELETE("/image/:id/tag/:tag", s.HandleRemoveImageTag)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	s.router = router
}

//Start server
func (s *Server) Start() *http.Server {
	// Listen and serve on port
	//s.router.Run(s.port)
	srv := &http.Server{
		Addr:    s.port,
		Handler: s.router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
	return srv
}
