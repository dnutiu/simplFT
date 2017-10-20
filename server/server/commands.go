package server

import (
	"net"
	"strings"
	"os"
	"log"
	"io/ioutil"
	"bytes"
	"strconv"
)

// PATH is the constant which should contain the fixed path where the simpleFTP server will run
// This will act like a root cage.
const PATH = "/Users/denis/GoglandProjects/golangBook/GoRoutines/"

// SendFile sends the file to the client and returns true if it succeeds and false otherwise.
func SendFile(c net.Conn, path string) (int, error) {
	var fileName string

	// Make sure the user can't request any files on the system.
	lastForwardSlash := strings.LastIndex(path, "/")
	if lastForwardSlash != -1 {
		// Eliminate the last forward slash i.e ../../asdas will become asdas
		fileName = path[lastForwardSlash+1:]
	} else {
		fileName = path
	}

	file, err := os.Open(PATH + fileName)
	if err != nil {
		// Open file failed.
		log.Println(err)
		return 0, err
	}
	defer file.Close() // Closing the fd when the function has exited.

	data, err := ioutil.ReadAll(file)
	n, err := c.Write(data)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	// How is this even possible?
	if n == 0 {
		log.Println("0 bits written for:", path)
		return 0, nil
	}

	return n, nil
}

// ListFiles list the files from path and sends them to the connection
func ListFiles(c net.Conn) error {
	files, err := ioutil.ReadDir(PATH)
	if err != nil {
		return err;
	}

	buffer := bytes.NewBufferString("Directory Mode Size LastModified Name\n")
	for _, f := range files {
		buffer.WriteString(strconv.FormatBool(f.IsDir()) + " " + string(f.Mode().String()) + " " +
			strconv.FormatInt(f.Size(), 10) + " " + f.ModTime().String() + " " + string(f.Name()) + " " + "\n")
	}

	_, err = c.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func ShowHelp(c net.Conn) error {
	var helpText string = `
The available commands are:
get <filename> - Download the requested filename.
ls - List the files in the current directory.
clear - Clear the screen.
exit - Close the connection with the server.
`
	_, err := c.Write([]byte(helpText))

	return err
}
