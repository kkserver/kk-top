package top

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TopCreateTaskResult struct {
	app.Result
	Top *Top `json:"top,omitempty"`
}

type TopCreateTask struct {
	app.Task
	Alias        string      `json:"alias"`        //别名
	ItemId       int64       `json:"itemId"`       //推荐项ID
	ItemType     string      `json:"itemType"`     //推荐项类型
	Tags         string      `json:"tags"`         //搜索标签
	CityId       int64       `json:"cityId"`       //城市ID
	CityPath     string      `json:"cityPath"`     //城市路径
	ClassifyId   int64       `json:"classifyId"`   //分类ID
	ClassifyPath string      `json:"classifyPath"` //分类路径
	GroupId      int64       `json:"groupId"`      //分组ID
	Options      interface{} `json:"options"`      //其他选项
	Result       TopCreateTaskResult
}

func (task *TopCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *TopCreateTask) GetInhertType() string {
	return "top"
}

func (task *TopCreateTask) GetClientName() string {
	return "Top.Create"
}
