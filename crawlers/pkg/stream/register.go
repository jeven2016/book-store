package stream

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/dao"
	"crawlers/pkg/models"
	"crawlers/pkg/website"
	"encoding/json"
	"github.com/reugn/go-streams/flow"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const novelConsumers = 5
const messageRetries = 2

func RegisterStream() error {
	// Register for page

	//consume catalog page ColumnUrl
	if err := consumePageUrl(); err != nil {
		return err
	}

	if err := consumeNovel(); err != nil {
		return err
	}

	return nil
}

// 解析page url得到每一个novel的url
// from: catalogPage stream => novel stream
func consumePageUrl() error {

	source, err := NewRedisStreamSource(context.Background(), common.GetSystem().RedisClient,
		CatalogPageUrlStream, CatalogPageUrlStreamConsumer)
	if err != nil {
		return err
	}

	//item url
	sink := NewRedisStreamSink(context.Background(), common.GetSystem().RedisClient,
		NovelUrlStream)

	paramsConvertFlow := flow.NewMap(func(jsonData string) []models.NovelTask {
		var catalogPage models.CatalogPageTask

		if err = json.Unmarshal([]byte(jsonData), &catalogPage); err != nil {
			zap.L().Error("json decode", zap.Error(err), zap.String("data", jsonData))
			return nil
		}
		zap.L().Info("catalogPage task", zap.String("json", jsonData))

		//check if page url is duplicated
		//for task.CatalogPageTask
		exists, shouldStopped := isDuplicatedUrl(catalogPage.SiteName, common.CatalogPageTaskCollection,
			catalogPage.Url,
			bson.M{
				common.ColumnCatalogId: catalogPage.CatalogId, //catalogPage.catalogId
				common.ColumnUrl:       catalogPage.Url,       //catalogPage.Url
			}, "catalog page link")

		if shouldStopped || exists {
			return nil
		}

		downloader := website.GetSiteCrawler(catalogPage.SiteName)
		if downloader == nil {
			zap.L().Error("site downloader not found", zap.String("SiteName", catalogPage.SiteName))
			return nil
		}
		if novelMsgs, err := downloader.HandleCatalogPage(context.Background(), &catalogPage); err != nil {
			//保存失败,把消息归还到队列中去或最终写入DB
			if err = returnMessageToQueue(&catalogPage, CatalogPageUrlStream); err != nil {
				zap.L().Error("failed to return this message to queue", zap.String("jsonData", jsonData), zap.Error(err))
			}
			return nil
		} else {
			//已经处理过，记录该url
			if err = website.GetSiteProcessor(catalogPage.SiteName).SaveCatalogPageUrl(&catalogPage); err != nil {
				zap.L().Error("failed to save catalog page url", zap.String("jsonData", jsonData), zap.Error(err))

				//保存失败,把消息归还到队列中去
				if err = returnMessageToQueue(&catalogPage, CatalogPageUrlStream); err != nil {
					zap.L().Error("failed to return this message to queue", zap.String("jsonData", jsonData), zap.Error(err))
				}
				return nil
			}

			return novelMsgs
		}

	}, 1)

	flowMap := flow.NewFlatMap(func(novelMsg []models.NovelTask) []models.NovelTask {
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
func consumeNovel() error {
	streamName := NovelUrlStream
	source, err := NewRedisStreamSource(context.Background(), common.GetSystem().RedisClient,
		streamName, NovelUrlStreamConsumer)
	if err != nil {
		return err
	}

	//item url
	sink := NewRedisStreamSink(context.Background(), common.GetSystem().RedisClient,
		ChapterUrlStream)

	paramsConvertFlow := flow.NewMap(func(jsonData string) []models.ChapterTask {
		var NovelTask models.NovelTask

		if err = json.Unmarshal([]byte(jsonData), &NovelTask); err != nil {
			zap.L().Error("json decode", zap.Error(err), zap.String("data", jsonData))
			return nil
		}
		zap.S().Info("Got jsonData", jsonData)

		exists, shouldStopped := isDuplicatedUrl(NovelTask.SiteName, common.NovelCollection, NovelTask.Url, bson.M{
			common.SiteName:  NovelTask.SiteName, //SiteName
			common.NovelLink: NovelTask.Url,      //novelLink
		}, "novel link")

		if shouldStopped || exists {
			return nil
		}

		downloader := website.GetSiteCrawler(NovelTask.SiteName)
		if downloader == nil {
			zap.L().Error("site downloader not found", zap.String("SiteName", NovelTask.SiteName))
			return nil
		}
		if chapterMgs, err := downloader.HandleNovelPage(context.Background(), &NovelTask); err != nil {
			//保存失败,把消息归还到队列中去
			if err = returnMessageToQueue(&NovelTask, streamName); err != nil {
				zap.L().Error("failed to return this message to queue", zap.String("jsonData", jsonData), zap.Error(err))
			}
			return nil
		} else {
			//已经处理过，记录该url
			if err = website.GetSiteProcessor(NovelTask.SiteName).SaveNovelPageUrl(&NovelTask); err != nil {
				zap.L().Error("failed to save novel page url", zap.String("jsonData", jsonData), zap.Error(err))

				//保存失败,把消息归还到队列中去
				if err = returnMessageToQueue(&NovelTask, streamName); err != nil {
					zap.L().Error("failed to return this message to queue", zap.String("jsonData", jsonData), zap.Error(err))
				}
				return nil
			}

			return chapterMgs
		}

	}, novelConsumers)

	flowMap := flow.NewFlatMap(func(novelMsg []models.ChapterTask) []models.ChapterTask {
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

// 检查是否已经处理过的url
func isDuplicatedUrl(SiteName, collectionName, url string, bsonFilter bson.M, logName string) (bool /*existence*/, bool /*interrupted*/) {
	//判断是否已经处理过
	exists, err := common.Exists(context.Background(), url, func() (any, error) {
		return dao.ExistsByMongoFilter(context.Background(), bsonFilter, collectionName,
			&options.FindOneOptions{Projection: bson.M{common.ColumId: 1}})
	})
	if err != nil {
		zap.L().Error("failed to check if the "+logName+" exists", zap.String("SiteName", SiteName), zap.Error(err))
		return false, true
	}
	if exists {
		zap.L().Info(logName+" already processed", zap.String("SiteName",
			SiteName), zap.String(logName, url))
		return true, false
	}
	return false, false
}

func returnMessageToQueue(msg models.Resource, streamName string) error {
	//if msg.GetErrorCount() >= messageRetries {
	//	return errors.New("max retires exceeded:" + strconv.Itoa(int(msg.GetErrorCount())))
	//}
	//
	////todo bug: 一直循环中
	////把消息归还到队列中去
	//msg.IncreaseErrorCount()
	//if err := common.GetSystem().RedisClient.PublishMessage(context.Background(),
	//	msg, streamName); err != nil {
	//	return err
	//}
	return nil
}
