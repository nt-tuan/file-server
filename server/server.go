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

	"github.com/ptcoffee/image-server/cloudflare"
	"github.com/ptcoffee/image-server/docs"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles"

	// swagger embed files
	"github.com/ptcoffee/image-server/database"
	localstorage "github.com/ptcoffee/image-server/storages/local"
)

//Server struct
type Server struct {
	db            *database.DB
	config        *Config
	storage       *localstorage.Storage
	router        *gin.Engine
	port          string
	cloudflareAPI *cloudflare.API
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
	sv.cloudflareAPI = cloudflare.New()
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
	imageGroup.Static("/static", s.storage.WorkingDir)
	imageGroup.GET("/size/:width/:height/*name", s.HandleResize)
	imageGroup.GET("/webp/:width/:height/*name", s.HandleGetWebpImage)

	// History route
	historyGroup := router.Group("history")
	historyGroup.Static("/static", s.storage.HistoryDir)

	// Register private route
	adminGroup := router.Group("/admin")
	adminGroup.Use(func(c *gin.Context) {
		var header struct {
			User string `header:"X-User"`
		}
		if err := c.ShouldBindHeader(&header); err != nil {
			log.Println(err)
		}
		log.Printf("User: %s", header.User)
		c.Set("User", header.User)
		c.Next()
	})
	adminGroup.GET("/images", s.HandleGetImages)
	adminGroup.GET("/images/count", s.HandleCountImages)
	adminGroup.GET("/image/:id", s.HandleGetImageByID)
	adminGroup.DELETE("/image/:id", s.HandleDeleteImage)
	adminGroup.POST("/image", s.HandleUploadImage)
	adminGroup.POST("/image/:id/rename", s.HandleRenameImage)
	adminGroup.POST("/image/:id/replace", s.HandleReplaceImage)
	adminGroup.PUT("/image/:id/tag/:tag", s.HandleAddImageTag)
	adminGroup.DELETE("/image/:id/tag/:tag", s.HandleRemoveImageTag)
	adminGroup.GET("/image/:id/history", s.HandleGetImageHistory)
	adminGroup.GET("/deletedImages", s.HandleGetDeletedFiles)
	adminGroup.POST("/deletedImage/:id/restore", s.HandleRecoverDeletedFile)
	adminGroup.POST("/image/:id/purgeCache", s.HandlePurgeCDNCache)
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
