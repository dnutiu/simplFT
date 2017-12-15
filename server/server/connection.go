package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"os"

	"time"

	"os/signal"
	"syscall"

	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/metonimie/simpleFTP/server/server/config"
	"github.com/spf13/viper"
)

// DataBufferSize the maximum size of the data buffer.
// The data buffer is used at reading from files, the buffer
// is also send to the client.
const DataBufferSize = 1024 * 1024

// ConfigPath is used by the config package to find the config file.
var ConfigPath string

// uploadDirectory is the directory where the files will be uploaded
var uploadDirectory string

// uploadTimeout is the amount in seconds the server will wait for a file to be uploaded
var uploadTimeout time.Duration

// uploadEnabled holds true of false, if in the config upload was enabled when the upload server was started.
var uploadEnabled bool

// Shutdown is the shutdown where SIGINT and SIGTERM is send too
var Shutdown = make(chan os.Signal)
var ftpShutdown = make(chan struct{})
var uploadShutdown = make(chan struct{})

var uploadListener net.Listener
var ftpListener net.Listener

// All connected clients
var clients map[Client]bool

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

func shutdownHandler() {
	for {
		select {
		case <-Shutdown:
			log.Println("Shutdown signal received")
			var wg sync.WaitGroup
			wg.Add(1)

			go func() { // Disconnect all the clients.
				for k := range clients {
					k.Disconnect()
				}
				wg.Done()
			}()
			wg.Wait()

			ShutdownUploadServer()
			ShutdownFtpServer()
			return
		}
	}
}

func ShutdownUploadServer() {
	if uploadEnabled == true {
		if uploadListener != nil {
			uploadListener.Close()
		}
		uploadShutdown <- struct{}{}
	}
}

func ShutdownFtpServer() {
	if ftpListener != nil {
		ftpListener.Close()
	}
	ftpShutdown <- struct{}{}
}

func Init() {
	signal.Notify(Shutdown, syscall.SIGINT, syscall.SIGTERM)

	clients = make(map[Client]bool)
	go shutdownHandler()

	config.InitializeConfiguration("config", ConfigPath)
	config.ChangeCallback(func(event fsnotify.Event) {
		log.Println("Configuration reloaded successfully!")
	})
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
	clients[client] = true

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
	delete(clients, client)
	log.Println(client.Connection().RemoteAddr(), "has disconnected.")
}

func StartFtpServer(wg *sync.WaitGroup) error {
	defer wg.Done()
	Addr := viper.GetString("address")
	Port := viper.GetInt("port")
	DirDepth := viper.GetInt("maxDirDepth")
	BasePath = viper.GetString("absoluteServePath")

	// Start the server
	var err error
	ftpListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", Addr, Port))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer ftpListener.Close()

	log.Println("Hello world!")
	log.Println("Ftp server running on:", Addr, "port", Port)

	for {
		conn, err := ftpListener.Accept()

		// Handle shutdown
		select {
		case <-ftpShutdown:
			goto exit
		default:
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
exit:
	log.Println("Ftp server exited.")
	return nil
}

func HandleUpload(conn net.Conn) {
	// Initialize Client
	client := FTPClient{}
	client.SetStack(MakeStringStack(2))

	// Upload directory
	client.Stack().Push(uploadDirectory)
	client.SetConnection(conn)
	defer client.Disconnect()

	// Create the file on the disk and make sure that the filename is random.
	var filename = randSeq(10)
	var _, err = os.Stat(MakePathFromStringStack(client.Stack()) + filename)

	for !os.IsNotExist(err) {
		filename = randSeq(10)
		_, err = os.Stat(MakePathFromStringStack(client.Stack()) + filename)
	}

	// This channel will be used to store the uploadResult
	c1 := make(chan error, 1)
	log.Println(conn.RemoteAddr().String() + " is uploading something.")
	clients[&client] = true

	// Create a new Go routine for uploading
	go func() {
		err := UploadFile(&client, filename)
		c1 <- err
	}()

	// Wait for either UploadResult or Timeout
	select {
	case result := <-c1:
		{
			if result == nil {
				io.WriteString(conn, filename)
				log.Println(conn.RemoteAddr().String() + "'s upload finished.")
			} else {
				log.Println(fmt.Sprintf("%s: %s %s", "HandleUpload", conn.RemoteAddr().String(), result.Error()))

				client.Stack().Push(filename)
				os.Remove(MakePathFromStringStack(client.Stack()))

				io.WriteString(conn, result.Error())
			}

			conn.Close()
		}
	case <-time.After(time.Second * uploadTimeout):
		{
			io.WriteString(conn, "Timeout")
			conn.Close()
		}
	}

	delete(clients, &client)
}

// StartUploadServer starts the uploading server
func StartUploadServer(wg *sync.WaitGroup) error {
	defer wg.Done()
	var err error
	uploadEnabled = viper.GetBool("upload.enabled")
	if uploadEnabled == false {
		log.Println("Uploading not enabled. To enable modify the config file and restart the server")
		return ErrUploadServerFailure
	}

	addr := viper.GetString("upload.address")
	port := viper.GetInt("upload.port")
	uploadDirectory = viper.GetString("upload.directory")
	uploadTimeout = time.Duration(viper.GetInt("upload.timeout"))

	uploadListener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Println(err)
		return err
	}
	defer uploadListener.Close()

	err = os.Mkdir(uploadDirectory, 0740)
	if err != nil {
		if _, err := os.Stat(uploadDirectory); err != nil {
			if os.IsNotExist(err) {
				log.Println("Can't create upload directory!")
				return err
			}
		}
	}

	log.Println("Upload server running on:", addr, "port", port)

	for {
		conn, err := uploadListener.Accept()
		// Handle shutdown
		select {
		case <-uploadShutdown:
			goto exit
		default:
			if err != nil {
				log.Print(err)
				continue
			}

			go HandleUpload(conn)
		}
	}

exit:
	log.Println("Upload server exited.")
	return nil
}
