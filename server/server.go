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

	"github.com/ptcoffee/image-server/docs"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"

	// swagger embed files
	"github.com/ptcoffee/image-server/database"
	localstorage "github.com/ptcoffee/image-server/storages/local"
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
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AddAllowHeaders("authorization")
	router.Use(cors.New(config))
	imageGroup := router.Group("images")
	// Register public route
	imageGroup.Use(cacheHeader).Static("/static", s.storage.WorkingDir)
	imageGroup.Use(cacheHeader).GET("/size/:width/:height/*name", s.HandleResize)

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
	router.GET("/health/ready", func(c *gin.Context) {
		c.JSON(200, map[string]bool{"ok": true})
	})
	s.router = router
}

//Start server
func (s *Server) Start() {
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
}

func cacheHeader(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=108000")
}
