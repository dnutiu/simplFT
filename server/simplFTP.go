package main

import (
	"flag"

	"github.com/metonimie/simpleFTP/server/server"
)

func main() {
	flag.StringVar(&server.ConfigPath, "config", ".", "Set the location of the config file.")
	flag.Parse()

	server.Init()

	go server.StartUploadServer()
	server.StartFtpServer()
}
