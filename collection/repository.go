package collection

import "context"

type GetById[T any] interface {
	GetById(ctx context.Context, id string) (T, error)
}

type GetList[T any] interface {
	GetList(ctx context.Context) (ListResult[T], error)
}

type Create[T any] interface {
	Create(ctx context.Context, data T) (string, error)
}

type Update[T any] interface {
	Update(ctx context.Context, id string, data T) error
}

type Delete interface {
	Delete(ctx context.Context, id string) error
}

type Count interface {
	Count(ctx context.Context) (int, error)
}
