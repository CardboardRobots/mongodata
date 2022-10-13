package repository

import "context"

type Repository[T any] interface {
	GetById[T]
	GetList[T]
	Insert[T]
	Update[T]
	Delete
	Count
}

type GetById[T any] interface {
	GetById(ctx context.Context, id string) (T, error)
}

type GetList[T any] interface {
	GetList(ctx context.Context) (ListResult[T], error)
}

type Insert[T any] interface {
	Insert(ctx context.Context, data T) (string, error)
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
