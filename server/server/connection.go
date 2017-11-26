package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/fsnotify/fsnotify"
	"github.com/metonimie/simpleFTP/server/server/config"
	"github.com/spf13/viper"
)

// DataBufferSize the maximum size of the data buffer.
// The data buffer is used at reading from files, the buffer
// is also send to the client.
const DataBufferSize = 1024 * 1024

// Client interface provides the blueprints for the Client that is used by the server.
type Client interface {
	Connection() net.Conn        // Connection returns the connection stream.
	SetConnection(conn net.Conn) // SetConnection sets the connection for the client.
	Disconnect()                 // Disconnect closes the Client's connections and clears up resources.
	Stack() *StringStack         // Returns the underlying String Stack.
}

// FTPClient represents a FTPClient connection, it holds a root cage and the underlying connection.
type FTPClient struct {
	rootCage   *StringStack // rootCage is a StringStack that is used to represent the current directory the client is in.
	connection net.Conn
}

// Stack returns the root cage stack.
func (c *FTPClient) Stack() *StringStack {
	return c.rootCage
}

// SetStack sets the stack for the FTPClient.
func (c *FTPClient) SetStack(stack *StringStack) {
	c.rootCage = stack
}

// Connection returns the Connection of the client.
func (c *FTPClient) Connection() net.Conn {
	return c.connection
}

// SetConnection sets the given connection to the client.
func (c *FTPClient) SetConnection(conn net.Conn) {
	c.connection = conn
}

// Disconnects the client.
func (c *FTPClient) Disconnect() {
	c.connection.Close()
}

func HandleConnection(client Client) {
	defer client.Disconnect()
	defer func() {
		if r := recover(); r != nil {
			log.Println("PANIC: ", r)

			recoveryError, ok := r.(string)
			if ok {
				io.WriteString(client.Connection(), fmt.Sprintf("PANIC: %s", recoveryError))
			}
		}
	}()

	log.Println(client.Connection().RemoteAddr(), "has connected.")

	// Process input
	input := bufio.NewScanner(client.Connection())

	for input.Scan() {
		log.Println(client.Connection().RemoteAddr(), ":", input.Text())

		err := ProcessInput(client, input.Text())
		if err != nil {
			log.Println(err)
			io.WriteString(client.Connection(), err.Error()+"\n")
		}
	}

	// Client has left.
	log.Println(client.Connection().RemoteAddr(), "has disconnected.")
}

func Init() {
	config.InitializeConfiguration()
	config.ConfigChangeCallback(func(event fsnotify.Event) {
		log.Println("Configuration reloaded successfully!")
	})
}

func StartFtpServer() {
	Addr := viper.GetString("address")
	Port := viper.GetInt("port")
	DirDepth := viper.GetInt("maxDirDepth")
	BasePath = viper.GetString("absoluteServePath")

	// Start the server
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", Addr, Port))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Hello world!")
	log.Println("Ftp server running on:", Addr, "port", Port)

	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		client := FTPClient{}
		client.SetStack(MakeStringStack(DirDepth))
		client.SetConnection(conn)

		go HandleConnection(&client)
	}
}

func StartUploadServer() {
	if viper.GetBool("upload.enabled") == false {
		log.Println("Uploading not enabled. To enable modify the config file and restart the server")
		return
	}

	Addr := viper.GetString("upload.address")
	Port := viper.GetInt("upload.port")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", Addr, Port))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Upload server running on:", Addr, "port", Port)

	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		client := FTPClient{}
		client.SetStack(MakeStringStack(1))
		client.SetConnection(conn)

		log.Println(conn.RemoteAddr().String() + " is uploading something.")

		filename, err := UploadFile(&client)
		if err == nil {
			io.WriteString(conn, filename)
		} else {
			log.Print(conn.RemoteAddr().String())
			log.Println(err)
		}

		log.Println(conn.RemoteAddr().String() + "'s upload finished.")
	}
}
