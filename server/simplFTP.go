package main

import (
	"github.com/metonimie/simpleFTP/server/server"
)

func main() {
	server.Init()

	go server.StartUploadServer()
	server.StartFtpServer()
}
