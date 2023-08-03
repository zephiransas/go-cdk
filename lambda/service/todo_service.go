package service

import (
	"app/domain"
	"app/domain/repository"
	"app/infra/dynamodb"
	"context"
)

type TodoService interface {
	List(ctx context.Context, sub string) (todos []domain.Todo, err error)
	Add(ctx context.Context, sub string, title string) (todo domain.Todo, err error)
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(ctx context.Context) (s TodoService, err error) {
	var repo repository.TodoRepository
	if repo, err = dynamodb.NewTodoRepository(ctx); err != nil {
		return
	}
	return &todoService{repo}, nil
}

func (s *todoService) List(ctx context.Context, sub string) (todos []domain.Todo, err error) {
	if todos, err = s.repo.List(ctx, sub); err != nil {
		return
	}
	return
}

func (s *todoService) Add(ctx context.Context, sub string, title string) (todo domain.Todo, err error) {
	if todo, err = s.repo.Add(ctx, sub, title); err != nil {
		return
	}
	return
}
