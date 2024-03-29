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
	"time"
)

type catalogPageTaskInterface interface {
	FindById(ctx context.Context, id primitive.ObjectID) (*model.CatalogPageTask, error)
	FindByUrl(ctx context.Context, url string) (*model.CatalogPageTask, error)
	ExistsById(ctx context.Context, id primitive.ObjectID) (bool, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	Save(ctx context.Context, task *model.CatalogPageTask) (*primitive.ObjectID, error)
}

type catalogPageTaskDaoImpl struct{}

func (c *catalogPageTaskDaoImpl) FindById(ctx context.Context, id primitive.ObjectID) (*model.CatalogPageTask, error) {
	return FindById(ctx, id, common.CollectionCatalogPageTask, &model.CatalogPageTask{})
}

func (c *catalogPageTaskDaoImpl) FindByUrl(ctx context.Context, url string) (*model.CatalogPageTask, error) {
	task, err := FindByMongoFilter(ctx, bson.M{common.ColumnUrl: url}, common.CollectionCatalogPageTask, &model.CatalogPageTask{})
	return task, err
}

func (s *catalogPageTaskDaoImpl) ExistsById(ctx context.Context, id primitive.ObjectID) (bool, error) {
	task, err := FindById(ctx, id, common.CollectionCatalogPageTask, &model.CatalogPageTask{},
		&options.FindOneOptions{Projection: bson.M{common.ColumId: 1}})
	return task != nil, err
}

func (s *catalogPageTaskDaoImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	task, err := FindByMongoFilter(ctx, bson.M{common.ColumnName: name}, common.CollectionCatalogPageTask, &model.CatalogPageTask{},
		&options.FindOneOptions{Projection: bson.M{common.ColumId: 1}})
	return task != nil, err
}

func (c *catalogPageTaskDaoImpl) Save(ctx context.Context, task *model.CatalogPageTask) (*primitive.ObjectID, error) {
	collection := common.GetSystem().GetCollection(common.CollectionCatalogPageTask)
	if collection == nil {
		zap.L().Error("collection not found: " + common.CollectionCatalogPageTask)
		return nil, errors.New("collection not found: " + common.CollectionCatalogPageTask)
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
		curTime := time.Now()
		task.LastUpdated = &curTime
		taskBytes, err := bson.Marshal(task)
		if err != nil {
			return nil, err
		}
		var doc bson.D
		if err = bson.Unmarshal(taskBytes, &doc); err != nil {
			return nil, err
		}
		_, err = collection.UpdateOne(ctx, bson.M{common.ColumId: task.Id, common.ColumnSiteName: task.SiteName}, bson.M{"$set": doc})
		return &task.Id, err
	}
}
