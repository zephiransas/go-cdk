package repository

import (
	"app/domain"
	"app/domain/vo"
	"context"
)

type TodoRepository interface {
	List(ctx context.Context, sub vo.SubId) (todos []domain.Todo, err error)
	Get(ctx context.Context, sub vo.SubId, id string) (todo domain.Todo, err error)
	Add(ctx context.Context, sub vo.SubId, title string) (todo domain.Todo, err error)
	Done(ctx context.Context, sub vo.SubId, id string) (todo domain.Todo, err error)
	Delete(ctx context.Context, sub vo.SubId, id string) (todo domain.Todo, err error)
}
