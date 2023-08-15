package main

import (
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

	todo, err := s.Done(ctx, req.RequestContext.Authorizer["sub"].(string), req.PathParameters["id"])
	if err != nil {
		return
	}

	return middleware.JsonResponse(200, todo)
}

func main() {
	m := middleware.NewMiddleware(middleware.DefaultMiddlewares()...)
	lambda.Start(m.Apply(HandleEvent))
}
