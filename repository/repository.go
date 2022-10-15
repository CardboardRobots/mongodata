package repository

import "context"

type Repository[T any] interface {
	GetById[T]
	GetList[T]
	Create[T]
	Update[T]
	Delete
}

type GetById[T any] interface {
	GetById(ctx context.Context, id string) (T, error)
}

type Query map[string]interface{}

type GetList[T any] interface {
	GetList(ctx context.Context, query Query) (*ListResult[T], error)
}

type Create[T any] interface {
	Create(ctx context.Context, data T) (string, error)
}

type Update[T any] interface {
	Update(ctx context.Context, id string, data T) (bool, error)
}

type Delete interface {
	Delete(ctx context.Context, id string) (bool, error)
}
