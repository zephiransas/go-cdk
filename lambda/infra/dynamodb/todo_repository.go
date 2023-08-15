package dynamodb

import (
	"app/domain"
	"app/domain/repository"
	"app/domain/vo"
	. "app/logger"
	"app/util"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type todoRepository struct {
	client  *dynamodb.Client
	counter repository.CounterRepository
}

func NewTodoRepository(ctx context.Context) (r repository.TodoRepository, err error) {
	c := NewConfig()
	cfg, _ := util.LoadDefaultConfig(ctx, c.region, c.endpoint, c.service)
	cr, err := NewCounterRepository(ctx)
	if err != nil {
		return
	}
	r = &todoRepository{
		client:  dynamodb.NewFromConfig(cfg),
		counter: cr,
	}
	return
}

func (r *todoRepository) List(ctx context.Context, sub vo.SubId) (todos []domain.Todo, err error) {
	var out *dynamodb.QueryOutput

	// ref: https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/gov2/dynamodb
	ex := expression.Key("user_id").Equal(expression.Value(string(sub)))
	expr, _ := expression.NewBuilder().WithKeyCondition(ex).Build()

	var input *dynamodb.QueryInput
	input = &dynamodb.QueryInput{
		TableName:                 aws.String("todos-go"),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}

	if out, err = r.client.Query(ctx, input); err != nil {
		return
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &todos)

	return
}

func (r *todoRepository) Add(ctx context.Context, sub vo.SubId, title string) (todo domain.Todo, err error) {
	var id int
	if id, err = r.counter.GenerateId(ctx, sub); err != nil {
		return
	}

	item := map[string]types.AttributeValue{
		"user_id": toS(string(sub)),
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

func (r *todoRepository) Get(ctx context.Context, sub vo.SubId, id string) (todo domain.Todo, err error) {
	var out *dynamodb.GetItemOutput

	Log(ctx).Debug(fmt.Sprintf("sub = %s, id = %s", sub, id))

	out, err = r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("todos-go"),
		Key: map[string]types.AttributeValue{
			"user_id": toS(string(sub)),
			"id":      toN(id),
		},
	})
	if err != nil {
		return
	}

	if out.Item == nil {
		err = NewResourceNotFoundError()
		return
	}

	if err = attributevalue.UnmarshalMap(out.Item, &todo); err != nil {
		return
	}
	return
}

func (r *todoRepository) Done(ctx context.Context, sub vo.SubId, id string) (todo domain.Todo, err error) {
	var res *dynamodb.UpdateItemOutput

	res, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String("todos-go"),
		Key: map[string]types.AttributeValue{
			"user_id": toS(string(sub)),
			"id":      toN(id),
		},
		UpdateExpression: aws.String("SET #value = :done"),
		ExpressionAttributeNames: map[string]string{
			"#value": "done",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":done": toBOOL(true),
		},
		ReturnValues: types.ReturnValueAllNew,
	})
	if err != nil {
		return
	}
	if err = attributevalue.UnmarshalMap(res.Attributes, &todo); err != nil {
		return
	}
	return
}

func (r *todoRepository) Delete(ctx context.Context, sub vo.SubId, id string) (todo domain.Todo, err error) {
	var res *dynamodb.DeleteItemOutput

	res, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String("todos-go"),
		Key: map[string]types.AttributeValue{
			"user_id": toS(string(sub)),
			"id":      toN(id),
		},
		ReturnValues: types.ReturnValueAllOld,
	})
	if err != nil {
		return
	}

	if res.Attributes == nil {
		err = NewResourceNotFoundError()
		return
	}

	if err = attributevalue.UnmarshalMap(res.Attributes, &todo); err != nil {
		return
	}
	return
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
