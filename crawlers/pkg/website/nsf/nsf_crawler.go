package nfs

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/dao"
	"crawlers/pkg/model"
	"crawlers/pkg/model/entity"
	"github.com/chromedp/chromedp"
	"github.com/go-creed/sat"
	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"strings"
	"time"
)

type NsfCrawler struct {
	redis       *common.Redis
	mongoClient *common.MongoClient
	colly       *colly.Collector
	siteCfg     *common.SiteConfig
	client      *resty.Client
	zhConvertor sat.Dicter
}

func NewNsfCrawler() *NsfCrawler {
	sys := common.GetSystem()
	cfg := common.GetSiteConfig(common.SiteNsf)
	if cfg == nil {
		zap.L().Sugar().Warn("Could not find site config", zap.String("siteName", common.SiteNsf))
	}

	return &NsfCrawler{
		redis:       sys.RedisClient,
		mongoClient: sys.MongoClient,
		colly:       common.NewCollector(zap.L()),
		siteCfg:     cfg,
		client:      resty.New(),
		zhConvertor: sat.DefaultDict(),
	}
}

var removeTexts = []string{
	"<p>更*多`精;彩'小*说'尽|在'ｗ'ｗ'ｗ．''Ｂ'．'Ｅ'第&amp;#*站</p><p>\");</p>",
	"<p>ThisfilewassavedusingUNREGISTEREDversionofChmDecompiler.</p><p>DownloadChmDecompilerat:（结尾英文忽略即可）</p>",
	"<p>##</p><p>ThefilewassavedusingTrialversionofChmDecompiler.</p><p>DownloadChmDecompilerfrom:（结尾英文忽略即可）</p>",
}

// CrawlCatalogPage 解析每一页
func (n *NsfCrawler) CrawlCatalogPage(ctx context.Context, catalogPageTask *model.CatalogPageTask) ([]model.NovelTask, error) {
	zap.L().Info("Got CatalogPageTask message", zap.String("url", catalogPageTask.Url))
	var novelTasks []model.NovelTask
	cly := n.colly.Clone()
	cly.OnHTML(".CGsectionTwo-right-content-unit .title", func(element *colly.HTMLElement) {
		href := element.Attr("href")
		novelUrl := common.BuildUrl(catalogPageTask.Url, href)
		novelTasks = append(novelTasks, model.NovelTask{
			Url:      novelUrl,
			SiteName: catalogPageTask.SiteName,
		})
	})

	if err := cly.Visit(catalogPageTask.Url); err != nil {
		return nil, err
	}
	return novelTasks, nil
}

// CrawlNovelPage 解析具体的Novel
func (n *NsfCrawler) CrawlNovelPage(ctx context.Context, novelTask *model.NovelTask, skipSaveIfPresent bool) ([]model.ChapterTask, error) {
	zap.L().Info("Got novel message", zap.String("url", novelTask.Url))
	var createdTime = time.Now()
	var novel = entity.Novel{Attributes: make(map[string]interface{}), CreatedTime: &createdTime}
	var chpTasks []model.ChapterTask
	cly := n.colly.Clone()
	//获取名称
	cly.OnHTML(".title", func(element *colly.HTMLElement) {
		novel.Name = n.zhConvertor.Read(element.Text)
		novelTask.Name = novel.Name
	})

	//获取作者
	cly.OnHTML(".author .b", func(element *colly.HTMLElement) {
		novel.Attributes[common.AttrAuthor] = n.zhConvertor.Read(element.Text)
	})

	//获取描述
	cly.OnHTML(".BGsectionTwo-bottom", func(element *colly.HTMLElement) {
		desc := n.zhConvertor.Read(element.Text)
		novel.Description = strings.TrimSpace(desc)
	})

	//获取“全部章节”按钮
	cly.OnHTML(".BGsectionOne-bottom li:nth-of-type(2) a", func(element *colly.HTMLElement) {
		allChaptersLink := common.BuildUrl(novelTask.Url, element.Attr("href"))
		if allChaptersLink == "" {
			zap.L().Warn("No chapters found", zap.String("novelUrl", novelTask.Url))
			return
		}

		if err := cly.Request("GET", allChaptersLink, nil, colly.NewContext(), nil); err != nil {
			zap.L().Error("failed to parser chapters link",
				zap.String("novelUrl", novelTask.Url), zap.Error(err))
			return
		}
	})

	//获取每一页上面的chapter内容
	cly.OnHTML(".BCsectionTwo-top-chapter a", func(a *colly.HTMLElement) {
		chapterName := n.zhConvertor.Read(a.Text)
		chpTask := model.ChapterTask{
			Name:     chapterName,
			SiteName: novelTask.SiteName,
			Url:      common.BuildUrl(novelTask.Url, a.Attr("href")),
		}
		chpTasks = append(chpTasks, chpTask)
	})

	//解析完当前页面，解析下一页
	cly.OnHTML(".CGsectionTwo-right-bottom-btn #next", func(nextBtn *colly.HTMLElement) {
		nextPageUrl := common.BuildUrl(novelTask.Url, nextBtn.Attr("href"))
		if err := cly.Visit(nextPageUrl); err != nil {
			zap.L().Error("error occurred while visiting the next page", zap.String("nextPageUrl", nextPageUrl),
				zap.String("novelName", novelTask.Name))
			return
		}
	})

	if err := cly.Visit(novelTask.Url); err != nil {
		return nil, err
	}

	var novelId *primitive.ObjectID
	var err error

	if novelId, err = dao.NovelDao.FindIdByName(ctx, novel.Name); err != nil {
		return nil, err
	}

	if !skipSaveIfPresent || novelId == nil {
		//保存novel
		novel.HasChapters = len(chpTasks) > 0
		if novelId != nil {
			novel.Id = *novelId
		}
		if novelId, err = dao.NovelDao.Save(ctx, &novel); err != nil {
			return nil, err
		}
	}

	if novelId != nil {
		for i := 0; i < len(chpTasks); i++ {
			chpTasks[i].NovelId = *novelId
			chpTasks[i].Order = i + 1
		}
	}

	if len(chpTasks) == 0 {
		zap.L().Error("no chapters found for novel", zap.String("novelName", novel.Name))
	} else {
		zap.L().Info("number of chapters found for novel", zap.String("novelName", novel.Name),
			zap.Int("number", len(chpTasks)))
	}
	return chpTasks, nil
}

func (n *NsfCrawler) CrawlHomePage(ctx context.Context, url string) error {
	//TODO implement me
	panic("implement me")
}

func (n *NsfCrawler) CrawlChapterPage(ctx context.Context, chapterTask *model.ChapterTask, skipSaveIfPresent bool) (err error) {
	zap.L().Info("Got chapter message", zap.String("url", chapterTask.Url))
	var createdTime = time.Now()

	chromeCtx, cleanFunc := common.OpenChrome(context.Background())
	defer cleanFunc()

	var text string //content of the chapter
	err = chromedp.Run(chromeCtx,
		chromedp.Navigate(chapterTask.Url),
		//chromedp.WaitNotPresent("//p[contains(text(),'内容未加载完成')]", chromedp.BySearch),
		chromedp.InnerHTML("//div[@class='RBGsectionThree-content']", &text, chromedp.BySearch),
	)
	if err != nil {
		return
	}

	var chapterId *primitive.ObjectID

	// for chapter
	existingChapter, err := dao.ChapterDao.FindByName(ctx, chapterTask.Name)
	if err != nil {
		return
	}
	if existingChapter != nil {
		chapterId = &existingChapter.Id
		existingChapter.NovelId = chapterTask.NovelId
		existingChapter.Order = chapterTask.Order
	} else {
		//create one
		existingChapter = &entity.Chapter{
			NovelId:     chapterTask.NovelId,
			Name:        chapterTask.Name,
			Order:       chapterTask.Order,
			CreatedTime: &createdTime,
		}
	}

	//todo
	if !skipSaveIfPresent || chapterId == nil || (*chapterId).IsZero() {
		if chapterId, err = dao.ChapterDao.Save(ctx, existingChapter); err != nil {
			return
		}
	}

	//for content
	//page for chapters, need an enhancement
	existingContent, err := dao.ContentDao.FindByParentIdAndPage(ctx, chapterId, 0)
	if err != nil {
		return
	}

	for _, txt := range removeTexts {
		text = strings.ReplaceAll(text, txt, "")
	}

	if existingContent != nil {
		existingContent.Content = text
	} else {
		//create one
		existingContent = &entity.Content{
			ParentId:    *chapterId,
			ParentType:  common.ParentTypeChapter,
			Content:     text,
			CreatedTime: &createdTime,
		}
	}

	if !skipSaveIfPresent || existingContent == nil || existingContent.Id.IsZero() {
		_, err = dao.ContentDao.Save(ctx, existingContent)
	}
	return
}
