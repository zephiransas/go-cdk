package main

import (
	. "app/logger"
	"app/middleware"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"os"
)

var poolId = os.Getenv("POOL_ID")
var clientId = os.Getenv("CLIENT_ID")

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int32  `json:"expires_in"`
	IdToken      string `json:"id_token"`
}

func HandleEvent(ctx context.Context, req events.APIGatewayProxyRequest) (res events.APIGatewayProxyResponse, err error) {

	cfg, _ := config.LoadDefaultConfig(ctx)
	client := cognitoidentityprovider.NewFromConfig(cfg)

	username := req.QueryStringParameters["username"]
	password := req.QueryStringParameters["password"]

	var (
		clientSecret string
		out          *cognitoidentityprovider.AdminInitiateAuthOutput
	)

	if clientSecret, err = getClientSecret(&ctx, &cfg); err != nil {
		Log(ctx).Error(err)
		return events.APIGatewayProxyResponse{StatusCode: 503}, err
	}

	hash := makeHMAC(username, clientSecret)

	if out, err = client.AdminInitiateAuth(ctx, &cognitoidentityprovider.AdminInitiateAuthInput{
		UserPoolId: &poolId,
		ClientId:   &clientId,
		AuthFlow:   types.AuthFlowTypeAdminUserPasswordAuth,
		AuthParameters: map[string]string{
			"USERNAME":    username,
			"PASSWORD":    password,
			"SECRET_HASH": hash,
		},
	}); err != nil {
		// if notAuth, ok := err.(*types.NotAuthorizedException); ok {... でハンドルできないか？
		var notAuth *types.NotAuthorizedException
		if errors.As(err, &notAuth) {
			Log(ctx).Info(err)
			return events.APIGatewayProxyResponse{StatusCode: 401}, nil
		}

		Log(ctx).Error(err)
		return events.APIGatewayProxyResponse{StatusCode: 503}, err
	}

	r := LoginResponse{
		AccessToken:  *out.AuthenticationResult.AccessToken,
		RefreshToken: *out.AuthenticationResult.RefreshToken,
		ExpiresIn:    out.AuthenticationResult.ExpiresIn,
		IdToken:      *out.AuthenticationResult.IdToken,
	}

	var j []byte
	if j, err = json.Marshal(r); err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 503}, err
	}

	return events.APIGatewayProxyResponse{Body: string(j), StatusCode: 200}, nil
}

func makeHMAC(username string, clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func getClientSecret(ctx *context.Context, cfg *aws.Config) (v string, err error) {
	s := ssm.NewFromConfig(*cfg)
	var clientSecret *ssm.GetParameterOutput
	if clientSecret, err = s.GetParameter(*ctx, &ssm.GetParameterInput{
		Name:           aws.String("/go-cdk/clientSecret"),
		WithDecryption: aws.Bool(true),
	}); err != nil {
		return
	}
	return *clientSecret.Parameter.Value, nil
}

func main() {
	m := middleware.NewMiddleware(middleware.DefaultMiddlewares()...)
	lambda.Start(m.Apply(HandleEvent))
}
