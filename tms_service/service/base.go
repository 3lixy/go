package service

import (
	"azoya/nova"
)

//BaseService 模型
type BaseService struct {
	C *nova.Context
}

type Result struct {
	Status  uint64      `json:"status"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}
