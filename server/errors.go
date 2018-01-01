package server

import "errors"

// InputError will be raised when the input given to the parser by the client is not right.
type InputError struct {
	// The operation that caused the error.
	Op string
	// The error that occurred during the operation.
	Err error
}

func (e InputError) Error() string { return "Error: " + e.Op + ": " + e.Err.Error() }

// Input Errors
var (
	InputInvalidCommand   = errors.New("invalid command")
	InputTooManyArguments = errors.New("too many arguments")
	InputTooFewArguments  = errors.New("too few arguments")
)

// Command Errors represent errors that occur when the server is executing commands
var (
	GetNoBitsError = errors.New("the file/directory contains zero bits")
)

type StackError struct {
	ErrorName string
	Err       error
}

func (e StackError) Error() string { return e.ErrorName + ": " + e.Err.Error() }

// Stack Errors
var (
	StackInvalidTypeError = StackError{"InvalidTypeError", errors.New("invalid item type for the Stack")}
	StackOverflowError    = StackError{"StackOverflowError", errors.New("stack capacity exceeded")}
	StackUnderflowError   = StackError{"StackUnderflowError", errors.New("stack is empty")}
	ErrStackCast          = StackError{"StackCastError", errors.New("stack can't be casted to selected type")}
)

type PathError struct {
	Err error
}

func (e PathError) Error() string { return "Error: " + e.Err.Error() }

// PathErrors
var (
	ErrInvalidDirectoryName   = PathError{errors.New("names should not contain / character")}
	ErrNotADirectory          = PathError{errors.New("file name is not a valid directory")}
	ErrAlreadyAtBaseDirectory = PathError{errors.New("can't go past beyond root directory")}
	ErrSlashNotAllowed        = PathError{errors.New("slash is not allowed in file names")}
)

// General Errors
var (
	ErrUploadServerFailure = errors.New("upload server failed to start")
)
