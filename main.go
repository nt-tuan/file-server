package main

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/thanhtuan260593/file-server/database"
	"github.com/thanhtuan260593/file-server/server"
)

//Main func
func main() {
	dbURL := os.Getenv("DATABASE_URL")
	db := database.New(dbURL)
	server := server.NewServer(db)
	server.Start()
}
