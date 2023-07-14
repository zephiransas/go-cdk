package main

import (
	"app/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleEvent(c context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var todos []domain.Todo
	for i := 0; i < 5; i++ {
		todos = append(todos, domain.Todo{Id: fmt.Sprintf("id-%d", i), Title: fmt.Sprintf("todo %d", i), Done: false})
	}
	j, _ := json.Marshal(domain.Todos{Todos: todos})
	return events.APIGatewayProxyResponse{Body: string(j), StatusCode: 200}, nil
}

func main() {
	lambda.Start(HandleEvent)
}
