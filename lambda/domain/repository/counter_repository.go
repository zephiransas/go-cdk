package repository

import "context"

type CounterRepository interface {
	GenerateId(ctx context.Context, sub string) (id int, err error)
}
