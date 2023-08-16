package dao

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/models/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CatalogInterface interface {
	FindById(ctx context.Context, id primitive.ObjectID) (*entity.Catalog, error)
	ExistsById(ctx context.Context, id primitive.ObjectID) (bool, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
}

type CatalogDaoImpl struct{}

func (c *CatalogDaoImpl) FindById(ctx context.Context, id primitive.ObjectID) (*entity.Catalog, error) {
	return FindById(ctx, id, common.CatalogCollection, &entity.Catalog{})
}

func (s *CatalogDaoImpl) ExistsById(ctx context.Context, id primitive.ObjectID) (bool, error) {
	site, err := FindById(ctx, id, common.CatalogCollection, &entity.Site{},
		&options.FindOneOptions{Projection: bson.M{common.ColumId: 1}})
	return site != nil, err
}

func (s *CatalogDaoImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	site, err := FindByMongoFilter(ctx, bson.M{common.ColumnName: name}, common.CatalogCollection, &entity.Site{},
		&options.FindOneOptions{Projection: bson.M{common.ColumId: 1}})
	return site != nil, err
}
