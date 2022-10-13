package db

import (
	"context"
	"errors"

	"github.com/cardboardrobots/mongodata/collection"
	"github.com/cardboardrobots/mongodata/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Collection[T utils.Valid] struct {
	Collection *mongo.Collection
}

func NewCollection[T utils.Valid](da *DataAccess, collectionName string) Collection[T] {
	collection := da.database.Collection(collectionName)
	return Collection[T]{
		collection,
	}
}

func (c *Collection[T]) GetList(
	ctx context.Context,
	query primitive.M,
) (collection.ListResult[T], error) {
	cursor, err := c.Collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	var results []T
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return collection.NewListResult(len(results), results), nil
}

func (c *Collection[T]) GetById(
	ctx context.Context,
	id string,
) (*T, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result T
	err = c.Collection.FindOne(ctx, bson.M{
		"_id": _id,
	}).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &result, err
}

func (c *Collection[T]) Insert(
	ctx context.Context,
	data T,
) (string, error) {
	error := data.Valid()
	if error != nil {
		return "", error
	}

	result, error := c.Collection.InsertOne(ctx, data)
	return result.InsertedID.(primitive.ObjectID).Hex(), error
}

func (c *Collection[T]) Update(
	ctx context.Context,
	id string,
	data T,
) (bool, error) {
	err := data.Valid()
	if err != nil {
		return false, err
	}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	updateDate := make(map[string]interface{})
	updateDate["$set"] = data
	result, err := c.Collection.UpdateByID(ctx, _id, updateDate)
	if err != nil {
		return false, err
	}

	return (result.ModifiedCount + result.UpsertedCount) > 0, err
}

func (c *Collection[T]) Delete(
	ctx context.Context,
	id string,
) (bool, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	result, err := c.Collection.DeleteOne(ctx, bson.M{
		"_id": _id,
	})
	return result.DeletedCount > 0, err
}