package main

import (
	"app/domain"
	"app/middleware"
	"app/service"
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleEvent(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {

	var s service.TodoService
	if s, err = service.NewTodoService(ctx); err != nil {
		return
	}

	var todos []domain.Todo
	if todos, err = s.List(ctx, req.RequestContext.Authorizer["sub"].(string)); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 503}, err
	}

	return middleware.JsonResponse(200, domain.Todos{Todos: todos})
}

func main() {
	m := middleware.NewMiddleware(middleware.DefaultMiddlewares()...)
	lambda.Start(m.Apply(HandleEvent))
}
