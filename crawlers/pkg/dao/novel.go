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

type NovelInterface interface {
	ExistsByName(ctx context.Context, name string) (bool, error)
	Insert(ctx context.Context, novel *entity.Novel) (*primitive.ObjectID, error)
}

type NovelDaoImpl struct{}

func (n *NovelDaoImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	task, err := FindByMongoFilter(ctx, bson.M{common.ColumnName: name}, common.CollectionNovel, &model.CatalogPageTask{},
		&options.FindOneOptions{Projection: bson.M{common.ColumId: 1}})
	return task != nil, err
}

func (n *NovelDaoImpl) Insert(ctx context.Context, novel *entity.Novel) (*primitive.ObjectID, error) {
	collection := common.GetSystem().GetCollection(common.CollectionNovel)
	//for creating
	if !novel.Id.IsZero() {
		return nil, common.ErrDocumentIdExists
	}
	//check if name conflicts
	exists, err := n.ExistsByName(ctx, novel.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, common.ErrDuplicatedDocument
	}
	//insert
	if result, err := collection.InsertOne(ctx, novel, &options.InsertOneOptions{}); err != nil {
		return nil, err
	} else {
		insertedId := result.InsertedID.(primitive.ObjectID)
		return &insertedId, nil
	}
}
