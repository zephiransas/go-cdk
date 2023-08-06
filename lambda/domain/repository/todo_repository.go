package repository

import (
	"app/domain"
	"context"
)

type TodoRepository interface {
	List(ctx context.Context, sub string) (todos []domain.Todo, err error)
	Get(ctx context.Context, sub string, id string) (todo domain.Todo, err error)
	Add(ctx context.Context, sub string, title string) (todo domain.Todo, err error)
	Done(ctx context.Context, sub string, title string) (todo domain.Todo, err error)
}
