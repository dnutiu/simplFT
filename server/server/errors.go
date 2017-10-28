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
	InputUnknownError     = errors.New("Unknown Error.")
	InputInvalidCommand   = errors.New("Invalid command.")
	InputTooManyArguments = errors.New("Too many arguments.")
	InputTooFewArguments  = errors.New("Too few arguments.")
)

// Command Errors represent errors that occur when the server is executing commands
var (
	GetNoBitsError = errors.New("The file/directory contains zero bits!")
)

type StackError struct {
	ErrorName string
	Err       error
}

func (e StackError) Error() string { return e.ErrorName + ": " + e.Err.Error() }

// Stack Errors
var (
	StackInvalidTypeError = StackError{"InvalidTypeError", errors.New("Invalid item type for the Stack")}
	StackOverflowError    = StackError{"StackOverflowError", errors.New("Stack capacity exceeded!")}
	StackUnderflowError   = StackError{"StackUnderflowError", errors.New("Stack is empty!")}
)
