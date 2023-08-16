package common

import "errors"

const (
	RedisStreamDataVar = "data"
	SiteOneJ           = "onej"
	SiteYzs8           = "yzs8"
	CollyMaxRetries    = 3
)

// db column
const (
	ColumId         = "_id"
	ColumnName      = "name"
	ColumnCatalogId = "catalogId"
	SiteName        = "siteName"
	NovelLink       = "novelLink"
	ColumnUrl       = "url"
)

// db collection
const (
	SiteCollection            = "site"
	CatalogCollection         = "catalog"
	NovelCollection           = "novel"
	CatalogPageTaskCollection = "catalogPageTask"
)

var ConfigFiles = []string{"/etc/crawlers.yaml", "./crawlers.yaml"}
var SupportedCatalogs = []string{SiteOneJ, SiteYzs8}

var CollectionNotFoundError = errors.New("collection not found")

type CrawlerResourceType int

const (
	SiteResourceType CrawlerResourceType = iota + 1
	CatalogResourceType
	CatalogPageResourceType
	ArticleResourceType
	ArticlePageResourceType
	ChapterResourceType
	ChapterPageResourceType
)

// CrawlerType 抓取资源类型
type CrawlerType int

const (
	BtCrawlerType = iota + 1
	ComicCrawlerType
	NovelCrawlerType
)

// cache key prefix

const (
	SiteKeyExistsPrefix    = "site:exists"
	CatalogKeyExistsPrefix = "catalog:exists"
)

var ErrDecodingDocument = errors.New("document retrieved without decoding process")
