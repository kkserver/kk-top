package top

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TopTaskResult struct {
	app.Result
	Top *Top `json:"top,omitempty"`
}

type TopTask struct {
	app.Task
	Id     int64 `json:"id"`
	Result TopTaskResult
}

func (task *TopTask) GetResult() interface{} {
	return &task.Result
}

func (task *TopTask) GetInhertType() string {
	return "top"
}

func (task *TopTask) GetClientName() string {
	return "Top.Get"
}
