package api

import (
	"crawlers/pkg/common"
	"crawlers/pkg/dao"
	"crawlers/pkg/model"
	"crawlers/pkg/model/dto"
	"crawlers/pkg/model/entity"
	"crawlers/pkg/stream"
	"crawlers/pkg/website"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Handler struct {
	sys *common.System
}

func NewHandler() *Handler {
	return &Handler{
		sys: common.GetSystem(),
	}
}

func (h *Handler) HandleCatalogPage(c *gin.Context) {
	var pageReq model.CatalogPageTask
	err := c.ShouldBindJSON(&pageReq)
	if err != nil {
		zap.L().Warn("failed to convert json", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest,
			common.FailsWithMessage(common.ErrCodeUnknown, err.Error()))
		return
	}

	var site *entity.Site
	var hasError bool
	if site, hasError = h.getTaskEntity(c, pageReq.CatalogId); hasError {
		return
	}

	//if multiple pages need to handle
	if sp := website.GetSiteTaskProcessor(site.Name); sp == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, common.Fails(common.ErrCodeUnSupportedCatalog))
		h.sys.Log.Warn("no processor found for this siteKey", zap.String("siteKey", site.Name))
		return
	} else {
		//parse all page urls if page parameter is specified in such format: page=1-5
		urls, err := sp.ParsePageUrls(site.Name, pageReq.Url)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, common.FailsWithParams(common.ErrParsePageUrl, err.Error()))
			h.sys.Log.Warn("failed to process pageUrl",
				zap.String("pageUrl", pageReq.Url), zap.Error(err))
			return
		}

		for _, url := range urls {
			pageMsg := &model.CatalogPageTask{
				SiteName:  site.Name,
				CatalogId: pageReq.CatalogId,
				Url:       url,
				Status:    common.TaskStatusNotStared,
			}
			if err = common.GetSystem().RedisClient.PublishMessage(c, pageMsg, stream.CatalogPageUrlStream); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError,
					common.FailsWithParams(common.ErrPublishMessage, err.Error()))
				h.sys.Log.Warn("failed to publish a message",
					zap.String("pageUrl", pageReq.Url), zap.Error(err))
				return
			}
		}

	}

	c.JSON(http.StatusAccepted, common.SuccessCode(common.ErrCodeTaskSubmitted))
}

func (h *Handler) getTaskEntity(c *gin.Context, catalogId primitive.ObjectID) (site *entity.Site, hasError bool) {
	var err error
	var catalog *entity.Catalog
	catalogStringId := catalogId.Hex()
	siteStringId := catalogId.Hex()
	if catalog, err = dao.CatalogDao.FindById(c, catalogId); err != nil {
		zap.L().Warn("catalog does not exist", zap.String("catalogId", catalogStringId), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, common.FailsWithParams(common.ErrCatalogNotFound, catalogStringId))
		hasError = true
		return
	}
	if site, err = dao.SiteDao.FindById(c, catalog.SiteId); err != nil {
		zap.L().Warn("site does not exist", zap.String("siteId", siteStringId), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, common.FailsWithParams(common.ErrSiteNotFound, siteStringId))
		hasError = true
		return
	}
	return
}

func (h *Handler) HandleNovelPage(c *gin.Context) {
	var novelTask model.NovelTask
	err := c.ShouldBindJSON(&novelTask)
	if err != nil {
		zap.L().Warn("failed to convert json", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest,
			common.FailsWithMessage(common.ErrCodeUnknown, err.Error()))
		return
	}

	var site *entity.Site
	var hasError bool

	if slice.Contain(common.GetConfig().CrawlerSettings.EcludedNovelUrls, novelTask.Url) {
		zap.L().Warn("excluded novel url", zap.String("url", novelTask.Url))
		c.AbortWithStatusJSON(http.StatusBadRequest,
			common.FailsWithMessage(common.ErrExcludedNovel, err.Error()))
		return
	}

	if site, hasError = h.getTaskEntity(c, novelTask.CatalogId); hasError {
		return
	}
	novelTask.Status = common.TaskStatusNotStared
	novelTask.SiteName = site.Name

	if err := common.GetSystem().RedisClient.PublishMessage(c, novelTask, stream.NovelUrlStream); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			common.FailsWithParams(common.ErrPublishMessage, err.Error()))
		h.sys.Log.Warn("failed to publish a message",
			zap.String("pageUrl", novelTask.Url), zap.Error(err))
		return
	}
}

func (h *Handler) CreateSite(c *gin.Context) {
	var site entity.Site
	if err := c.ShouldBindJSON(&site); err != nil {
		//自定义error， https://juejin.cn/post/7015517416608235534
		zap.L().Warn("failed to convert json", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest,
			common.FailsWithMessage(common.ErrCodeUnknown, err.Error()))
		return
	}
	currentTime := time.Now()
	site.CreatedTime = &currentTime

	h.doCreate(c, &dto.CreateRequest{
		Key:           "site",
		Name:          site.Name,
		Entity:        site,
		Collection:    common.CollectionSite,
		RedisCacheKey: common.GenKey(common.SiteKeyExistsPrefix, site.Name),
	})
}

func (h *Handler) CreateCatalog(c *gin.Context) {
	var catalog entity.Catalog
	if err := c.ShouldBindJSON(&catalog); err != nil {
		zap.L().Warn("failed to convert json", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest,
			common.FailsWithMessage(common.ErrCodeUnknown, err.Error()))
		return
	}

	//check if the site exists and cache the result
	exists, err := common.Exists(c, common.GenKey(common.SiteKeyExistsPrefix, catalog.SiteId.Hex()), func() (any, error) {
		return dao.SiteDao.ExistsById(c, catalog.SiteId)
	})
	if err != nil {
		zap.L().Warn("failed to check if any sites exist with this siteId", zap.String("siteId", catalog.SiteId.Hex()),
			zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			common.FailsWithMessage(common.ErrCodeUnknown, err.Error()))
		return
	}
	if !exists {
		zap.L().Warn("no site exists with this siteId", zap.String("siteId", catalog.SiteId.Hex()))
		c.AbortWithStatusJSON(http.StatusBadRequest,
			common.FailsWithParams(common.ErrSiteNotFound, catalog.SiteId.Hex()))
		return
	}

	h.doCreate(c, &dto.CreateRequest{
		Key:           "catalog",
		Name:          catalog.Name,
		Entity:        catalog,
		Collection:    common.CollectionCatalog,
		RedisCacheKey: common.GenKey(common.CatalogKeyExistsPrefix, catalog.Name),
	})
}

func (h *Handler) doCreate(c *gin.Context, req *dto.CreateRequest) {
	col := common.GetSystem().GetCollection(req.Collection)

	exists, err := common.Exists(c, req.RedisCacheKey, func() (any, error) {
		return dao.CatalogDao.ExistsByName(c, req.Name)
	})
	if err != nil {
		zap.L().Warn("failed to check if it exists", zap.Error(err), zap.Any("request", req.Entity))
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			common.FailsWithMessage(common.ErrCodeUnknown, err.Error()))
		return
	}

	if exists {
		zap.L().Warn("it's duplicated to save", zap.Any(req.Key, req.Name), zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest,
			common.FailsWithParams(common.ErrDuplicated, req.Key, req.Name))
		return
	}

	if obj, err := col.InsertOne(c, req.Entity); err != nil {
		zap.L().Warn("failed to save", zap.Error(err), zap.Any(req.Key, req.Name))
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			common.FailsWithMessage(common.ErrCodeUnknown, err.Error()))
		return
	} else {
		c.JSON(http.StatusCreated, obj)
	}
}

func (h *Handler) RunScheduleTask(c *gin.Context) {

}
