package website

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/model"
	nfs "crawlers/pkg/website/nsf"
	"crawlers/pkg/website/onej"
)

type SiteCrawler interface {
	CrawlHomePage(ctx context.Context, url string) error
	CrawlCatalogPage(ctx context.Context, catalogPageMsg *model.CatalogPageTask) ([]model.NovelTask, error)
	CrawlNovelPage(ctx context.Context, novelPageMsg *model.NovelTask, skipSaveIfPresent bool) ([]model.ChapterTask, error)
	CrawlChapterPage(ctx context.Context, chapterMsg *model.ChapterTask, skipSaveIfPresent bool) error
}

var siteCrawlerMap = make(map[string]SiteCrawler)
var siteTaskProcessorMap = make(map[string]TaskProcessor)

func RegisterProcessors() {
	siteCrawlerMap[common.SiteOneJ] = onej.NewSiteOnej()
	siteCrawlerMap[common.SiteNsf] = nfs.NewNsfCrawler()

	defaultTaskProcessor := NewTaskProcessor()
	siteTaskProcessorMap[common.SiteOneJ] = defaultTaskProcessor
	siteTaskProcessorMap[common.SiteNsf] = defaultTaskProcessor
}

func GetSiteCrawler(siteName string) SiteCrawler {
	return siteCrawlerMap[siteName]
}

func GetSiteTaskProcessor(siteName string) TaskProcessor {
	return siteTaskProcessorMap[siteName]
}
