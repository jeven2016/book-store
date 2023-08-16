package yzs8

import (
	"context"
	"crawlers/pkg/common"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"crawlers/pkg/models"
)

// converting
func Convert() {
	ConvertCatalog()
	ConvertArticle()
}

// ConvertCatalog 将ID、 parentId从string转为objectId
func ConvertCatalog() {
	// mongo
	client, err := common.CreateMongoClient(log)
	handleError(err)

	database := client.Database("books")
	catalogCol := database.Collection("catalog")

	cur, err := catalogCol.Find(context.Background(), bson.M{})

	var results = new([]models.CatalogDocMap)
	err = cur.All(context.Background(), results)
	handleError(err)

	for _, c := range *results {
		oldId := c.Id
		shouldCorrect := false
		if reflect.TypeOf(c.Id) == reflect.TypeOf(string("")) {
			shouldCorrect = true
			// 转换ID
			stringId := c.Id.(string)
			objectId, err := primitive.ObjectIDFromHex(stringId)
			handleError(err)
			c.Id = objectId
		}

		// 转换parentId
		if reflect.TypeOf(c.ParentId) == reflect.TypeOf(string("")) {
			if len(c.ParentId.(string)) > 0 {
				shouldCorrect = true
				// 转换ID
				stringId := c.ParentId.(string)
				objectId, err := primitive.ObjectIDFromHex(stringId)
				handleError(err)
				if err == nil {
					c.ParentId = objectId
				}
			}
		}

		// mongodb无法更新_id, 只能新建一个同时删除旧的
		if err == nil && shouldCorrect {
			// 删除旧文档
			_, err := catalogCol.DeleteOne(context.Background(), bson.M{"_id": oldId})
			if err != nil {
				panic(err)
			}

			// 插入新的
			_, err = catalogCol.InsertOne(context.Background(), c)
			handleError(err)
		}
	}

	println("finished")
}

// ConvertArticle 将ID、 parentId从string转为objectId
func ConvertArticle() {
	// mongo
	client, err := common.CreateMongoClient(log)
	handleError(err)

	database := client.Database("books")
	articleCol := database.Collection("article")

	cur, err := articleCol.Find(context.Background(), bson.M{})
	defer cur.Close(context.Background())

	var i = 0
	for cur.Next(context.TODO()) {
		if i%10000 == 0 {
			print("completed ", i)
		}
		var c *models.ArticleMap
		err := cur.Decode(&c)
		handleError(err)
		if err == nil {
			oldId := c.Id
			shouldCorrect := false
			if reflect.TypeOf(c.Id) == reflect.TypeOf(string("")) {
				shouldCorrect = true
				// 转换ID
				stringId := c.Id.(string)
				objectId, err := primitive.ObjectIDFromHex(stringId)
				handleError(err)
				c.Id = objectId
			}

			// 转换parentId
			if reflect.TypeOf(c.CatalogId) == reflect.TypeOf(string("")) {
				if len(c.CatalogId.(string)) > 0 {
					shouldCorrect = true
					// 转换ID
					stringId := c.CatalogId.(string)
					objectId, err := primitive.ObjectIDFromHex(stringId)
					handleError(err)
					if err == nil {
						c.CatalogId = objectId
					}
				}
			}

			// mongodb无法更新_id, 只能新建一个同时删除旧的
			if err == nil && shouldCorrect {
				// 删除旧文档
				_, err := articleCol.DeleteOne(context.Background(), bson.M{"_id": oldId})
				if err != nil {
					panic(err)
				}

				// 插入新的
				_, err = articleCol.InsertOne(context.Background(), c)
				handleError(err)
			}
		}
		i++
	}

	println("finished:", i)
}
