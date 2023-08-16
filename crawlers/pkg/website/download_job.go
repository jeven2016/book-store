package website

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/models"
	"crawlers/pkg/website/onej"
)

type SiteProcessor interface {
	ProcessPageUrls(originPageUrl string) ([]string, error)
	SaveCatalogPageUrl(catalogPageMsg *models.CatalogPageTask) error
	SaveNovelPageUrl(catalogPageMsg *models.NovelTask) error
}

type SiteCrawler interface {
	HandleCatalogPage(ctx context.Context, catalogPageMsg *models.CatalogPageTask) ([]models.NovelTask, error)
	HandleNovelPage(ctx context.Context, novelPageMsg *models.NovelTask) ([]models.ChapterTask, error)
	DownloadHomePage(ctx context.Context, url string) error
}

var siteDownloaderMap = make(map[string]SiteCrawler)
var siteProcessorMap = make(map[string]SiteProcessor)

func InitJobHandlers() {
	siteDownloaderMap[common.SiteOneJ] = onej.NewSiteOnej()

	siteProcessorMap[common.SiteOneJ] = onej.NewSiteOnejProcessor()
}

func GetSiteCrawler(catalogName string) SiteCrawler {
	return siteDownloaderMap[catalogName]
}

func GetSiteProcessor(catalogName string) SiteProcessor {
	return siteProcessorMap[catalogName]
}
