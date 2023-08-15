package middleware

import (
	appContext "app/context"
	. "app/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
)

type LambdaHandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

type LambdaMiddlewareFunc func(next LambdaHandlerFunc) LambdaHandlerFunc

type Middleware struct {
	middlewares []LambdaMiddlewareFunc
}

func (m *Middleware) Use(middlewares ...LambdaMiddlewareFunc) {
	m.middlewares = append(m.middlewares, middlewares...)
}

func (m *Middleware) Apply(handler LambdaHandlerFunc) LambdaHandlerFunc {
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		handler = m.middlewares[i](handler)
	}
	return handler
}

func DefaultMiddlewares() []LambdaMiddlewareFunc {
	return []LambdaMiddlewareFunc{
		RequestId(),
		Logging(),
	}
}

func NewMiddleware(middlewares ...LambdaMiddlewareFunc) *Middleware {
	return &Middleware{middlewares: middlewares}
}

func RequestId() LambdaMiddlewareFunc {
	return func(next LambdaHandlerFunc) LambdaHandlerFunc {
		return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			c := appContext.SetRequestId(ctx)
			res, err := next(c, request)
			return res, err
		}
	}
}

func Logging() LambdaMiddlewareFunc {
	return func(next LambdaHandlerFunc) LambdaHandlerFunc {
		return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			res, err := next(ctx, request)
			Log(ctx).Info(fmt.Sprintf("req: %+v, res: %+v", request, res))
			return res, err
		}
	}
}

func JsonResponse(status int, body any) (events.APIGatewayProxyResponse, error) {
	j, err := json.Marshal(body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}
	return events.APIGatewayProxyResponse{StatusCode: status, Body: string(j)}, nil
}
