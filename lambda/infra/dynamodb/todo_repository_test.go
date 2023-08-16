package dynamodb

import (
	"app/domain/repository"
	"app/domain/vo"
	"app/testutil"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupTestTodoRepository(t *testing.T) (repository.TodoRepository, context.Context) {
	testutil.SetLocalMockEnv(t)
	ctx := context.TODO()
	r, err := NewTodoRepository(ctx)
	assert.NoError(t, err)
	return r, ctx
}

func TestTodoRepository_List(t *testing.T) {
	r, ctx := setupTestTodoRepository(t)
	CleanTable(t, ctx, "todos-go")

	PutItem(t, ctx, "todos-go", map[string]types.AttributeValue{
		"user_id": toS("sub"),
		"id":      toN("1"),
		"title":   toS("todo1"),
		"done":    toBOOL(false),
	})
	PutItem(t, ctx, "todos-go", map[string]types.AttributeValue{
		"user_id": toS("sub"),
		"id":      toN("2"),
		"title":   toS("todo2"),
		"done":    toBOOL(false),
	})

	// 他ユーザのデータ
	PutItem(t, ctx, "todos-go", map[string]types.AttributeValue{
		"user_id": toS("DUMMY"),
		"id":      toN("2"),
		"title":   toS("todo2"),
		"done":    toBOOL(false),
	})

	res, err := r.List(ctx, "sub")
	assert.NoError(t, err)

	// 正しく自分のデータのみ取得できていること
	assert.Len(t, res, 2)

	var subs []string
	for _, v := range res {
		subs = append(subs, string(v.UserId))
	}

	// 他ユーザのデータが存在しないこと
	assert.NotContains(t, subs, "DUMMY")
}

func TestTodoRepository_Add(t *testing.T) {
	r, ctx := setupTestTodoRepository(t)
	CleanTable(t, ctx, "todos-go")

	PutItem(t, ctx, "todos-go", map[string]types.AttributeValue{
		"user_id": toS("sub"),
		"id":      toN("1"),
		"title":   toS("todo1"),
		"done":    toBOOL(false),
	})
	PutItem(t, ctx, "todos-go-counter", map[string]types.AttributeValue{
		"user_id": toS("sub"),
		"id":      toN("1"),
	})

	todo, err := r.Add(ctx, vo.NewSubId("sub"), "test_title")
	assert.NoError(t, err)

	assert.Equal(t, vo.NewSubId("sub"), todo.UserId)
	assert.Equal(t, 2, todo.Id)
	assert.Equal(t, "test_title", todo.Title)
	assert.Equal(t, false, todo.Done)
}

func TestTodoRepository_Get(t *testing.T) {
	r, ctx := setupTestTodoRepository(t)
	CleanTable(t, ctx, "todos-go")

	PutItem(t, ctx, "todos-go", map[string]types.AttributeValue{
		"user_id": toS("sub"),
		"id":      toN("1"),
		"title":   toS("todo1"),
		"done":    toBOOL(false),
	})

	todo, err := r.Get(ctx, vo.NewSubId("sub"), "1")
	assert.NoError(t, err)

	assert.Equal(t, vo.NewSubId("sub"), todo.UserId)
	assert.Equal(t, 1, todo.Id)
	assert.Equal(t, "todo1", todo.Title)
	assert.Equal(t, false, todo.Done)
}

func TestTodoRepository_Done(t *testing.T) {
	r, ctx := setupTestTodoRepository(t)
	CleanTable(t, ctx, "todos-go")

	PutItem(t, ctx, "todos-go", map[string]types.AttributeValue{
		"user_id": toS("sub"),
		"id":      toN("1"),
		"title":   toS("todo1"),
		"done":    toBOOL(false),
	})

	todo, err := r.Done(ctx, vo.NewSubId("sub"), "1")
	assert.NoError(t, err)

	assert.Equal(t, true, todo.Done)
}