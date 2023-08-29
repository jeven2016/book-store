package common

import "errors"

const (
	RedisStreamDataVar = "data"
	SiteOneJ           = "onej"
	SiteYzs8           = "yzs8"
	SiteNsf            = "nsf"

	CollyMaxRetries = 3

	BulkBatchSize = 10 // 批量保存的文档数量

	ParentTypeChapter = "chapter"
	ParentTypeNovel   = "novel"
)

// db column
const (
	ColumId         = "_id"
	ColumnName      = "name"
	ColumnCatalogId = "catalogId"
	ColumnUrl       = "url"
	ColumnSiteName  = "siteName"
	ColumnNovelId   = "novelId"

	AttrAuthor = "author"
)

// db collection
const (
	CollectionSite            = "site"
	CollectionCatalog         = "catalog"
	CollectionNovel           = "novel"
	CollectionNovelTask       = "novelTask"
	CollectionChapter         = "chapter"
	CollectionChapterTask     = "chapterTask"
	CollectionCatalogPageTask = "catalogPageTask"
	CollectionContent         = "content"
)

var ConfigFiles = []string{"/etc/crawlers.yaml", "./crawlers.yaml"}
var SupportedCatalogs = []string{SiteOneJ, SiteYzs8}

var CollectionNotFoundError = errors.New("collection not found")
var ErrNonValueProvided = errors.New("null value provided for a specified key")

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

type TaskStatus int

const (
	TaskStatusNotStared TaskStatus = iota + 1
	TaskStatusProcessing
	TaskStatusFinished
	TaskStatusFailed
	TaskStatusRetryFailed
)

var ErrDecodingDocument = errors.New("document retrieved without decoding process")
var ErrDuplicatedDocument = errors.New("document is duplicated")
var ErrDocumentIdExists = errors.New("document's ID exists")
