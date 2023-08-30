package nfs

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/dao"
	"crawlers/pkg/model"
	"crawlers/pkg/model/entity"
	"github.com/chromedp/chromedp"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-creed/sat"
	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly/v2"
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
	cfg, err := common.GetSiteConfig(common.SiteNsf)
	if err != nil {
		sys.Log.Sugar().Error("Could not find site config", zap.Error(err))
	}

	return &NsfCrawler{
		redis:       sys.RedisClient,
		mongoClient: sys.MongoClient,
		colly:       common.NewCollector(sys.Log),
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

// HandleCatalogPage 解析每一页
func (n *NsfCrawler) HandleCatalogPage(ctx context.Context, catalogPageMsg *model.CatalogPageTask) ([]model.NovelTask, error) {
	panic("Not implemented")
}

// HandleNovelPage 解析具体的Novel
func (n *NsfCrawler) HandleNovelPage(ctx context.Context, novelTask *model.NovelTask) ([]model.ChapterTask, error) {
	zap.L().Info("Got novel message", zap.String("url", novelTask.Url))
	var createdTime = time.Now()
	var novel = entity.Novel{Attributes: make(map[string]interface{}), CreatedTime: &createdTime}
	var chapters []*entity.Chapter
	var chpTasks []model.ChapterTask

	//获取名称
	n.colly.OnHTML(".title", func(element *colly.HTMLElement) {
		novel.Name = element.Text
	})

	//获取作者
	n.colly.OnHTML(".author .b", func(element *colly.HTMLElement) {
		novel.Attributes[common.AttrAuthor] = element.Text
	})

	//获取描述
	n.colly.OnHTML(".BGsectionTwo-bottom", func(element *colly.HTMLElement) {
		desc := n.zhConvertor.Read(element.Text)
		novel.Description = strings.TrimSpace(desc)
	})

	//获取“全部章节”按钮
	n.colly.OnHTML(".BGsectionOne-bottom li:nth-of-type(2) a", func(element *colly.HTMLElement) {
		allChaptersLink := common.BuildUrl(novelTask.Url, element.Attr("href"))
		if allChaptersLink == "" {
			zap.L().Warn("No chapters found", zap.String("novelUrl", novelTask.Url))
			return
		}

		if err := n.colly.Request("GET", allChaptersLink, nil, colly.NewContext(), nil); err != nil {
			zap.L().Error("failed to parser chapters link",
				zap.String("novelUrl", novelTask.Url), zap.Error(err))
			return
		}
	})

	//获取每一页上面的chapter内容
	var index = 1
	n.colly.OnHTML(".BCsectionTwo-top-chapter a", func(a *colly.HTMLElement) {
		chapterName := n.zhConvertor.Read(a.Text)
		chp := &entity.Chapter{
			Name:        chapterName,
			Order:       index,
			Description: "",
			CreatedTime: &createdTime,
			UpdatedTime: nil,
		}
		chapters = append(chapters, chp)

		chpTask := model.ChapterTask{
			Name:     chp.Name,
			SiteName: novelTask.SiteName,
			Url:      common.BuildUrl(novelTask.Url, a.Attr("href")),
		}
		chpTasks = append(chpTasks, chpTask)

		index++
	})

	//解析完当前页面，解析下一页
	n.colly.OnHTML("#next", func(nextBtn *colly.HTMLElement) {
		nextPageUrl := common.BuildUrl(novelTask.Url, nextBtn.Attr("href"))
		if err := n.colly.Visit(nextPageUrl); err != nil {
			zap.L().Error("error occurred while visiting the next page", zap.String("nextPageUrl", nextPageUrl),
				zap.String("novelName", novelTask.Name))
			return
		}
	})

	if err := n.colly.Visit(novelTask.Url); err != nil {
		return nil, err
	}
	//保存novel
	novel.HasChapters = index > 0
	novelId, err := dao.NovelDao.Insert(ctx, &novel)
	if err != nil {
		return nil, err
	}
	if index > 0 {
		if err = dao.ChapterDao.BulkInsert(ctx, chapters, novelId); err != nil {
			return nil, err
		}
	}

	slice.ForEach(chpTasks, func(index int, item model.ChapterTask) {
		item.NovelId = *novelId
	})
	return chpTasks, nil
}

func (n *NsfCrawler) HandleHomePage(ctx context.Context, url string) error {
	//TODO implement me
	panic("implement me")
}
func (n *NsfCrawler) HandleChapterPage(ctx context.Context, chapterMsg *model.ChapterTask) (err error) {
	zap.L().Info("Got chapter message", zap.String("url", chapterMsg.Url))
	var createdTime = time.Now()
	var content = &entity.Content{
		ParentType:  common.ParentTypeChapter,
		ParentId:    chapterMsg.Id,
		CreatedTime: &createdTime,
	}
	chromeCtx, cleanFunc := common.OpenChrome(context.Background())
	defer cleanFunc()

	var text string
	err = chromedp.Run(chromeCtx,
		chromedp.Navigate(chapterMsg.Url),
		//chromedp.WaitNotPresent("//p[contains(text(),'内容未加载完成')]", chromedp.BySearch),
		chromedp.InnerHTML("//div[@class='RBGsectionThree-content']", &text, chromedp.BySearch),
	)
	if err != nil {
		return
	}
	for _, txt := range removeTexts {
		text = strings.ReplaceAll(text, txt, "")
	}

	content.Content = text
	_, err = dao.ContentDao.Insert(ctx, content)
	return
}
