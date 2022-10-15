package memory

import (
	"context"
	"testing"
)

func TestGet(t *testing.T) {
	ctx := context.TODO()
	c := NewCollection(func(value, id string) {})
	want := "abcde"
	id, err := c.Create(ctx, want)
	if err != nil {
		t.Errorf("%v", err)
	}

	value, err := c.GetById(context.TODO(), id)
	if err != nil {
		t.Errorf("%v", err)
	}

	if value != want {
		t.Errorf("Value is wrong.  Received: %v, Expected: %v", value, want)
	}
}

func TestInsert(t *testing.T) {
	ctx := context.TODO()
	valueCallback := ""
	idCallback := ""
	c := NewCollection(func(value, id string) {
		valueCallback = value
		idCallback = id
	})
	id0, _ := c.Create(ctx, "abcde")
	id1, _ := c.Create(ctx, "fghij")
	if id0 == "" || id1 == "" || id0 == id1 {
		t.Errorf("ids should be non-nil, and unique.  Received: %v, %v", id0, id1)
	}
	if valueCallback != "fghij" || idCallback != id1 {
		t.Errorf("Callback is wrong.  Received %v, %v", valueCallback, idCallback)
	}
	item0, _ := c.GetById(ctx, id0)
	item1, _ := c.GetById(ctx, id1)
	if item0 == "" || item1 == "" {
		t.Errorf("Items not inserted: %v, %v", item0, item1)
	}
}

func TestReplace(t *testing.T) {
	ctx := context.TODO()
	c := NewCollection(func(value, id string) {})
	want := "fghij"
	id, _ := c.Create(ctx, "abcde")

	_, err := c.Update(ctx, id, want)
	if err != nil {
		t.Errorf("%v", err)
	}

	value, err := c.GetById(ctx, id)
	if err != nil {
		t.Errorf("%v", err)
	}

	if value != want {
		t.Errorf("Value is wrong.  Received: %v, Expected: %v", value, want)
	}
}

func TestDelete(t *testing.T) {
	ctx := context.TODO()
	c := NewCollection(func(value, id string) {})
	id, err := c.Create(ctx, "abcde")
	if err != nil {
		t.Errorf("%v", err)
	}

	_, err = c.Delete(ctx, id)
	if err != nil {
		t.Errorf("%v", err)
	}

	_, err = c.Delete(ctx, id)
	if err == nil {
		t.Errorf("Error should be NotFoundError.  Received: %v", err)
	}
}

// func sliceEqual[T comparable](a []T, b []T) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}

// 	for i, value := range a {
// 		if value != b[i] {
// 			return false
// 		}
// 	}

// 	return true
// }
