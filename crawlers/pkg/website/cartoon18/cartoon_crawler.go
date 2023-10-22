package cartoon18

import (
	"context"
	"crawlers/pkg/common"
	"crawlers/pkg/dao"
	"crawlers/pkg/model"
	"crawlers/pkg/model/entity"
	"fmt"
	"github.com/go-creed/sat"
	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type CartoonCrawler struct {
	redis       *common.Redis
	mongoClient *common.MongoClient
	colly       *colly.Collector
	siteCfg     *common.SiteConfig
	client      *resty.Client
	zhConvertor sat.Dicter
}

func NewCartoonCrawler() *CartoonCrawler {
	sys := common.GetSystem()
	cfg := common.GetSiteConfig(common.Cartoon18)
	if cfg == nil {
		sys.Log.Sugar().Warn("Could not find site config", zap.String("siteName", common.SiteNsf))
	}

	return &CartoonCrawler{
		redis:       sys.RedisClient,
		mongoClient: sys.MongoClient,
		colly:       common.NewCollector(sys.Log),
		siteCfg:     cfg,
		client:      resty.New(),
		zhConvertor: sat.DefaultDict(),
	}
}

func (c CartoonCrawler) CrawlHomePage(ctx context.Context, url string) error {
	//TODO implement me
	panic("implement me")
}

func (c CartoonCrawler) CrawlCatalogPage(ctx context.Context, catalogPageTask *model.CatalogPageTask) ([]model.NovelTask, error) {
	zap.L().Info("Got CatalogPageTask message", zap.String("url", catalogPageTask.Url))
	var novelTasks []model.NovelTask
	cly := c.colly.Clone()
	cly.OnHTML(".card .lines.lines-2 a.visited", func(element *colly.HTMLElement) {
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

func (c CartoonCrawler) CrawlNovelPage(ctx context.Context, novelTask *model.NovelTask, skipSaveIfPresent bool) ([]model.ChapterTask, error) {
	zap.L().Info("Got novel message", zap.String("url", novelTask.Url))
	var createdTime = time.Now()
	var novel = entity.Novel{Attributes: make(map[string]interface{}), CreatedTime: &createdTime}
	var chpTasks []model.ChapterTask
	cly := c.colly.Clone()
	//获取名称
	cly.OnHTML(".title.py-1", func(element *colly.HTMLElement) {
		name := c.zhConvertor.Read(element.Text)
		name = strings.Split(name, "\t\n\t\t")[0]
		name = strings.ReplaceAll(name, "\n\t", "")
		name = strings.TrimSpace(name)

		if strings.Contains(name, "/") {
			name = strings.Split(name, "/")[1]
		}
		novel.Name = name
	})

	//获取每一页上面的chapter内容
	cly.OnHTML(".btn.btn-info.mr-2.mb-2", func(a *colly.HTMLElement) {
		chapterName := c.zhConvertor.Read(a.Text)
		chpTask := model.ChapterTask{
			Name:     chapterName,
			SiteName: novelTask.SiteName,
			Url:      common.BuildUrl(novelTask.Url, a.Attr("href")),
		}
		chpTasks = append(chpTasks, chpTask)
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

	//create directory
	if novelDir, ok := c.siteCfg.Attributes["directory"]; ok {
		if err = os.MkdirAll(filepath.Join(novelDir, novel.Name), 0755); err != nil {
			return chpTasks, err
		}
	}

	return chpTasks, nil
}

func (c CartoonCrawler) CrawlChapterPage(ctx context.Context, chapterTask *model.ChapterTask, skipSaveIfPresent bool) error {
	var err error
	var client *resty.Client
	var novel *entity.Novel

	cly := c.colly.Clone()
	zap.L().Info("Got chapter message", zap.String("url", chapterTask.Url))

	if novel, err = dao.NovelDao.FindById(ctx, chapterTask.NovelId); err != nil {
		return err
	}

	//以novel名称为根目录，chapter目录为子目录
	var chapterDir string
	if novelDir, ok := c.siteCfg.Attributes["directory"]; ok {
		chapterDir = filepath.Join(novelDir, novel.Name, chapterTask.Name)
		if err = os.MkdirAll(chapterDir, 0755); err != nil {
			return err
		}
	}

	if chapterDir == "" {
		return fmt.Errorf("no chapter directory specified %v", c.siteCfg.Attributes["directory"])
	}

	var i = 1
	cly.OnHTML(".cartoon-image", func(img *colly.HTMLElement) {
		if err != nil {
			return
		}

		if i%100 == 0 {
			time.Sleep(4 * time.Second)
		}

		picUrl := img.Attr("data-src")
		client, err = common.GetRestyClient(picUrl)
		if err != nil {
			return
		}

		var fileFormat = ".webp"
		if !strings.Contains(picUrl, ".webp") {
			fileFormat = ".jpg"
		}

		destFile := filepath.Join(chapterDir, fmt.Sprintf("%04d", i)+fileFormat)
		i++

		if _, err = client.R().SetOutput(destFile).Get(picUrl); err != nil {
			zap.L().Error("download picture error", zap.String("url", picUrl), zap.Error(err))
			return
		} else {
			zap.L().Info("picture downloaded", zap.String("url", picUrl), zap.String("localFile", destFile))
		}
	})
	if err = cly.Visit(chapterTask.Url); err != nil {
		return err
	}
	return nil
}
