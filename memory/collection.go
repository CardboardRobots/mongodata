package collection

import (
	"context"
	"fmt"

	"github.com/cardboardrobots/mongodata/repository"
)

type Collection[T any] struct {
	Items             map[string]T
	currentId         int
	preInsertCallback func(value T, id string)
}

func NewCollection[T any](preInsertCallback func(value T, id string)) repository.Repository[T] {
	return &Collection[T]{
		Items:             map[string]T{},
		preInsertCallback: preInsertCallback,
	}
}

func (c *Collection[T]) GetList(ctx context.Context) (repository.ListResult[T], error) {
	count, _ := c.Count(ctx)
	return repository.NewListResult(count, c.Slice()), nil
}

func (c *Collection[T]) Count(ctx context.Context) (int, error) {
	return len(c.Items), nil
}

func (c *Collection[T]) Slice() []T {
	data := make([]T, len(c.Items))
	i := 0
	for _, value := range c.Items {
		data[i] = value
		i++
	}
	return data
}

func (c *Collection[T]) GetById(ctx context.Context, id string) (T, error) {
	value, ok := c.Items[id]
	if !ok {
		var zero T
		return zero, repository.NewErrNotFound()
	}

	return value, nil
}

func (c *Collection[T]) Insert(ctx context.Context, value T) (string, error) {
	id := fmt.Sprintf("%v", c.currentId)
	c.currentId = c.currentId + 1
	c.preInsertCallback(value, id)
	c.Items[id] = value
	return id, nil
}

func (c *Collection[T]) Update(ctx context.Context, id string, value T) error {
	_, ok := c.Items[id]
	if !ok {
		return repository.NewErrNotFound()
	}

	c.Items[id] = value
	return nil
}

func (c *Collection[T]) Delete(ctx context.Context, id string) error {
	_, ok := c.Items[id]
	if !ok {
		return repository.NewErrNotFound()
	}

	delete(c.Items, id)
	return nil
}
