package main

import (
	appContext "app/context"
	. "app/logger"
	"context"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func HandleEvent(c context.Context, event events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	ctx := appContext.SetRequestId(c)
	Log(ctx).Info("START: authorizer")

	Log(ctx).Debug("Method ARN: " + event.MethodArn)

	tmp := strings.Split(event.MethodArn, ":")
	region := tmp[3]
	awsAccountID := tmp[4]
	apiGatewayArnTmp := strings.Split(tmp[5], "/")

	res := NewAuthorizeResponse("*", awsAccountID)

	res.Region = region
	res.APIID = apiGatewayArnTmp[0]
	res.Stage = apiGatewayArnTmp[1]

	b, err := verifyJWTToken(ctx, event)
	if b {
		res.addMethod(Allow, apiGatewayArnTmp[2], "*")
	} else {
		Log(ctx).Debug(err)
		res.addMethod(Deny, apiGatewayArnTmp[2], "*")
	}

	return res.APIGatewayCustomAuthorizerResponse, nil
}

func verifyJWTToken(ctx context.Context, event events.APIGatewayCustomAuthorizerRequest) (bool, error) {
	cfg, _ := config.LoadDefaultConfig(ctx)
	client := cognitoidentityprovider.NewFromConfig(cfg)

	u, err := client.GetUser(ctx, &cognitoidentityprovider.GetUserInput{
		AccessToken: &event.AuthorizationToken,
	})
	if err != nil {
		return false, err
	}

	Log(ctx).Info(*u.UserAttributes[0].Value) //sub
	Log(ctx).Info(*u.UserAttributes[1].Value) //email

	return true, nil
}

type Effect int

const (
	Allow Effect = iota
	Deny
)

func (e Effect) String() string {
	switch e {
	case Allow:
		return "Allow"
	case Deny:
		return "Deny"
	}
	return ""
}

type AuthorizeResponse struct {
	events.APIGatewayCustomAuthorizerResponse
	Region    string
	AccountID string
	APIID     string
	Stage     string
}

func (r *AuthorizeResponse) addMethod(effect Effect, verb string, resource string) {
	resourceArn := "arn:aws:execute-api:" +
		r.Region + ":" +
		r.AccountID + ":" +
		r.APIID + "/" +
		r.Stage + "/" +
		verb + "/" +
		strings.TrimLeft(resource, "/")

	s := events.IAMPolicyStatement{
		Effect:   effect.String(),
		Action:   []string{"execute-api:Invoke"},
		Resource: []string{resourceArn},
	}

	r.PolicyDocument.Statement = append(r.PolicyDocument.Statement, s)
}

func NewAuthorizeResponse(principalID string, accountID string) *AuthorizeResponse {
	return &AuthorizeResponse{
		APIGatewayCustomAuthorizerResponse: events.APIGatewayCustomAuthorizerResponse{
			PrincipalID: principalID,
			PolicyDocument: events.APIGatewayCustomAuthorizerPolicy{
				Version: "2012-10-17",
			},
		},
		Region:    "*",
		AccountID: accountID,
		APIID:     "*",
		Stage:     "*",
	}
}

func main() {
	lambda.Start(HandleEvent)
}
