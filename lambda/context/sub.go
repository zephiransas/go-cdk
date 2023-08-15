package context

import (
	"app/domain/vo"
	"context"
	"github.com/aws/aws-lambda-go/events"
)

const subKey = "SUB_ID"

func SetSub(c context.Context, req events.APIGatewayProxyRequest) context.Context {
	sub, ok := req.RequestContext.Authorizer["sub"].(string)
	if !ok {
		return c
	}
	return context.WithValue(c, subKey, sub)
}

func GetSub(c context.Context) vo.SubId {
	v, ok := c.Value(subKey).(string)
	if !ok {
		return "Not Authorized"
	}
	return vo.NewSubId(v)
}
