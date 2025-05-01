package mongodata

import (
	"context"
	"errors"
	"iter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection[M any] struct {
	collection *mongo.Collection
	decoders   map[string]Decoder[M]
}

func NewCollection[M any](
	collection *mongo.Collection,
) *Collection[M] {
	return &Collection[M]{
		collection: collection,
	}
}

func (c *Collection[M]) Insert(ctx context.Context, id string, m *M) error {
	_, err := c.collection.InsertOne(ctx,
		m)
	return err
}

func (c *Collection[M]) Upsert(ctx context.Context, id string, m *M) error {
	_, err := c.collection.ReplaceOne(ctx,
		bson.M{"_id": id},
		m,
		options.Replace().SetUpsert(true))
	return err
}

func (c *Collection[M]) Replace(ctx context.Context, filter FilterBuilder, m *M) error {
	result, err := c.collection.ReplaceOne(ctx,
		filter.Build(),
		m)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrNoMatch
	}

	return nil
}

func (c *Collection[M]) Get(ctx context.Context, id string) (*M, error) {
	m, err := c.decodeSingle(c.collection.FindOne(ctx,
		bson.M{"_id": id}))
	if errors.Is(err, mongo.ErrNoDocuments) {
		err = ErrNotFound
	}

	return m, err
}

func (c Collection[M]) decodeSingle(sr *mongo.SingleResult) (*M, error) {
	if len(c.decoders) == 0 {
		return Decode[M](sr)
	}

	b, err := sr.Raw()
	if err != nil {
		return nil, err
	}

	d, ok := c.decoders[b.Lookup("_schema").String()]
	if !ok {
		return Decode[M](sr)
	}

	return d.Decode(sr)
}

func (c *Collection[M]) Delete(ctx context.Context, id string) error {
	_, err := c.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (c *Collection[M]) GetList(ctx context.Context, filter FilterBuilder, sort *SortBuilder) iter.Seq2[*M, error] {
	return func(yield func(*M, error) bool) {
		f := filter.Build()
		o := sort.Build()
		result, err := c.collection.Find(ctx, f, o)
		if err != nil {
			yield(nil, err)
			return
		}

		defer result.Close(ctx)

		for result.Next(ctx) {
			m, err := c.decodeCursor(result)
			if err != nil {
				if errors.Is(err, mongo.ErrNoDocuments) {
					err = ErrNotFound
				}
				yield(nil, err)
				return
			}

			if !yield(m, nil) {
				return
			}
		}
	}
}

func (c Collection[M]) decodeCursor(cr *mongo.Cursor) (*M, error) {
	if len(c.decoders) == 0 {
		return Decode[M](cr)
	}

	d, ok := c.decoders[cr.Current.Lookup("_schema").String()]
	if !ok {
		return Decode[M](cr)
	}

	return d.Decode(cr)
}

type Change[M any] struct {
	Value M
	Token string
}

func newChange[M any](cs *mongo.ChangeStream) (Change[M], error) {
	m, err := Decode[M](cs)
	return Change[M]{Value: *m, Token: cs.ResumeToken().String()}, err
}

func (c *Collection[M]) Watch(ctx context.Context, pipeline any, token string) iter.Seq2[Change[M], error] {
	return func(yield func(Change[M], error) bool) {
		o := options.ChangeStream()
		if token != "" {
			o.SetResumeAfter(bson.Raw(token))
		}

		cs, err := c.collection.Watch(ctx, pipeline, o)
		if err != nil {
			yield(Change[M]{Token: token}, err)
			return
		}

		defer cs.Close(ctx)

		for cs.Next(ctx) {
			if c, err := newChange[M](cs); err != nil {
				yield(c, err)
				return
			} else if !yield(c, nil) {
				return
			}
		}
	}
}
