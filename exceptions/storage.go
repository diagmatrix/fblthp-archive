package exceptions

import (
	"errors"
	"strconv"
)

// -----------------------------------------------------------------------------
// Card Not Found Error

type CardNotFoundError struct {
	ErrorCode string
	CardID    int
	Err       error
}

func (e *CardNotFoundError) Error() string {
	return e.ErrorCode + ": Card with ID " + strconv.Itoa(e.CardID) + " not found"
}

func NewCardNotFoundError(id int) *CardNotFoundError {
	return &CardNotFoundError{
		ErrorCode: "ST001",
		CardID:    id,
		Err:       errors.New("Card not found"),
	}
}
