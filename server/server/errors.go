package server

import "errors"

// InputError will be raised when the input is not right.
type InputError struct {
	// The operation that caused the error.
	Op  string
	// The error that occurred during the operation.
	Err error
}

func (e *InputError) Error() string { return "Error: " + e.Op + ": " + e.Err.Error() }

// Input Errors
var (
	UnknownError     = errors.New("Unknown Error.")
	InvalidCommand   = errors.New("Invalid command.")
	TooManyArguments = errors.New("Too many arguments.")
	TooFewArguments  = errors.New("Too few arguments.")
)

// Command Errors
var (
	GetNoBitsError = errors.New("The file/directory contains zero bits!")
)