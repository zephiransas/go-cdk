package main

import (
	appContext "app/context"
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
	Log(ctx).Info(fmt.Sprintf("START: todos/%s/_done", req.PathParameters["id"]))

	var s service.TodoService
	if s, err = service.NewTodoService(ctx); err != nil {
		return
	}

	todo, err := s.Done(ctx, req.RequestContext.Authorizer["sub"].(string), req.PathParameters["id"])
	if err != nil {
		return
	}

	var j []byte
	if j, err = json.Marshal(todo); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 503}, err
	}
	return events.APIGatewayProxyResponse{Body: string(j), StatusCode: 200}, nil

}

func main() {
	lambda.Start(HandleEvent)
}
