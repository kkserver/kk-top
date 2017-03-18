package top

import (
	"database/sql"
	"github.com/kkserver/kk-lib/kk"
	"github.com/kkserver/kk-lib/kk/app"
	"github.com/kkserver/kk-lib/kk/app/remote"
	"time"
)

type Top struct {
	Id           int64  `json:"id"`
	Alias        string `json:"alias"`        //别名
	ItemId       int64  `json:"itemId"`       //推荐项ID
	ItemType     string `json:"itemType"`     //推荐项类型
	Tags         string `json:"tags"`         //搜索标签
	CityId       int64  `json:"cityId"`       //城市ID
	CityPath     string `json:"cityPath"`     //城市路径
	ClassifyId   int64  `json:"classifyId"`   //分类ID
	ClassifyPath string `json:"classifyPath"` //分类路径
	GroupId      int64  `json:"groupId"`      //分组ID
	Options      string `json:"options"`      //其他选项
	Oid          int64  `json:"oid"`
}

type ITopApp interface {
	app.IApp
	GetDB() (*sql.DB, error)
	GetPrefix() string
	GetTopTable() *kk.DBTable
}

type TopApp struct {
	app.App

	DB *app.DBConfig

	Remote *remote.Service

	Top      *TopService
	TopTable kk.DBTable
}

const twepoch = int64(1424016000000)

func milliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func NewOid() int64 {
	return milliseconds() - twepoch
}

func (C *TopApp) GetDB() (*sql.DB, error) {
	return C.DB.Get(C)
}

func (C *TopApp) GetPrefix() string {
	return C.DB.Prefix
}

func (C *TopApp) GetTopTable() *kk.DBTable {
	return &C.TopTable
}
