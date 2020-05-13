package server

import (
	"os"

	"github.com/thanhtuan260593/file-server/database"
	"github.com/thanhtuan260593/file-server/storages/local"
)

//Server struct
type Server struct {
	db      *database.DB
	config  *Config
	storage *local.Local
}

//NewServer will instantiate a new server
func NewServer() *Server {
	var sv = Server{}

	dbURL := os.Getenv("DATABASE_URL")
	sv.db = database.New(dbURL)
	sv.config = NewConfig()
	sv.storage = local.NewLocal(sv.db)
	return &sv
}
