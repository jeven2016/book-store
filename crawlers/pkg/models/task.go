package models

import (
	"crawlers/pkg/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Status int

const (
	StatusNotStared Status = iota + 1
	StatusProcessing
	StatusFinished
	StatusFailed
)

// Resource 资源，标记网站上需要下载的的资源，供爬虫使用
type Resource interface {
	ResourceType() common.CrawlerResourceType
}

// OperationDate 操作日期
type OperationDate struct {
	CreatedDate *time.Time `bson:"createdDate" json:"createdDate"`
	LastUpdated *time.Time `bson:"lastUpdated" json:"lastUpdated"`
}

// SiteTask 站点
type SiteTask struct {
	Id          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name        string                 `bson:"name" json:"name" binding:"required"`
	Description string                 `bson:"description" json:"description"`
	Url         string                 `bson:"url" json:"url" binding:"required"`
	Attributes  map[string]interface{} `bson:"attributes" json:"attributes"`
	Status      Status                 `bson:"status" json:"status"`

	OperationDate
}

func (s *SiteTask) ResourceType() common.CrawlerResourceType {
	return common.SiteResourceType
}

// CatalogTask 分类目录
type CatalogTask struct {
	// 添加omitempty，当为空时，mongo driver会自动生成
	Id              primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	ParentCatalogId primitive.ObjectID     `bson:"parentCatalogId,omitempty" json:"parentCatalogId"`
	Name            string                 `bson:"name" json:"name" binding:"required"`
	Url             string                 `bson:"url" json:"url" binding:"required"`
	Description     string                 `bson:"description" json:"description"`
	Attributes      map[string]interface{} `bson:"attributes" json:"attributes"`
	Status          Status                 `bson:"status" json:"status"`
	SiteName        string                 `bson:"siteName" json:"siteName"` //便于后续的日志输出
	OperationDate
}

func (s *CatalogTask) ResourceType() common.CrawlerResourceType {
	return common.CatalogResourceType
}

type CatalogPageTask struct {
	// 添加omitempty，当为空时，mongo driver会自动生成
	Id         primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	CatalogId  primitive.ObjectID     `json:"catalogId" bson:"catalogId" binding:"required"`
	Url        string                 `bson:"url" json:"url" binding:"required" binding:"required"`
	Attributes map[string]interface{} `bson:"attributes" json:"attributes"`
	Status     Status                 `bson:"status" json:"status"`
	SiteName   string                 `bson:"siteName" json:"siteName"`
	OperationDate
}

func (s *CatalogPageTask) ResourceType() common.CrawlerResourceType {
	return common.CatalogPageResourceType
}

type NovelTask struct {
	Id          primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name        string                 `bson:"name" json:"name" binding:"required"`
	CatalogId   primitive.ObjectID     `bson:"catalogId,omitempty" json:"catalogId" binding:"required"`
	Url         string                 `bson:"url" json:"url" binding:"required"`
	Content     string                 `bson:"content,omitempty" json:"content"`
	HasChapters bool                   `bson:"hasChapters,omitempty" json:"hasChapters"`
	Attributes  map[string]interface{} `bson:"attributes" json:"attributes"`
	Status      Status                 `bson:"status" json:"status"`
	SiteName    string                 `bson:"siteName" json:"siteName"`
	OperationDate
}

func (s *NovelTask) ResourceType() common.CrawlerResourceType {
	return common.ArticleResourceType
}

type NovelPageTask struct {
	// 添加omitempty，当为空时，mongo driver会自动生成
	Id         primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	ArticleId  primitive.ObjectID     `bson:"articleId,omitempty" json:"articleId" binding:"required"`
	Url        string                 `bson:"url,omitempty" json:"url" binding:"required"`
	Attributes map[string]interface{} `bson:"attributes" json:"attributes"`
	Status     Status                 `bson:"status" json:"status"`
	SiteName   string                 `bson:"siteName" json:"siteName"`
	OperationDate
}

func (s *NovelPageTask) ResourceType() common.CrawlerResourceType {
	return common.ArticlePageResourceType
}

type ChapterTask struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	ArticleId primitive.ObjectID `bson:"articleId,omitempty" json:"articleId"`
	Url       string             `bson:"url" json:"url"`
	Status    Status             `bson:"status" json:"status"`
	SiteName  string             `bson:"siteName" json:"siteName"`
	OperationDate
}

func (s *ChapterTask) ResourceType() common.CrawlerResourceType {
	return common.ChapterResourceType
}

type ChapterPageTask struct {
	// 添加omitempty，当为空时，mongo driver会自动生成
	Id         primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	ArticleId  primitive.ObjectID     `bson:"articleId,omitempty" json:"articleId"`
	Url        string                 `bson:"url,omitempty" json:"url"`
	Attributes map[string]interface{} `bson:"attributes" json:"attributes"`
	Status     Status                 `bson:"status" json:"status"`
	SiteName   string                 `bson:"siteName" json:"siteName"`
	OperationDate
}

func (s *ChapterPageTask) ResourceType() common.CrawlerResourceType {
	return common.ChapterPageResourceType
}
