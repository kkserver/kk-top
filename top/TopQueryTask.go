package top

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TopQueryCounter struct {
	PageIndex int `json:"p"`
	PageSize  int `json:"size"`
	PageCount int `json:"count"`
	RowCount  int `json:"rowCount"`
}

type TopQueryTaskResult struct {
	app.Result
	Counter *TopQueryCounter `json:"counter,omitempty"`
	Tops    []Top            `json:"tops,omitempty"`
}

type TopQueryTask struct {
	app.Task
	Id             int64  `json:"id"`
	Alias          string `json:"alias"`
	ItemId         int64  `json:"itemId"`   //推荐项ID
	ItemType       string `json:"itemType"` //推荐项类型
	Keyword        string `json:"q"`
	CityId         int64  `json:"cityId"`       //城市ID
	CityPrefix     string `json:"cityPath"`     //城市路径
	ClassifyId     int64  `json:"classifyId"`   //分类ID
	ClassifyPrefix string `json:"classifyPath"` //分类路径
	GroupId        int64  `json:"groupId"`      //分组ID
	PageIndex      int    `json:"p"`
	PageSize       int    `json:"size"`
	Counter        bool   `json:"counter"`
	Result         TopQueryTaskResult
}

func (task *TopQueryTask) GetResult() interface{} {
	return &task.Result
}

func (task *TopQueryTask) GetInhertType() string {
	return "top"
}

func (task *TopQueryTask) GetClientName() string {
	return "Top.Query"
}
