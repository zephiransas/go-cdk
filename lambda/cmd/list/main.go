package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

type MyResponse struct {
	Message string `json:"message"`
}

func HandleEvent(c context.Context, event MyEvent) (string, error) {
	r := MyResponse{
		Message: fmt.Sprintf("Hello %s!", event.Name),
	}
	j, _ := json.Marshal(r)
	return string(j), nil
}

func main() {
	lambda.Start(HandleEvent)
}
