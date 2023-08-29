package website

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/model"
	"crawlers/pkg/stream"
	"github.com/reugn/go-streams/flow"
	"go.uber.org/zap"
)

const consumersNumber = 5
const messageRetries = 2

func RegisterStream() error {
	pr := NewTaskProcessor()
	// Register for page

	//consume catalog page ColumnUrl
	if err := catalogPageStream(pr); err != nil {
		return err
	}

	if err := novelStream(pr); err != nil {
		return err
	}

	return nil
}

// 解析page url得到每一个novel的url
// from: catalogPage stream => novel stream
func catalogPageStream(pr TaskProcessor) error {
	source, err := stream.NewRedisStreamSource(context.Background(), common.GetSystem().RedisClient,
		stream.CatalogPageUrlStream, stream.CatalogPageUrlStreamConsumer)
	if err != nil {
		return err
	}

	//item url
	sink := stream.NewRedisStreamSink(context.Background(), common.GetSystem().RedisClient,
		stream.NovelUrlStream)

	//convert the catalogPageTask message
	paramsConvertFlow := flow.NewMap(pr.HandleCatalogPageTask, 1)

	flowMap := flow.NewFlatMap(func(novelMsg []model.NovelTask) []model.NovelTask {
		return novelMsg
	}, 1)

	err = common.GetSystem().TaskPool.Submit(func() {
		source.Via(paramsConvertFlow).Via(flowMap).To(sink)
	})
	if err != nil {
		zap.S().Error("failed to submit task", zap.Error(err))
		return err
	}
	return nil
}

// 处理每一个novel
func novelStream(pr TaskProcessor) error {
	source, err := stream.NewRedisStreamSource(context.Background(), common.GetSystem().RedisClient,
		stream.NovelUrlStream, stream.NovelUrlStreamConsumer)
	if err != nil {
		return err
	}

	//item url
	sink := stream.NewRedisStreamSink(context.Background(), common.GetSystem().RedisClient,
		stream.ChapterUrlStream)

	paramsConvertFlow := flow.NewMap(pr.HandleNovelTask, consumersNumber)

	flowMap := flow.NewFlatMap(func(novelMsg []model.ChapterTask) []model.ChapterTask {
		return novelMsg
	}, 1)

	err = common.GetSystem().TaskPool.Submit(func() {
		source.Via(paramsConvertFlow).Via(flowMap).To(sink)
	})
	if err != nil {
		zap.L().Error("failed to submit task", zap.Error(err))
		return err
	}
	return nil
}

// 处理每一个novel
func chapterStream(pr TaskProcessor) error {
	source, err := stream.NewRedisStreamSource(context.Background(), common.GetSystem().RedisClient,
		stream.ChapterUrlStream, stream.ChapterUrlStreamConsumer)
	if err != nil {
		return err
	}

	//item url
	//sink := stream.NewRedisStreamSink(context.Background(), common.GetSystem().RedisClient,
	//	stream.ChapterUrlStream)

	paramsConvertFlow := flow.NewMap(pr.HandleChapterTask, consumersNumber)

	flowMap := flow.NewFlatMap(func(novelMsg []model.ChapterTask) []model.ChapterTask {
		return novelMsg
	}, 1)

	err = common.GetSystem().TaskPool.Submit(func() {
		source.Via(paramsConvertFlow).Via(flowMap)
	})
	if err != nil {
		zap.L().Error("failed to submit task", zap.Error(err))
		return err
	}
	return nil
}
