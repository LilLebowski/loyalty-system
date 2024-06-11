package utils

import (
	"fmt"
)

type OrderIsExistAnotherUserError struct {
	Message string
	Err     error
}

func NewOrderIsExistAnotherUserError(msg string, err error) error {
	return &OrderIsExistAnotherUserError{
		Message: msg,
		Err:     err,
	}
}

func (e *OrderIsExistAnotherUserError) Error() string {
	return fmt.Sprintf("[%s] %v", e.Message, e.Err)
}

type OrderIsExistThisUserError struct {
	Message string
	Err     error
}

func NewOrderIsExistThisUserError(msg string, err error) error {
	return &OrderIsExistThisUserError{
		Message: msg,
		Err:     err,
	}
}

func (e *OrderIsExistThisUserError) Error() string {
	return fmt.Sprintf("[%s] %v", e.Message, e.Err)
}

type LessBonusError struct {
	Message string
	Err     error
}

func NewLessBonusErrorError(msg string, err error) error {
	return &LessBonusError{
		Message: msg,
		Err:     err,
	}
}

func (e *LessBonusError) Error() string {
	return fmt.Sprintf("[%s] %v", e.Message, e.Err)
}

var ErrTooManyRequests = fmt.Errorf("too many requests")
var ErrNoContent = fmt.Errorf("no content")
