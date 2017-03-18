package top

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TopSetTaskResult struct {
	app.Result
	Top *Top `json:"top,omitempty"`
}

type TopSetTask struct {
	app.Task
	Id           int64       `json:"id"`
	ItemId       interface{} `json:"itemId"`       //推荐项ID
	ItemType     interface{} `json:"itemType"`     //推荐项类型
	Options      interface{} `json:"options"`      //其他选项
	Tags         interface{} `json:"tags"`         //搜索标签
	CityId       interface{} `json:"cityId"`       //城市ID
	CityPath     interface{} `json:"cityPath"`     //城市路径
	ClassifyId   interface{} `json:"classifyId"`   //分类ID
	ClassifyPath interface{} `json:"classifyPath"` //分类路径
	GroupId      interface{} `json:"groupId"`      //分组ID
	Result       TopSetTaskResult
}

func (task *TopSetTask) GetResult() interface{} {
	return &task.Result
}

func (task *TopSetTask) GetInhertType() string {
	return "top"
}

func (task *TopSetTask) GetClientName() string {
	return "Top.Create"
}
