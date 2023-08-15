package main

import (
	"app/domain"
	"app/middleware"
	"app/service"
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type addRequestBody struct {
	Title string `json:"title"`
}

func HandleEvent(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {
	var (
		s    service.TodoService
		body addRequestBody
		todo domain.Todo
	)

	if s, err = service.NewTodoService(ctx); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 503}, err
	}

	if err := json.Unmarshal([]byte(req.Body), &body); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 503}, err
	}

	if todo, err = s.Add(ctx, req.RequestContext.Authorizer["sub"].(string), body.Title); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 503}, err
	}

	return middleware.JsonResponse(200, todo)
}

func main() {
	m := middleware.NewMiddleware(middleware.DefaultMiddlewares()...)
	lambda.Start(m.Apply(HandleEvent))
}
