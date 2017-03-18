package top

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TopBatchCreateTaskResult struct {
	app.Result
	Top *Top `json:"top,omitempty"`
}

type TopBatchCreateTask struct {
	app.Task
	Alias        string      `json:"alias"`        //别名
	ItemIds      string      `json:"itemIds"`      //推荐项ID
	ItemType     string      `json:"itemType"`     //推荐项类型
	Tags         string      `json:"tags"`         //搜索标签
	CityId       int64       `json:"cityId"`       //城市ID
	CityPath     string      `json:"cityPath"`     //城市路径
	ClassifyId   int64       `json:"classifyId"`   //分类ID
	ClassifyPath string      `json:"classifyPath"` //分类路径
	GroupId      int64       `json:"groupId"`      //分组ID
	Options      interface{} `json:"options"`      //其他选项
	Result       TopBatchCreateTaskResult
}

func (task *TopBatchCreateTask) GetResult() interface{} {
	return &task.Result
}

func (task *TopBatchCreateTask) GetInhertType() string {
	return "top"
}

func (task *TopBatchCreateTask) GetClientName() string {
	return "Top.BatchCreate"
}
