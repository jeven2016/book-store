package dao

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/models/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SiteInterface interface {
	FindById(ctx context.Context, id primitive.ObjectID) (*entity.Site, error)
	ExistsById(ctx context.Context, id primitive.ObjectID) (bool, error)
}

type SiteDaoImpl struct{}

func (s *SiteDaoImpl) FindById(ctx context.Context, id primitive.ObjectID) (*entity.Site, error) {
	return FindById(ctx, id, common.SiteCollection, &entity.Site{})
}

func (s *SiteDaoImpl) ExistsById(ctx context.Context, id primitive.ObjectID) (bool, error) {
	site, err := FindById(ctx, id, common.SiteCollection, &entity.Site{},
		&options.FindOneOptions{Projection: bson.M{common.ColumId: 1}})
	return site != nil, err
}
