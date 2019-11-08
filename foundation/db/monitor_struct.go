package db

import (
	"time"

	"hank.com/web-monitor/log"
)

//通用监控结构体
type Monitor struct {
	log.CommonLog
	Date      time.Time `json:"@timestamp"` //发生的时间
	FieldsTag string    `json:"fields.tag"`
	Message   string    `json:"message"`
}

type ErrMessage struct {
	Bid   int64  `json:"bid"`  //集团Bid
	Soid  int64  `json:"soid"` //店铺Soid
	Sid   int64  `json:"sid"`
	Name  string `json:"name"`  //商户的ID
	Code  string `json:"code"`  //错误代码
	Info  string `json:"info"`  //错误信息
	Stack string `json:"stack"` //堆栈错误
}
