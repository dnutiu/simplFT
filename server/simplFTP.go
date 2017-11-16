package main

import (
	"log"
	"net"

	"github.com/metonimie/simpleFTP/server/server"
	"github.com/spf13/viper"
)

func main() {
	server.InitializedConfiguration()

	Addr := viper.GetString("Address")
	Port := viper.GetString("Port")
	DirDepth := viper.GetInt("MaxDirDepth")

	// Start the server
	listener, err := net.Listen("tcp", Addr+":"+Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Hello world!")
	log.Println("Running on:", Addr, "port", Port)

	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		client := server.FTPClient{}
		client.SetStack(server.MakeStringStack(DirDepth))
		client.SetConnection(conn)

		go server.HandleConnection(&client)
	}
}
