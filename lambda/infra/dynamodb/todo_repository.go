package dynamodb

import (
	"app/domain"
	"app/domain/repository"
	"context"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type todoRepository struct {
	client *dynamodb.Client
}

func NewTodoRepository(ctx context.Context) (r repository.TodoRepository, err error) {
	var cfg aws.Config
	cfg, _ = config.LoadDefaultConfig(ctx, func(o *config.LoadOptions) error {
		o.Region = "ap-northeast-1"
		return nil
	})
	r = &todoRepository{
		client: dynamodb.NewFromConfig(cfg),
	}
	return
}

func (r *todoRepository) List(ctx context.Context, sub string) (todos []domain.Todo, err error) {
	var out *dynamodb.QueryOutput

	// TODO: Replace github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression
	expression := "user_id = :user_id"
	expValue := map[string]types.AttributeValue{":user_id": toS(sub)}

	var input *dynamodb.QueryInput
	input = &dynamodb.QueryInput{
		TableName:                 aws.String("todos-go"),
		KeyConditionExpression:    &expression,
		ExpressionAttributeValues: expValue,
	}

	if out, err = r.client.Query(ctx, input); err != nil {
		return
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &todos)

	return
}

func (r *todoRepository) Add(ctx context.Context, sub string, title string) (todo domain.Todo, err error) {
	var id int
	if id, err = r.generateId(ctx, sub); err != nil {
		return
	}

	item := map[string]types.AttributeValue{
		"user_id": toS(sub),
		"id":      toN(strconv.Itoa(id)),
		"title":   toS(title),
		"done":    toBOOL(false),
	}

	if _, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("todos-go"),
	}); err != nil {
		return
	}

	if err = attributevalue.UnmarshalMap(item, &todo); err != nil {
		return
	}
	return
}

func (r *todoRepository) generateId(ctx context.Context, sub string) (id int, err error) {
	var res *dynamodb.UpdateItemOutput

	res, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("todos-go-counter"),
		Key: map[string]types.AttributeValue{
			"user_id": toS(sub),
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
		return 0, err
	}

	c := Response{}
	if err = attributevalue.UnmarshalMap(res.Attributes, &c); err != nil {
		return 0, err
	}

	return c.Id, nil
}

type Response struct {
	Id int `dynamodbav:"id"`
}

func toS(v string) *types.AttributeValueMemberS {
	return &types.AttributeValueMemberS{Value: v}
}

func toN(v string) *types.AttributeValueMemberN {
	return &types.AttributeValueMemberN{Value: v}
}

func toBOOL(v bool) *types.AttributeValueMemberBOOL {
	return &types.AttributeValueMemberBOOL{Value: v}
}
