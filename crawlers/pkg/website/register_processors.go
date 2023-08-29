package website

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/model"
	nfs "crawlers/pkg/website/nsf"
	"crawlers/pkg/website/onej"
)

type SiteCrawler interface {
	HandleHomePage(ctx context.Context, url string) error
	HandleCatalogPage(ctx context.Context, catalogPageMsg *model.CatalogPageTask) ([]model.NovelTask, error)
	HandleNovelPage(ctx context.Context, novelPageMsg *model.NovelTask) ([]model.ChapterTask, error)
	HandleChapterPage(ctx context.Context, chapterMsg *model.ChapterTask) error
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
