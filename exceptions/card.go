package exceptions

import (
	"errors"
)

// -----------------------------------------------------------------------------
// Cannot Parse Typeline Error

type CannotParseTypelineError struct {
	ErrorCode string
	Typeline  string
	Reason    string
	Err       error
}

func (e *CannotParseTypelineError) Error() string {
	return e.ErrorCode + ": Cannot parse typeline " + e.Typeline + "(" + e.Reason + ")"
}

func NewCannotParseTypelineError(typeline string, err string) *CannotParseTypelineError {
	return &CannotParseTypelineError{
		ErrorCode: "CA001",
		Typeline:  typeline,
		Reason:    err,
		Err:       errors.New("Cannot parse typeline"),
	}
}

// -----------------------------------------------------------------------------
// No Card Name Error

type NoCardNameError struct {
	ErrorCode string
	Err       error
}

func (e *NoCardNameError) Error() string {
	return e.ErrorCode + ": " + e.Err.Error()
}

func NewNoCardNameError() *NoCardNameError {
	return &NoCardNameError{
		ErrorCode: "CA002",
		Err:       errors.New("No card name provided"),
	}
}

// -----------------------------------------------------------------------------
// No Set ID Error

type NoSetIDError struct {
	ErrorCode string
	Err       error
}

func (e *NoSetIDError) Error() string {
	return e.ErrorCode + ": " + e.Err.Error()
}

func NewNoSetIDError() *NoSetIDError {
	return &NoSetIDError{
		ErrorCode: "CA003",
		Err:       errors.New("No set ID provided"),
	}
}

// -----------------------------------------------------------------------------
// No Collector Number Error

type NoCollectorNumberError struct {
	ErrorCode string
	Err       error
}

func (e *NoCollectorNumberError) Error() string {
	return e.ErrorCode + ": " + e.Err.Error()
}

func NewNoCollectorNumberError() *NoCollectorNumberError {
	return &NoCollectorNumberError{
		ErrorCode: "CA004",
		Err:       errors.New("No collector number provided"),
	}
}
