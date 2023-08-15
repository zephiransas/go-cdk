package service

import (
	appContext "app/context"
	"app/domain"
	"app/domain/repository"
	"app/infra/dynamodb"
	"context"
)

type TodoService interface {
	List(ctx context.Context) (todos []domain.Todo, err error)
	Add(ctx context.Context, title string) (todo domain.Todo, err error)
	Get(ctx context.Context, id string) (todo domain.Todo, err error)
	Done(ctx context.Context, id string) (todo domain.Todo, err error)
	Delete(ctx context.Context, id string) (todo domain.Todo, err error)
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

func (s *todoService) List(ctx context.Context) (todos []domain.Todo, err error) {
	if todos, err = s.repo.List(ctx, appContext.GetSub(ctx)); err != nil {
		return
	}
	return
}

func (s *todoService) Add(ctx context.Context, title string) (todo domain.Todo, err error) {
	if todo, err = s.repo.Add(ctx, appContext.GetSub(ctx), title); err != nil {
		return
	}
	return
}

func (s *todoService) Get(ctx context.Context, id string) (todo domain.Todo, err error) {
	if todo, err = s.repo.Get(ctx, appContext.GetSub(ctx), id); err != nil {
		return
	}
	return
}

func (s *todoService) Done(ctx context.Context, id string) (todo domain.Todo, err error) {
	if todo, err = s.repo.Done(ctx, appContext.GetSub(ctx), id); err != nil {
		return
	}
	return
}

func (s *todoService) Delete(ctx context.Context, id string) (todo domain.Todo, err error) {
	if todo, err = s.repo.Delete(ctx, appContext.GetSub(ctx), id); err != nil {
		return
	}
	return
}
