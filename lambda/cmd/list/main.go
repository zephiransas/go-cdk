package main

import (
	appContext "app/context"
	"app/domain"
	. "app/logger"
	"app/service"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleEvent(c context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	ctx := appContext.SetRequestId(c)
	Log(ctx).Info("START: todos/list")

	Log(ctx).Debug(fmt.Sprintf("sub = %s", req.RequestContext.Authorizer["sub"]))

	var s service.TodoService
	if s, err = service.NewTodoService(ctx); err != nil {
		return
	}

	var todos []domain.Todo
	if todos, err = s.List(ctx); err != nil {
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
