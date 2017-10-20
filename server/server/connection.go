package server

import (
	"net"
	"io"
	"log"
	"bufio"
)

func HandleConnection(c net.Conn) {
	defer c.Close()
	io.WriteString(c, "Hello and welcome to simple ftp\n")

	log.Println(c.RemoteAddr(), "has connected.")

	// Process input
	input := bufio.NewScanner(c)
	for input.Scan() {
		log.Println(c.RemoteAddr(), ":", input.Text())

		err := ProcessInput(c, input.Text())
		if err != nil {
			log.Println(err)
			io.WriteString(c, err.Error()+"\n")
		}
	}

	// Client has left.
	log.Println(c.RemoteAddr(), "has disconnected.")
}
