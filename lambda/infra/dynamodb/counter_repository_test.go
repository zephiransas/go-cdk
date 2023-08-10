package dynamodb

import (
	"app/domain/repository"
	"app/testutil"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setupTestCounterRepository(t *testing.T) (repository.CounterRepository, context.Context) {
	testutil.SetLocalMockEnv(t)
	ctx := context.TODO()
	r, err := NewCounterRepository(ctx)
	assert.NoError(t, err)
	return r, ctx
}

func TestCounterRepository_GenerateId(t *testing.T) {
	r, ctx := setupTestCounterRepository(t)
	CleanTable(t, ctx, "todos-go-counter")

	PutItem(t, ctx, "todos-go-counter", map[string]types.AttributeValue{
		"user_id": toS("test"),
		"id":      toN("10"),
	})

	id, err := r.GenerateId(ctx, "test")
	assert.NoError(t, err)

	// 10 -> 11にインクリメントされていること
	assert.Equal(t, 11, id)

	id, err = r.GenerateId(ctx, "other")
	assert.NoError(t, err)

	// 存在しない場合1から始まること
	assert.Equal(t, 1, id)
}
