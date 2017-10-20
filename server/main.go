package main

import (
	"net"
	"log"
	"github.com/metonimie/simpleFTP/server/server"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Hello world!")
	log.Println("Running on:", "localhost", "port", "8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go server.HandleConnection(conn)
	}
}
