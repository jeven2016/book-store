package dao

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/model"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type novelTaskInterface interface {
	FindByUrl(ctx context.Context, url string) (*model.NovelTask, error)
	Save(ctx context.Context, task *model.NovelTask) (*primitive.ObjectID, error)
}

type novelTaskDaoImpl struct{}

func (c *novelTaskDaoImpl) FindByUrl(ctx context.Context, url string) (*model.NovelTask, error) {
	task, err := FindByMongoFilter(ctx, bson.M{common.ColumnUrl: url}, common.CollectionNovelTask, &model.NovelTask{})
	return task, err
}

func (c *novelTaskDaoImpl) Save(ctx context.Context, task *model.NovelTask) (*primitive.ObjectID, error) {
	collection := common.GetSystem().GetCollection(common.CollectionNovelTask)
	if collection == nil {
		zap.L().Error("collection not found: " + common.CollectionNovelTask)
		return nil, errors.New("collection not found: " + common.CollectionNovelTask)
	}
	if task.Id.IsZero() {
		//insert
		if result, err := collection.InsertOne(ctx, task, &options.InsertOneOptions{}); err != nil {
			return nil, err
		} else {
			insertedId := result.InsertedID.(primitive.ObjectID)
			return &insertedId, nil
		}
	} else {
		//update
		taskBytes, err := bson.Marshal(task)
		if err != nil {
			return nil, err
		}
		var doc bson.D
		if err = bson.Unmarshal(taskBytes, &doc); err != nil {
			return nil, err
		}
		_, err = collection.UpdateOne(ctx,
			bson.M{common.ColumId: task.Id, common.ColumnSiteName: task.SiteName}, bson.M{"$set": doc})
		return &task.Id, err
	}
}
