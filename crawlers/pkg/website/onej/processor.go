package onej

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/models"
	"go.uber.org/zap"
)

type OnejProcessor struct {
	logger  *zap.Logger
	siteCfg *common.SiteConfig
}

func NewSiteOnejProcessor() *OnejProcessor {
	sys := common.GetSystem()
	cfg, err := common.GetSiteConfig(common.SiteOneJ)
	if err != nil {
		sys.Log.Sugar().Error("Could not find site config", zap.Error(err))
		return nil
	}

	return &OnejProcessor{
		logger:  sys.Log,
		siteCfg: cfg,
	}
}

func (p *OnejProcessor) ProcessPageUrls(originPageUrl string) ([]string, error) {
	return common.GenPageUrls(p.siteCfg.RegexSettings.ParsePageRegex, originPageUrl,
		p.siteCfg.RegexSettings.PagePrefix)
}

func (p *OnejProcessor) SaveCatalogPageUrl(catalogPageMsg *models.CatalogPageTask) error {
	return p.insertDocument(common.CatalogPageTaskCollection, catalogPageMsg)
}

func (p *OnejProcessor) insertDocument(collectionName string, msg models.Resource) error {
	col := common.GetSystem().GetCollection(collectionName)
	if col == nil {
		zap.L().Error("no collection found", zap.String("siteKey", "TODO"),
			zap.String("collectionName", collectionName))
		return common.CollectionNotFoundError
	}
	if _, err := col.InsertOne(context.Background(), msg); err != nil {
		return err
	}
	return nil
}

func (p *OnejProcessor) SaveNovelPageUrl(novelMsg *models.NovelTask) error {
	return p.insertDocument(common.NovelCollection, novelMsg)
}
