package top

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TopRemoveTaskResult struct {
	app.Result
	Top *Top `json:"top,omitempty"`
}

type TopRemoveTask struct {
	app.Task
	Id     int64 `json:"id"`
	Result TopRemoveTaskResult
}

func (task *TopRemoveTask) GetResult() interface{} {
	return &task.Result
}

func (task *TopRemoveTask) GetInhertType() string {
	return "top"
}

func (task *TopRemoveTask) GetClientName() string {
	return "Top.Remove"
}
