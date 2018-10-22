package main

import (
	"flag"

	"log"
	"sync"

	"github.com/dnutiu/simplFT/server"
)

func main() {
	flag.StringVar(&server.ConfigName, "config-name", "config",
		"Set the name of the config file.")
	flag.StringVar(&server.ConfigPath, "config", ".", "Set the location of the config file.")
	flag.Parse()

	var wg = new(sync.WaitGroup)

	server.Init()

	wg.Add(2)
	go server.StartUploadServer(wg)
	go server.StartFtpServer(wg)
	wg.Wait()

	log.Println("bye")
}
