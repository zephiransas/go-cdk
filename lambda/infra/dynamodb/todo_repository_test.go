package dynamodb

import (
	"app/domain"
	"app/domain/repository"
	"app/testutil"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup(t *testing.T) (repository.TodoRepository, context.Context) {
	testutil.SetLocalMockEnv(t)
	ctx := context.TODO()
	r, err := NewTodoRepository(ctx)
	assert.NoError(t, err)
	return r, ctx
}

func TestTodoRepository_List(t *testing.T) {
	r, ctx := setup(t)

	PutItem(t, ctx, domain.Todo{
		UserId: "sub",
		Id:     1,
		Title:  "todo1",
		Done:   false,
	})
	PutItem(t, ctx, domain.Todo{
		UserId: "sub",
		Id:     2,
		Title:  "todo2",
		Done:   false,
	})

	// 他ユーザのデータ
	PutItem(t, ctx, domain.Todo{
		UserId: "DUMMY",
		Id:     1,
		Title:  "may_not_exists",
		Done:   false,
	})

	res, err := r.List(ctx, "sub")
	assert.NoError(t, err)

	// 正しく自分のデータのみ取得できていること
	assert.Len(t, res, 2)

	var subs []string
	for _, v := range res {
		subs = append(subs, v.UserId)
	}

	// 他ユーザのデータが存在しないこと
	assert.NotContains(t, subs, "DUMMY")

}
