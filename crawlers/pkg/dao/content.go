package dao

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/model"
	"crawlers/pkg/model/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContentInterface interface {
	ExistsByParentId(ctx context.Context, id primitive.ObjectID) (bool, error)
	Insert(ctx context.Context, content *entity.Content) (*primitive.ObjectID, error)
}

type ContentDaoImpl struct{}

func (c *ContentDaoImpl) Insert(ctx context.Context, content *entity.Content) (*primitive.ObjectID, error) {
	collection := common.GetSystem().GetCollection(common.CollectionContent)
	//for creating
	if !content.Id.IsZero() {
		return nil, common.ErrDocumentIdExists
	}
	//check if name conflicts
	exists, err := c.ExistsByParentId(ctx, content.ParentId)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ErrDuplicatedDocument
	}
	//insert
	if result, err := collection.InsertOne(ctx, content, &options.InsertOneOptions{}); err != nil {
		return nil, err
	} else {
		insertedId := result.InsertedID.(primitive.ObjectID)
		return &insertedId, nil
	}
}

func (c *ContentDaoImpl) ExistsByParentId(ctx context.Context, id primitive.ObjectID) (bool, error) {
	task, err := FindByMongoFilter(ctx, bson.M{common.ColumnName: name}, common.CollectionChapter, &model.CatalogPageTask{},
		&options.FindOneOptions{Projection: bson.M{common.ColumId: 1}})
	return task != nil, err
}
