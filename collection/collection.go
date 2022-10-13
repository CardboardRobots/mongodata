package collection

import (
	"context"
	"fmt"
)

type Collection[T any] struct {
	Items     map[string]T
	currentId int
}

func NewCollection[T any]() Collection[T] {
	return Collection[T]{
		Items: map[string]T{},
	}
}

func (c *Collection[T]) GetList(ctx context.Context) (ListResult[T], error) {
	count, _ := c.Count(ctx)
	return NewListResult(count, c.Slice()), nil
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
		return zero, NewErrNotFound()
	}

	return value, nil
}

func (c *Collection[T]) Insert(value T, callback func(value T, id string)) string {
	id := fmt.Sprintf("%v", c.currentId)
	c.currentId = c.currentId + 1
	callback(value, id)
	c.Items[id] = value
	return id
}

func (c *Collection[T]) InsertAtId(id string, value T) {
	c.Items[id] = value
}

func (c *Collection[T]) Update(ctx context.Context, id string, value T) error {
	_, ok := c.Items[id]
	if !ok {
		return NewErrNotFound()
	}

	c.Items[id] = value
	return nil
}

func (c *Collection[T]) Delete(ctx context.Context, id string) error {
	_, ok := c.Items[id]
	if !ok {
		return NewErrNotFound()
	}

	delete(c.Items, id)
	return nil
}
