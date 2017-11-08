package server

import (
	"bufio"
	"io"
	"log"
	"net"
)

// ClientStack interface holds the function needed to work with a stack.
type ClientStack interface {
	Stack() Stack // Returns the underlying Stack.
}

// Client interface provides the blueprints for the Client that is used by the server.
type Client interface {
	Connection() net.Conn        // Connection returns the connection stream.
	SetConnection(conn net.Conn) // SetConnection sets the connection for the client.
	Disconnect()                 // Disconnect closes the Client's connections and clears up resources.
}

// FTPClient represents a FTPClient connection, it holds a root cage and the underlying connection.
type FTPClient struct {
	rootCage   *StringStack // rootCage is a StringStack that is used to represent the current directory the client is in.
	connection net.Conn
}

// Stack returns the root cage stack.
func (c *FTPClient) Stack() Stack {
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
	io.WriteString(client.Connection(), "Hello and welcome to simple ftp\n")

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
