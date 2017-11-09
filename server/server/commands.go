package server

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// GetFile sends the file to the client and returns true if it succeeds and false otherwise.
func GetFile(c Client, path string) (int, error) {
	fileName, sanitized := sanitizeFilePath(path)
	if sanitized {
		return 0, ErrSlashNotAllowed
	}

	stack, ok := c.Stack().(*StringStack)
	if !ok {
		return 0, ErrStackCast
	}

	file, err := os.Open(MakePathFromStringStack(stack) + fileName)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer file.Close()

	data, err := readFileData(file)
	if err != nil {
		return 0, err
	}

	n, err := c.Connection().Write(data)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if n == 0 {
		// This happens when the user ties to get the current directory
		return 0, GetNoBitsError
	}
	return n, nil
}

func readFileData(file *os.File) ([]byte, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

func sanitizeFilePath(path string) (string, bool) {
	var fileName string
	var sanitized bool
	// Make sure the user can't request any files on the system.
	lastForwardSlash := strings.LastIndex(path, "/")
	if lastForwardSlash != -1 {
		// Eliminate the last forward slash i.e ../../asdas will become asdas
		fileName = path[lastForwardSlash+1:]
		sanitized = true
	} else {
		fileName = path
		sanitized = false
	}
	return fileName, sanitized
}

// ListFiles list the files from path and sends them to the connection
func ListFiles(c Client) error {
	stack, ok := c.Stack().(*StringStack)
	if !ok {
		return ErrStackCast
	}

	files, err := ioutil.ReadDir(MakePathFromStringStack(stack))
	if err != nil {
		return err
	}

	buffer := bytes.NewBufferString("Directory Mode Size LastModified Name\n")
	for _, f := range files {
		buffer.WriteString(strconv.FormatBool(f.IsDir()) + " " + string(f.Mode().String()) + " " +
			strconv.FormatInt(f.Size(), 10) + " " + f.ModTime().String() + " " + string(f.Name()) + " " + "\n")
	}

	_, err = c.Connection().Write(buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// ClearScreen cleans the client's screen by sending clear to the terminal.
func ClearScreen(c Client) error {
	// Ansi clear: 1b 5b 48 1b 5b 4a
	// clear | hexdump -C
	var b = []byte{0x1b, 0x5b, 0x48, 0x1b, 0x5b, 0x4a}
	_, err := c.Connection().Write(b)

	return err
}

// ChangeDirectoryCommand changes the directory to the given directory
func ChangeDirectoryCommand(c Client, directory string) error {
	var err error
	path, sanitized := sanitizeFilePath(directory)
	if sanitized {
		return ErrSlashNotAllowed
	}

	stack, ok := c.Stack().(*StringStack)
	if !ok {
		return ErrStackCast
	}

	if path == "." {
		err = nil
	} else if path == ".." {
		err = ChangeDirectoryToPrevious(stack)
	} else {
		err = ChangeDirectory(stack, path)
	}

	return err
}

func ShowHelp(c Client) error {
	var helpText = `
The available commands are:
get <filename> - Download the requested filename.
ls             - List the files in the current directory.
cd             - Changes the directory.
clear          - Clear the screen.
exit           - Close the connection with the server.c
`
	_, err := c.Connection().Write([]byte(helpText))

	return err
}
