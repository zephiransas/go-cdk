package dynamodb

import (
	"app/util"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
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

func PutItem(t *testing.T, ctx context.Context, tableName string, item map[string]types.AttributeValue) {
	c := createClient(t, ctx)
	_, err := c.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	})
	assert.NoError(t, err)
	return
}

func CleanTable(t *testing.T, ctx context.Context, tableName string) {
	c := createClient(t, ctx)

	desc, err := c.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	assert.NoError(t, err)

	out, err := c.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	assert.NoError(t, err)

	for _, i := range out.Items {
		key := make(map[string]types.AttributeValue)

		for _, k := range desc.Table.KeySchema {
			key[*k.AttributeName] = i[*k.AttributeName]
		}
		_, err := c.DeleteItem(ctx, &dynamodb.DeleteItemInput{
			Key:       key,
			TableName: aws.String(tableName),
		})
		assert.NoError(t, err)
	}
}
