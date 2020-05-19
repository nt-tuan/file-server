package main

import (
	_ "github.com/lib/pq"
	"github.com/thanhtuan260593/file-server/server"
)

//Main func
func main() {
	server := server.NewServer()
	server.Start()
}
