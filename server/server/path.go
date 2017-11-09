package server

// The Path module should work with the stack by providing functions that update the stack
// with information that helps keeping track of the directories that the clients use.

import (
	"bytes"
	"os"
)

func containsSlash(dir string) bool {
	for _, c := range dir {
		if c == '/' {
			return true
		}
	}
	return false
}

// PATH is the constant which should contain the fixed path where the simpleFTP server will run
// This will act like a root cage. It must end with a forward slash!
var BasePath = "/Users/denis/GoglandProjects/golangBook/"

// GetPath will return the complete path
//func GetPath(stack *StringStack) string {
//	if stack.IsEmpty() {
//		return BasePath
//	}
//
//	return BasePath
//}

// MakePathFromStringStack gets a StringStack and makes a path.
func MakePathFromStringStack(stack *StringStack) string {
	var buffer bytes.Buffer

	buffer.WriteString(BasePath)
	if BasePath[len(BasePath)-1] != '/' {
		buffer.WriteString("/")
	}

	for i := 0; i < stack.Size(); i++ {
		buffer.WriteString(stack.Items()[i] + "/")
	}

	return buffer.String()
}

// ChangeDirectory changes the current working directory with respect to BasePath
func ChangeDirectory(stack *StringStack, directory string) error {
	if containsSlash(directory) == true {
		return ErrInvalidDirectoryName
	}
	stack.Push(directory)

	path := MakePathFromStringStack(stack)
	fileInfo, err := os.Stat(path)
	if err != nil {
		stack.Pop()
		return err
	}

	if fileInfo.IsDir() == false {
		stack.Pop()
		return ErrNotADirectory
	}

	// The last 9 bits represent the Unix perms format rwxrwxrwx
	perms := fileInfo.Mode().Perm()

	// The user has permission to view the directory
	if (perms&1 != 0) && (perms&(1<<2) != 0) || (perms&(1<<6) != 0) && (perms&(1<<8) != 0) {
		return nil
	}
	stack.Pop()
	return os.ErrPermission
}

// ChangeDirectoryToPrevious changes the current working directory to the previous one,
// doesn't go past the BasePath
func ChangeDirectoryToPrevious(stack *StringStack) error {
	if stack.IsEmpty() {
		return ErrAlreadyAtBaseDirectory
	}
	stack.Pop()
	return nil
}
