package elastic

import (
	"context"
	"fmt"
	"time"

	"hank.com/web-monitor/filelog"
)

//KibanaDiscover- 对应Kibana的模板
type KibanaDiscover struct {
	Date      time.Time `json:"@timestamp"` //发生的时间
	FieldsTag string    `json:"fields.tag"`
	Message   string    `json:"message"`
}

// ErrMonitor is a structure used for serializing/deserializing data in Elasticsearch.
type ErrMonitor struct {
	Bid       int64  `json:"bid"`  //集团Bid
	Soid      int64  `json:"soid"` //店铺Soid
	Sid       int64  `json:"sid"`
	Name      string `json:"name"`      //商户的ID
	Module    string `json:"module"`    //出错的模块
	ViewUrl   string `json:"viewUrl"`   //url
	Address   string `json:"address"`   //所在位置
	Platform  string `json:"platform"`  //系统架构
	Ua        string `json:"ua"`        //UserAgent浏览器信息
	File      string `json:"file"`      //出错的文件
	Line      int64  `json:"line"`      //出错文件所在行
	Col       int64  `json:"col"`       //出错文件所在列
	Lang      string `json:"lang"`      //使用的语言
	Screen    string `json:"screen"`    //分辨率
	Carset    string `json:"carset"`    //浏览器编码环境
	Errlevel  int    `json:"errlevel"`  //错误等级 3err 4Warning 5Notice 7Debug
	Code      string `json:"code"`      //错误代码
	Info      string `json:"info"`      //错误信息
	Stack     string `json:"stack"`     //堆栈错误
	Date      string `json:"date"`      //发生的时间
	Timestamp int64  `json:"timestamp"` //发生的时间戳
}

func CreateTable(ctx context.Context, tableName string) error {
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

//NewKibanaDiscoverByErrMonitor- 根据errMonitor转化KibanaDiscover
func NewKibanaDiscoverByErrMonitor(errMonitor *ErrMonitor) *KibanaDiscover {
	kibanaDiscover := &KibanaDiscover{Date: time.Now()}
	kibanaDiscover.FieldsTag = errMonitor.Module
	kibanaDiscover.Message = filelog.FormatErrMonitorMessage(errMonitor)
	return kibanaDiscover
}

//BuildKibanaDiscover- 创建索引并生成数据
func BuildKibanaDiscover(ctx context.Context, indexName string, kibanaDiscover *KibanaDiscover) {

	//创建elastic索引
	err := CreateTable(ctx, indexName)
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
