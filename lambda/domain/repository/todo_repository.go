package repository

import (
	"app/domain"
	"context"
)

type TodoRepository interface {
	List(ctx context.Context) (todos []domain.Todo, err error)

	Add(ctx context.Context, title string) (todo domain.Todo, err error)
}
