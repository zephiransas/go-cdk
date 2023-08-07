package dynamodb

import (
	"encoding/json"
	"errors"
)

type resourceNotFoundErrorBody struct {
	Message string `json:"message"`
}

type ResourceNotFoundError struct {
	err error
}

func NewResourceNotFoundError() error {
	return &ResourceNotFoundError{
		errors.New("Resource Not Found."),
	}
}

func (e *ResourceNotFoundError) Error() string {
	j, err := json.Marshal(resourceNotFoundErrorBody{e.err.Error()})
	if err != nil {
		panic(err)
	}
	return string(j)
}

func (e *ResourceNotFoundError) Unwrap() error {
	return e.err
}
