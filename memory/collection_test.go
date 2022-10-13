package collection

import (
	"context"
	"testing"
)

func TestLength(t *testing.T) {
	ctx := context.TODO()
	c := NewCollection[string](func(value, id string) {})
	length, _ := c.Count(ctx)
	if length != 0 {
		t.Errorf("Length should be 0, received: %v", length)
	}
	c.Insert(ctx, "abcde")
	length, _ = c.Count(ctx)
	if length != 1 {
		t.Errorf("Length should be 1, received: %v", length)
	}
}

// func TestSlice(t *testing.T) {
// 	c := NewCollection[string](func(value, id string) {})

// 	slice := c.Slice()
// 	length := len(slice)
// 	if length != 0 {
// 		t.Errorf("Length should be 0, received: %v", length)
// 	}

// 	c.Insert("abcde")

// 	slice = c.Slice()
// 	length = len(slice)
// 	if length != 1 {
// 		t.Errorf("Length should be 1, received: %v", length)
// 	}

// 	want := []string{"abcde"}
// 	if !sliceEqual(want, slice) {
// 		t.Errorf("Slices are not equal.  Received: %v, Expected: %v", slice, want)
// 	}
// }

func TestGet(t *testing.T) {
	ctx := context.TODO()
	c := NewCollection[string](func(value, id string) {})
	want := "abcde"
	id, err := c.Insert(ctx, want)
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
	c := NewCollection[string]()
	valueCallback := ""
	idCallback := ""
	id0 := c.Insert("abcde", func(value, id string) {
		valueCallback = value
		idCallback = id
	})
	id1 := c.Insert("fghij", func(value, id string) {})
	if id0 == "" || id1 == "" || id0 == id1 {
		t.Errorf("ids should be non-nil, and unique.  Received: %v, %v", id0, id1)
	}
	if valueCallback != "abcde" || idCallback != id0 {
		t.Errorf("Callback is wrong.  Received %v, %v", valueCallback, idCallback)
	}
	length, _ := c.Count(ctx)
	want := 2
	if length != want {
		t.Errorf("length is wrong.  Received: %v, Expected: %v", length, want)
	}
}

func TestReplace(t *testing.T) {
	ctx := context.TODO()
	c := NewCollection[string]()
	want := "fghij"
	id := c.Insert("abcde", func(value, id string) {})

	err := c.Update(ctx, id, want)
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
	c := NewCollection[string](func(value, id string) {})
	id, err := c.Insert(ctx, "abcde")
	if err != nil {
		t.Errorf("%v", err)
	}

	err = c.Delete(ctx, id)
	if err != nil {
		t.Errorf("%v", err)
	}

	err = c.Delete(ctx, id)
	if err == nil {
		t.Errorf("Error should be NotFoundError.  Received: %v", err)
	}
}

func sliceEqual[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i, value := range a {
		if value != b[i] {
			return false
		}
	}

	return true
}
