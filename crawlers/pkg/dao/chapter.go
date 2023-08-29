package dao

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/model"
	"crawlers/pkg/model/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type ChapterInterface interface {
	ExistsByName(ctx context.Context, name string) (bool, error)
	Insert(ctx context.Context, novel *entity.Chapter) (*primitive.ObjectID, error)
	BulkInsert(ctx context.Context, chapters []*entity.Chapter, novelId *primitive.ObjectID) error
}

type ChapterDaoImpl struct{}

func (n *ChapterDaoImpl) ExistsByName(ctx context.Context, name string) (bool, error) {
	task, err := FindByMongoFilter(ctx, bson.M{common.ColumnName: name}, common.CollectionChapter, &model.CatalogPageTask{},
		&options.FindOneOptions{Projection: bson.M{common.ColumId: 1}})
	return task != nil, err
}

func (n *ChapterDaoImpl) Insert(ctx context.Context, novel *entity.Chapter) (*primitive.ObjectID, error) {
	collection := common.GetSystem().GetCollection(common.CollectionChapter)
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

func (n *ChapterDaoImpl) BulkInsert(ctx context.Context, chapters []*entity.Chapter, novelId *primitive.ObjectID) error {
	collection := common.GetSystem().GetCollection(common.CollectionChapter)

	documents := make([]interface{}, len(chapters))

	//保存chapters
	for i := 0; i < len(chapters); i++ {
		chapters[i].NovelId = *novelId
		documents[i] = chapters[i]
	}

	// 指定每个批次的文档数量
	BulkBatchSize := 10

	// 计算批次数量
	numBatches := (len(documents) + BulkBatchSize - 1) / BulkBatchSize
	// 分批插入文档
	for i := 0; i < numBatches; i++ {
		// 计算当前批次的起始和结束索引
		startIndex := i * BulkBatchSize
		endIndex := (i + 1) * BulkBatchSize
		if endIndex > len(chapters) {
			endIndex = len(chapters)
		}

		// 获取当前批次的文档
		batch := documents[startIndex:endIndex]

		// 执行批量插入操作
		_, err := collection.InsertMany(ctx, batch)
		if err != nil {
			zap.L().Error("failed to insert chapters", zap.String("novelId", novelId.Hex()))
			return err
		}
		zap.S().Info("the number of inserted chapters: ", numBatches*(i+1))
	}
	return nil
}
