package top

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type TopExchangeTaskResult struct {
	app.Result
	Top *Top `json:"top,omitempty"`
}

type TopExchangeTask struct {
	app.Task
	FromId int64 `json:"fromId"`
	ToId   int64 `json:"toId"` // 0 时至顶
	Result TopExchangeTaskResult
}

func (task *TopExchangeTask) GetResult() interface{} {
	return &task.Result
}

func (task *TopExchangeTask) GetInhertType() string {
	return "top"
}

func (task *TopExchangeTask) GetClientName() string {
	return "Top.Exchange"
}
