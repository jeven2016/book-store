package dao

var CatalogDao CatalogInterface
var SiteDao SiteInterface
var CatalogPageTaskDao CatalogPageTaskInterface
var NovelTaskDao NovelTaskInterface
var NovelDao NovelInterface
var ChapterDao ChapterInterface
var ChapterTaskDao ChapterTaskInterface

func InitDao() {
	CatalogDao = &CatalogDaoImpl{}
	SiteDao = &SiteDaoImpl{}
	CatalogPageTaskDao = &CatalogPageTaskDaoImpl{}
	NovelTaskDao = &NovelTaskDaoImpl{}
	NovelDao = &NovelDaoImpl{}
	ChapterDao = &ChapterDaoImpl{}
	ChapterTaskDao = &ChapterTaskDaoImpl{}
}
