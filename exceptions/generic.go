package exceptions

import (
	"errors"
)

// -----------------------------------------------------------------------------
// Not Implemented Error (for development purposes)

type NotImplementedError struct {
	ErrorCode string
	Err       error
}

func (e *NotImplementedError) Error() string {
	return e.ErrorCode + ": " + e.Err.Error()
}

func NewNotImplementedError() *NotImplementedError {
	return &NotImplementedError{
		ErrorCode: "GE001",
		Err:       errors.New("Not implemented"),
	}
}
