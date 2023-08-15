package main

import (
	"app/infra/dynamodb"
	"app/middleware"
	"app/service"
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleEvent(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {

	var s service.TodoService
	if s, err = service.NewTodoService(ctx); err != nil {
		return
	}

	todo, err := s.Delete(ctx, req.RequestContext.Authorizer["sub"].(string), req.PathParameters["id"])
	if err != nil {
		var notFound *dynamodb.ResourceNotFoundError
		if errors.As(err, &notFound) {
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       err.Error(),
			}, nil
		} else {
			return
		}
	}

	return middleware.JsonResponse(200, todo)
}

func main() {
	m := middleware.NewMiddleware(middleware.DefaultMiddlewares()...)
	lambda.Start(m.Apply(HandleEvent))
}
