package dynamodb

import (
	"app/domain/repository"
	"app/domain/vo"
	"app/util"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type counterRepository struct {
	client *dynamodb.Client
}

func NewCounterRepository(ctx context.Context) (r repository.CounterRepository, err error) {
	cfg := NewConfig()
	awsCfg, err := util.LoadDefaultConfig(ctx, cfg.region, cfg.endpoint, cfg.service)
	if err != nil {
		return
	}
	r = &counterRepository{
		client: dynamodb.NewFromConfig(awsCfg),
	}
	return
}

func (r *counterRepository) GenerateId(ctx context.Context, sub vo.SubId) (id int, err error) {
	var res *dynamodb.UpdateItemOutput

	res, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("todos-go-counter"),
		Key: map[string]types.AttributeValue{
			"user_id": toS(string(sub)),
		},
		UpdateExpression: aws.String("SET #value = if_not_exists(#value, :start) + :incr"),
		ExpressionAttributeNames: map[string]string{
			"#value": "id",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":incr":  toN("1"),
			":start": toN("0"),
		},
		ReturnValues: types.ReturnValueUpdatedNew,
	})
	if err != nil {
		return
	}

	c := response{}
	if err = attributevalue.UnmarshalMap(res.Attributes, &c); err != nil {
		return
	}

	return c.Id, nil
}

type response struct {
	Id int `dynamodbav:"id"`
}
