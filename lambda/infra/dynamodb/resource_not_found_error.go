package dynamodb

import (
	"errors"
)

type ResourceNotFoundError struct {
	err error
}

func NewResourceNotFoundError() error {
	return &ResourceNotFoundError{
		errors.New("Resource Not Found."),
	}
}

func (e *ResourceNotFoundError) Error() string {
	return e.err.Error()
}

func (e *ResourceNotFoundError) Unwrap() error {
	return e.err
}
