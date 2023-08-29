package website

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/dao"
	"crawlers/pkg/model"
	"encoding/json"
	"github.com/duke-git/lancet/v2/convertor"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"time"
)

type TaskProcessor interface {
	ParsePageUrls(siteName, originPageUrl string) ([]string, error)
	HandleCatalogPageTask(jsonData string) []model.NovelTask
	HandleNovelTask(jsonData string) []model.ChapterTask
	HandleChapterTask(jsonData string) interface{}
}

type DefaultTaskProcessor struct{}

func (d DefaultTaskProcessor) ParsePageUrls(siteName, originPageUrl string) ([]string, error) {
	sys := common.GetSystem()
	cfg, err := common.GetSiteConfig(siteName)
	if err != nil {
		sys.Log.Sugar().Error("Could not find site config", zap.Error(err))
		return nil, err
	}
	return common.GenPageUrls(cfg.RegexSettings.ParsePageRegex, originPageUrl, cfg.RegexSettings.PagePrefix)
}

func NewTaskProcessor() TaskProcessor {
	return &DefaultTaskProcessor{}
}

func (d DefaultTaskProcessor) HandleCatalogPageTask(jsonData string) (novelMsgs []model.NovelTask) {
	var catalogPageTask model.CatalogPageTask
	var err error

	if err = json.Unmarshal([]byte(jsonData), &catalogPageTask); err != nil {
		zap.L().Error("json decode", zap.Error(err), zap.String("data", jsonData))
		return nil
	}
	zap.L().Info("catalogPageTask task", zap.String("json", jsonData))

	//check if page url is duplicated
	exists, err := d.isDuplicatedCatalogPageTask(&catalogPageTask,
		common.CollectionCatalogPageTask,
		catalogPageTask.Url,
		bson.M{
			common.ColumnCatalogId: catalogPageTask.CatalogId, //catalogPageTask.catalogId
			common.ColumnUrl:       catalogPageTask.Url,       //catalogPageTask.Url
		})

	if err != nil || exists {
		return nil
	}

	downloader := GetSiteCrawler(catalogPageTask.SiteName)
	if downloader == nil {
		zap.L().Error("site downloader not found", zap.String("SiteName", catalogPageTask.SiteName))
		return nil
	}

	//check if it exists in db
	var existingTask *model.CatalogPageTask
	if existingTask, err = dao.CatalogPageTaskDao.FindByUrl(context.Background(), catalogPageTask.Url); err != nil {
		zap.L().Error("failed to retrieve catalog page task", zap.String("jsonData", jsonData), zap.Error(err))
		return nil
	}

	currentTime := time.Now()
	if novelMsgs, err = downloader.HandleCatalogPage(context.Background(), &catalogPageTask); err != nil {
		//save failed, update the status
		if existingTask != nil {
			if err = convertor.CopyProperties(&catalogPageTask, existingTask); err != nil {
				zap.L().Error("failed to copy properties of catalog page task", zap.Error(err))
				return nil
			}
			//如果之前重试过，重试次数加1
			if catalogPageTask.Status == common.TaskStatusFailed ||
				catalogPageTask.Status == common.TaskStatusRetryFailed {
				catalogPageTask.Retries++
				catalogPageTask.Status = common.TaskStatusRetryFailed
			}
		} else {
			catalogPageTask.Status = common.TaskStatusFailed
		}
		catalogPageTask.LastUpdated = &currentTime
	} else {
		//已经处理过，记录该url
		catalogPageTask.Status = common.TaskStatusFinished
		catalogPageTask.CreatedDate = &currentTime
	}

	if _, err = dao.CatalogPageTaskDao.Save(context.Background(), &catalogPageTask); err != nil {
		zap.L().Error("failed to save catalogPageTask", zap.Error(err))
	}

	return

}

func (d DefaultTaskProcessor) HandleNovelTask(jsonData string) (chapterMessages []model.ChapterTask) {
	var NovelTask model.NovelTask
	var err error

	if err = json.Unmarshal([]byte(jsonData), &NovelTask); err != nil {
		zap.L().Error("json decode", zap.Error(err), zap.String("data", jsonData))
		return nil
	}
	zap.L().Info("novel task", zap.String("json", jsonData))

	//check if page url is duplicated
	exists, err := d.isDuplicatedNovelTask(&NovelTask,
		common.CollectionNovelTask,
		NovelTask.Url,
		bson.M{
			common.ColumnCatalogId: NovelTask.CatalogId, //catalogPageTask.catalogId
			common.ColumnUrl:       NovelTask.Url,       //catalogPageTask.Url
		})
	if err != nil {
		zap.L().Warn("error occurs", zap.Error(err))
		return nil
	}
	if exists {
		zap.L().Warn("duplicated novel to crawl", zap.String("jsonData", jsonData))
		return nil
	}

	downloader := GetSiteCrawler(NovelTask.SiteName)
	if downloader == nil {
		zap.L().Error("site downloader not found", zap.String("SiteName", NovelTask.SiteName))
		return nil
	}

	//check if it exists in db
	var existingTask *model.NovelTask
	if existingTask, err = dao.NovelTaskDao.FindByUrl(context.Background(), NovelTask.Url); err != nil {
		zap.L().Error("failed to retrieve novel page task", zap.String("jsonData", jsonData), zap.Error(err))
		return nil
	}

	currentTime := time.Now()
	if chapterMessages, err = downloader.HandleNovelPage(context.Background(), &NovelTask); err != nil {
		//save failed, update the status
		if existingTask != nil {
			if err = convertor.CopyProperties(&NovelTask, existingTask); err != nil {
				zap.L().Error("failed to copy properties of catalog page task", zap.Error(err))
				return nil
			}
			//如果之前重试过，重试次数加1
			if NovelTask.Status == common.TaskStatusFailed ||
				NovelTask.Status == common.TaskStatusRetryFailed {
				NovelTask.Retries++
				NovelTask.Status = common.TaskStatusRetryFailed
			}
		} else {
			NovelTask.Status = common.TaskStatusFailed
		}
		NovelTask.LastUpdated = &currentTime
	} else {
		//已经处理过，记录该url
		NovelTask.Status = common.TaskStatusFinished
		NovelTask.CreatedDate = &currentTime
	}

	if _, err = dao.NovelTaskDao.Save(context.Background(), &NovelTask); err != nil {
		zap.L().Error("failed to save catalogPageTask", zap.Error(err))
	}

	return
}

// 检查是否已经处理过的url
func (t DefaultTaskProcessor) isDuplicatedNovelTask(cpTask *model.NovelTask, collectionName,
	url string, bsonFilter bson.M) (bool /*existence*/, error /*interrupted*/) {
	jsonString, err := common.GetAndSet(context.Background(), url, func() (*string, error) {
		if data, err := dao.FindByMongoFilter(context.Background(), bsonFilter, collectionName, cpTask, &options.FindOneOptions{}); err != nil {
			return nil, err
		} else {
			taskString := convertor.ToString(data)
			if taskString == "" {
				return nil, nil
			}
			return &taskString, nil
		}
	})

	if err != nil || jsonString == nil {
		return false, err
	}
	if err = json.Unmarshal([]byte(*jsonString), cpTask); err != nil {
		return false, err
	}

	return cpTask.GetStatus() == common.TaskStatusFinished, err
}

// 检查是否已经处理过的url
func (t DefaultTaskProcessor) isDuplicatedChapterTask(cpTask *model.ChapterTask, collectionName,
	url string, bsonFilter bson.M) (bool /*existence*/, error /*interrupted*/) {
	jsonString, err := common.GetAndSet(context.Background(), url, func() (*string, error) {
		if data, err := dao.FindByMongoFilter(context.Background(), bsonFilter, collectionName,
			cpTask, &options.FindOneOptions{}); err != nil {
			return nil, err
		} else {
			taskString := convertor.ToString(data)
			if taskString == "" {
				return nil, nil
			}
			return &taskString, nil
		}
	})

	if err != nil || jsonString == nil {
		return false, err
	}
	if err = json.Unmarshal([]byte(*jsonString), cpTask); err != nil {
		return false, err
	}

	return cpTask.GetStatus() == common.TaskStatusFinished, err
}

// 检查是否已经处理过的url
func (t DefaultTaskProcessor) isDuplicatedCatalogPageTask(cpTask *model.CatalogPageTask, collectionName,
	url string, bsonFilter bson.M) (bool /*existence*/, error /*interrupted*/) {
	jsonString, err := common.GetAndSet(context.Background(), url, func() (*string, error) {
		if data, err := dao.FindByMongoFilter(context.Background(), bsonFilter, collectionName, cpTask, &options.FindOneOptions{}); err != nil {
			return nil, err
		} else {
			taskString := convertor.ToString(data)
			if taskString == "" {
				return nil, nil
			}
			return &taskString, nil
		}
	})

	if err != nil || jsonString == nil {
		return false, err
	}
	if err = json.Unmarshal([]byte(*jsonString), cpTask); err != nil {
		return false, err
	}
	return (model.Resource(cpTask)).GetStatus() == common.TaskStatusFinished, err
}

func (d DefaultTaskProcessor) HandleChapterTask(jsonData string) interface{} {
	var chapterTask model.ChapterTask
	var err error

	if err = json.Unmarshal([]byte(jsonData), &chapterTask); err != nil {
		zap.L().Error("json decode", zap.Error(err), zap.String("data", jsonData))
		return nil
	}
	zap.L().Info("novel task", zap.String("json", jsonData))

	//check if page url is duplicated
	exists, err := d.isDuplicatedChapterTask(&chapterTask,
		common.CollectionChapterTask,
		chapterTask.Url,
		bson.M{
			common.ColumnNovelId: chapterTask.NovelId, //catalogPageTask.catalogId
			common.ColumnUrl:     chapterTask.Url,     //catalogPageTask.Url
		})
	if err != nil {
		zap.L().Warn("error occurs", zap.Error(err))
		return nil
	}
	if exists {
		zap.L().Warn("duplicated novel to crawl", zap.String("jsonData", jsonData))
		return nil
	}

	downloader := GetSiteCrawler(chapterTask.SiteName)
	if downloader == nil {
		zap.L().Error("site downloader not found", zap.String("SiteName", chapterTask.SiteName))
		return nil
	}

	//check if it exists in db
	var existingTask *model.ChapterTask
	if existingTask, err = dao.ChapterTaskDao.FindByUrl(context.Background(), chapterTask.Url); err != nil {
		zap.L().Error("failed to retrieve chapter page task", zap.String("jsonData", jsonData), zap.Error(err))
		return nil
	}

	currentTime := time.Now()
	if err = downloader.HandleChapterPage(context.Background(), &chapterTask); err != nil {
		//save failed, update the status
		if existingTask != nil {
			if err = convertor.CopyProperties(&chapterTask, existingTask); err != nil {
				zap.L().Error("failed to copy properties of catalog page task", zap.Error(err))
				return nil
			}
			//如果之前重试过，重试次数加1
			if chapterTask.Status == common.TaskStatusFailed ||
				chapterTask.Status == common.TaskStatusRetryFailed {
				chapterTask.Retries++
				chapterTask.Status = common.TaskStatusRetryFailed
			}
		} else {
			chapterTask.Status = common.TaskStatusFailed
		}
		chapterTask.LastUpdated = &currentTime
	} else {
		//已经处理过，记录该url
		chapterTask.Status = common.TaskStatusFinished
		chapterTask.CreatedDate = &currentTime
	}

	if _, err = dao.ChapterTaskDao.Save(context.Background(), &chapterTask); err != nil {
		zap.L().Error("failed to save chapterTask", zap.Error(err))
	}

	return nil
}
