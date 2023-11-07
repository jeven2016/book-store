package dao

import (
	"context"
	"crawlers/pkg/common"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func FindById[T any](ctx context.Context, id primitive.ObjectID, collection string,
	obj *T, opts ...*options.FindOneOptions) (*T, error) {
	return FindByMongoFilter(ctx, bson.M{common.ColumId: id}, collection, obj, opts...)
}

func FindByMongoFilter[T any](ctx context.Context, mongoFilter interface{}, collection string,
	decodedObj *T, opts ...*options.FindOneOptions) (*T, error) {
	col := common.GetSystem().GetCollection(collection)
	if col == nil {
		return nil, errors.New("collection not found: " + collection)
	}
	if result := col.FindOne(ctx, mongoFilter, opts...); result.Err() != nil {
		err := result.Err()
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	} else {
		if decodedObj != nil {
			err := result.Decode(decodedObj)
			return decodedObj, err
		}
		return nil, common.ErrDecodingDocument
	}
}

func ExistsByMongoFilter(ctx context.Context, mongoFilter interface{},
	collection string, opts ...*options.FindOneOptions) (bool, error) {
	var obj *interface{}
	_, err := FindByMongoFilter(ctx, mongoFilter, collection, obj, opts...)
	if err != nil && errors.Is(err, common.ErrDecodingDocument) {
		return true, nil
	}
	return false, err
}
