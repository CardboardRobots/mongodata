package mongodata

import (
	"context"
	"encoding/json"
	"testing"
)

func TestRepository(t *testing.T) {
	t.Skip()
	r := NewRepository(nil, bookMapper{})
	r.Insert(context.Background(), Book{}, nil)
}

type BookID string

func (b BookID) String() string { return string(b) }

type BookVersion int

func (bv BookVersion) Value() int { return int(bv) }

type Book struct {
	id      BookID
	version int
	title   string
}

func (b Book) ID() BookID   { return b.id }
func (b Book) Version() int { return b.version }

type BookModel struct {
	Title string `bson:"title"`
}

func (BookModel) Schema() string { return "" }

type bookMapper struct{}

func (b bookMapper) FromModel(id string, version int, m BookModel) Book {
	return Book{
		id:      BookID(id),
		version: version,
		title:   m.Title,
	}
}

func (bookMapper) ToModel(b Book) BookModel {
	return BookModel{
		Title: b.title,
	}
}

func (bookMapper) FromModelEvent(e Event) ([]byte, error) {
	return json.Marshal(e)
}

func (bookMapper) ToModelEvent(eType string, data []byte) (Event, bool) {
	return nil, false
}
