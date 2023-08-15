package repository

import (
	"app/domain/vo"
	"context"
)

type CounterRepository interface {
	GenerateId(ctx context.Context, sub vo.SubId) (id int, err error)
}
