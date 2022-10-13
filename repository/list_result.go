package repository

type ListResult[T any] interface {
	Count() int
	Data() []T
}

func NewListResult[T any](count int, data []T) ListResult[T] {
	return &listResult[T]{count, data}
}

type listResult[T any] struct {
	count int
	data  []T
}

func (r *listResult[T]) Count() int {
	return r.count
}

func (r *listResult[T]) Data() []T {
	return r.data
}
