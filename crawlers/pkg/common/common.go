package common

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

func CreateMongoClient(log *zap.Logger) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接MongoDB
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://db_user:db_pwd@localhost:27017/books?retryWrites=true&w=majority&authSource=admin&maxPoolSize=20"))

	if err != nil {
		log.Error("Cannot connect to mongodb", zap.Error(err))
		return nil, err
	}

	// 检测MongoDB是否连接成功
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Error("Failed to Ping for mongodb", zap.Error(err))
		return nil, err
	}

	return client, nil
}
