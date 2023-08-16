package models

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CatalogDoc struct {
	// 添加omitempty，当为空时，mongo driver会自动生成
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ParentId     primitive.ObjectID `bson:"parentId,omitempty" json:"parentId"`
	Name         string             `bson:"name" json:"name"`
	Order        int32              `bson:"order" json:"order"`
	ArticleCount int32              `bson:"articleCount" json:"articleCount"`
	Description  string             `bson:"description" json:"description"`
	CreateDate   time.Time          `bson:"createDate" json:"createDate"`
	LastUpdate   time.Time          `bson:"lastUpdate" json:"lastUpdate"`
}

type Article struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	CatalogId  primitive.ObjectID `bson:"catalogId,omitempty" json:"catalogId"`
	Content    string             `bson:"content,omitempty" json:"content"`
	CreateDate time.Time          `bson:"createDate" json:"createDate"`
}

type ArticlePage struct {
	Name    string
	Url     string
	PageUrl string
}

// MarshalBinary 实现该方法，以便对象可以序列化成子字符串保存到redis
func (a *ArticlePage) MarshalBinary() ([]byte, error) {
	bytes, err := json.Marshal(a)
	return bytes, err
}

func (a *ArticlePage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, a)
}

// CatalogDocMap Deprecated 当数据库中保存的是string不是objectId时，使用该类去解析
type CatalogDocMap struct {
	// 添加omitempty，当为空时，mongo driver会自动生成
	Id           any       `bson:"_id,omitempty" json:"id"`
	ParentId     any       `bson:"parentId,omitempty" json:"parentId"`
	Name         string    `bson:"name" json:"name"`
	Order        int32     `bson:"order" json:"order"`
	ArticleCount int32     `bson:"articleCount" json:"articleCount"`
	Description  string    `bson:"description" json:"description"`
	CreateDate   time.Time `bson:"createDate" json:"createDate"`
	LastUpdate   time.Time `bson:"lastUpdate" json:"lastUpdate"`
}

type ArticleMap struct {
	Id         any       `bson:"_id,omitempty" json:"id"`
	Name       string    `bson:"name" json:"name"`
	CatalogId  any       `bson:"catalogId,omitempty" json:"catalogId"`
	Content    string    `bson:"content,omitempty" json:"content"`
	CreateDate time.Time `bson:"createDate" json:"createDate"`
}
