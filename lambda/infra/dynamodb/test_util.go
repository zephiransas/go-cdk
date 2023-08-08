package dynamodb

import (
	"app/domain"
	"app/util"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type Config struct {
	region   string
	endpoint string
	service  string
}

func NewConfig() *Config {
	return &Config{
		region:   os.Getenv("AWS_DEFAULT_REGION"),
		endpoint: os.Getenv("DDB_ENDPOINT"),
		service:  "DynamoDB",
	}
}

func createClient(t *testing.T, ctx context.Context) *dynamodb.Client {
	cfg := NewConfig()
	awsCfg, err := util.LoadDefaultConfig(ctx, cfg.region, cfg.endpoint, cfg.service)
	assert.NoError(t, err)
	return dynamodb.NewFromConfig(awsCfg)
}

func PutItem(t *testing.T, ctx context.Context, todo domain.Todo) {
	c := createClient(t, ctx)
	var av map[string]types.AttributeValue

	av, err := attributevalue.MarshalMap(todo)
	assert.NoError(t, err)

	_, err = c.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("todos-go"),
	})
	assert.NoError(t, err)
	return
}
