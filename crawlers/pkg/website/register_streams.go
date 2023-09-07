package website

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/model"
	"crawlers/pkg/stream"
	"github.com/reugn/go-streams/extension"
	"github.com/reugn/go-streams/flow"
	"go.uber.org/zap"
)

const consumersNumber = 5

func RegisterStream(ctx context.Context) error {
	pr := NewTaskProcessor()

	//consume catalog page ColumnUrl
	if err := catalogPageStream(ctx, pr); err != nil {
		return err
	}

	if err := novelStream(ctx, pr); err != nil {
		return err
	}

	if err := chapterStream(ctx, pr); err != nil {
		return err
	}
	return nil
}

type StreamStepDefinition[T, R, E, U any] struct {
	sourceStream        string
	sourceConsumerGroup string
	destinationStream   string
	convertFunc         flow.MapFunction[T, R]
	flowFlatMap         flow.FlatMap[E, U]
}

// 解析page url得到每一个novel的url
// from: catalogPage stream => novel stream
func catalogPageStream(ctx context.Context, pr TaskProcessor) error {
	source, err := stream.NewRedisStreamSource(context.Background(), common.GetSystem().RedisClient,
		stream.CatalogPageUrlStream, stream.CatalogPageUrlStreamConsumer)
	if err != nil {
		return err
	}

	err = common.GetSystem().TaskPool.Submit(func() {
		source.
			Via(flow.NewMap(pr.HandleCatalogPageTask, 1)).
			Via(flow.NewFlatMap(func(novelMsg []model.NovelTask) []model.NovelTask {
				return novelMsg
			}, 1)).
			To(stream.NewRedisStreamSink(ctx, common.GetSystem().RedisClient,
				stream.NovelUrlStream))
	})
	if err != nil {
		zap.S().Error("failed to submit task", zap.Error(err))
		return err
	}
	return nil
}

// 处理每一个novel
func novelStream(ctx context.Context, pr TaskProcessor) error {
	source, err := stream.NewRedisStreamSource(context.Background(), common.GetSystem().RedisClient,
		stream.NovelUrlStream, stream.NovelUrlStreamConsumer)
	if err != nil {
		return err
	}

	//item url
	sink := stream.NewRedisStreamSink(ctx, common.GetSystem().RedisClient,
		stream.ChapterUrlStream)

	err = common.GetSystem().TaskPool.Submit(func() {
		source.
			Via(flow.NewMap(pr.HandleNovelTask, consumersNumber)).
			Via(flow.NewFlatMap(func(novelMsg []model.ChapterTask) []model.ChapterTask {
				return novelMsg
			}, 1)).
			To(sink)
	})
	if err != nil {
		zap.L().Error("failed to submit task", zap.Error(err))
		return err
	}
	return nil
}

// 处理每一个novel
func chapterStream(ctx context.Context, pr TaskProcessor) error {
	source, err := stream.NewRedisStreamSource(ctx, common.GetSystem().RedisClient,
		stream.ChapterUrlStream, stream.ChapterUrlStreamConsumer)
	if err != nil {
		return err
	}

	err = common.GetSystem().TaskPool.Submit(func() {
		source.
			Via(flow.NewMap(pr.HandleChapterTask, consumersNumber)).
			To(extension.NewIgnoreSink())
	})
	if err != nil {
		zap.L().Error("failed to submit task", zap.Error(err))
		return err
	}
	return nil
}
