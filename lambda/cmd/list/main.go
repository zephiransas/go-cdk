package main

import (
	"app/domain"
	"app/service"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleEvent(c context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	var s service.TodoService
	if s, err = service.NewTodoService(c); err != nil {
		return
	}

	var todos []domain.Todo
	if todos, err = s.List(c); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 503}, err
	}

	var j []byte
	if j, err = json.Marshal(domain.Todos{Todos: todos}); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 503}, err
	}
	return events.APIGatewayProxyResponse{Body: string(j), StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleEvent)
}
