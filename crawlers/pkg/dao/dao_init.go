package dao

var CatalogDao CatalogInterface
var SiteDao SiteInterface

func InitDao() {
	CatalogDao = &CatalogDaoImpl{}
	SiteDao = &SiteDaoImpl{}
}
