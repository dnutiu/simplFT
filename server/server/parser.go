package server

import (
	"net"
	"strings"
)

// checkArgumentsLength returns an error if length is not equal to expected.
func checkArgumentsLength(length int, expected int) error {
	if length > expected {
		return InputTooManyArguments
	} else if length < expected {
		return InputTooFewArguments
	}
	return nil
}

func ProcessInput(c net.Conn, text string) error {
	commands := strings.Fields(text)
	commandsLen := len(commands)

	if commandsLen == 0 {
		return nil
	}

	thisCommand := commands[0]

	switch thisCommand {
	case "get":
		err := checkArgumentsLength(commandsLen, 2)
		if err != nil {
			return &InputError{thisCommand, err}
		}

		// Get the file
		_, err = GetFile(c, commands[1])
		if err != nil {
			return &InputError{thisCommand, err}
		}
	case "ls":
		err := checkArgumentsLength(commandsLen, 1)
		if err != nil {
			return &InputError{thisCommand, err}
		}

		err = ListFiles(c)
		if err != nil {
			return &InputError{thisCommand, err}
		}
	case "clear":
		err := checkArgumentsLength(commandsLen, 1)
		if err != nil {
			return &InputError{thisCommand, err}
		}

		// Ansi clear: 1b 5b 48 1b 5b 4a
		// clear | hexdump -C
		var b = []byte{0x1b, 0x5b, 0x48, 0x1b, 0x5b, 0x4a}

		c.Write(b)
	case "help":
		// Check arguments
		err := checkArgumentsLength(commandsLen, 1)
		if err != nil {
			return &InputError{thisCommand, err}
		}

		err = ShowHelp(c)
		if err != nil {
			return &InputError{thisCommand, err}
		}
	case "exit":
		err := checkArgumentsLength(commandsLen, 1)
		if err != nil {
			return &InputError{thisCommand, err}
		}

		c.Close()
	default:
		return &InputError{thisCommand, InputInvalidCommand}
	}

	return nil
}
