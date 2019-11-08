package kibanadiscover

import (
	"context"
	"time"

	"hank.com/goelastic/esc"

	"hank.com/web-monitor/foundation/db"

	"github.com/astaxie/beego"
	"hank.com/web-monitor/log"
)

type KibanaDiscover struct {
	Config    *Config
	IndexName string
	CommonLog *log.CommonLog
}

func NewKibanaDiscover(commonLog *log.CommonLog) *KibanaDiscover {
	indexName := beego.AppConfig.DefaultString("elastic.indexname", "weberr")
	return &KibanaDiscover{Config: nil, CommonLog: commonLog, IndexName: indexName}
}

func (kd *KibanaDiscover) GetIndexName() string {
	return kd.IndexName + "-" + time.Unix(kd.CommonLog.Timestamp, 0).Format("2006-01-02")
}

func (kd *KibanaDiscover) Start(ctx context.Context) {
	//indexName 获取索引
	indexName := kd.GetIndexName()

	//CreateTable 获取Table
	err := esc.GetElasticDefault().CreateTable(ctx, indexName)
	if err != nil {
		panic(err)
	}

	l := log.NewLoggerByCommonLog(kd.CommonLog)

	msg := l.WriteMsg(kd.CommonLog.LogLevel, kd.CommonLog.Message)

	//Build KibanaDiscover
	kibanaDiscover := &db.Monitor{Date: time.Now(), FieldsTag: kd.CommonLog.Module, Message: msg}

	esc.GetElasticDefault().Insert(ctx, indexName, "", kibanaDiscover)
}
