package elastic

import (
	"context"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"hank.com/web-monitor/log"
)

//KibanaDiscover- 对应Kibana的模板
type KibanaDiscover struct {
	Date      time.Time `json:"@timestamp"` //发生的时间
	FieldsTag string    `json:"fields.tag"`
	Message   string    `json:"message"`
}

// ErrMessage is a structure used for serializing/deserializing data in Elasticsearch.
type ErrMessage struct {
	Bid   int64  `json:"bid"`  //集团Bid
	Soid  int64  `json:"soid"` //店铺Soid
	Sid   int64  `json:"sid"`
	Name  string `json:"name"`  //商户的ID
	Code  string `json:"code"`  //错误代码
	Info  string `json:"info"`  //错误信息
	Stack string `json:"stack"` //堆栈错误
}

type Elastic struct {
	Config    *Config
	CommonLog *log.CommonLog
}

func NewElastic(commonLog *log.CommonLog) *Elastic {
	config := &Config{
		IndexName: beego.AppConfig.DefaultString("elastic.indexname", "weberr"),
	}
	return &Elastic{Config: config, CommonLog: commonLog}
}

func (ec *Elastic) GetIndexName() string {
	return ec.Config.IndexName + "-" + time.Unix(ec.CommonLog.Timestamp, 0).Format("2006-01-02")
}

func (ec *Elastic) CreateTable(ctx context.Context, tableName string) error {
	exists, err := ElasticClient().IndexExists(tableName).Do(ctx)
	if err != nil {
		return err
	}

	if !exists {
		createIndex, err := ElasticClient().CreateIndex(tableName).Do(ctx)
		if err != nil {
			return err
		}
		if createIndex.Acknowledged {
			//Not acknowledged
		}
	}
	return nil
}

//BuildKibanaDiscover- 创建索引并生成数据
func (ec *Elastic) BuildKibanaDiscover(ctx context.Context, indexName string, kibanaDiscover *KibanaDiscover) {

	//创建elastic索引
	err := ec.CreateTable(ctx, indexName)
	if err != nil {
		panic(err)
	}

	//设置数据
	put1, err := ElasticClient().Index().
		Index(indexName).
		BodyJson(kibanaDiscover).
		Do(ctx)

	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Indexed errmonitor %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
}

func (ec *Elastic) Start(ctx context.Context) {
	//indexName 获取索引
	indexName := ec.GetIndexName()

	//CreateTable 获取Table
	err := ec.CreateTable(ctx, indexName)
	if err != nil {
		panic(err)
	}

	l := log.NewLoggerByCommonLog(ec.CommonLog)

	msg := l.WriteMsg(ec.CommonLog.LogLevel, ec.CommonLog.Message)

	//Build KibanaDiscover
	kibanaDiscover := &KibanaDiscover{Date: time.Now(), FieldsTag: ec.CommonLog.Module, Message: msg}
	ec.BuildKibanaDiscover(ctx, indexName, kibanaDiscover)
}
