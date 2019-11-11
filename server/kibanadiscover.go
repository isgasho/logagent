package server

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/astaxie/beego"

	"hank.com/goelastic/esc"
)

var ChanLog = make(chan string, 1)

type CommonLog struct {
	Module              string `json:"module"`              //出错的模块 应用的名称例如:xmiss
	ViewUrl             string `json:"viewurl"`             //请求的url
	LogLevel            int    `json:"loglevel"`            //错误等级 3err 4Warning 5Notice 7Debug
	FileName            string `json:"filename"`            //文件名称
	Line                int64  `json:"line"`                //文件所在的行
	Col                 int64  `json:"col"`                 //文件所在的列
	EnableFileDepthType int    `json:"enablefiledepthtype"` //是否需要格式化输出message 0不处理 1处理 2函数处理
	Message             string `json:"message"`             //自定义消息
	Platform            string `json:"platform"`            //系统架构
	Ua                  string `json:"ua"`                  //UserAgent浏览器信息
	Lang                string `json:"lang"`                //使用的语言
	Screen              string `json:"screen"`              //分辨率
	Carset              string `json:"carset"`              //浏览器编码环境
	Address             string `json:"address"`             //所在位置
	Date                string `json:"date"`                //发生的时间
	Timestamp           int64  `json:"timestamp"`           //发生的时间戳
}

type KibanaDiscover struct {
	Config    *Config
	IndexName string
	CommonLog *CommonLog
}

//NewKibanaDiscover -
func NewKibanaDiscover(indexName string) *KibanaDiscover {
	return &KibanaDiscover{Config: nil, IndexName: indexName}
}

func (kd *KibanaDiscover) GetIndexName() string {
	return kd.IndexName + "-" + time.Now().Format("2006-01-02")
}

func Run() {
	//启动KibanaDiscover TODO indexName写死
	indexName := beego.AppConfig.DefaultString("elastic.indexname", "weberr")
	kd := NewKibanaDiscover(indexName)
	//开始运行
	for {
		select {
		case bodyJson := <-ChanLog:
			//接收到了信息
			log.Println(bodyJson)
			kd.RunPush(bodyJson)
		}
	}
}

func (kd *KibanaDiscover) RunPush(bodyJson string) {

	ctx := context.Background()

	//indexName 获取索引
	indexName := kd.GetIndexName()

	//CreateTable 获取Table
	err := esc.GetElasticDefault().CreateTable(ctx, indexName)
	if err != nil {
		panic(err)
	}

	var bodyMap map[string]interface{}
	json.Unmarshal([]byte(bodyJson), &bodyMap)

	err = esc.GetElasticDefault().Insert(ctx, indexName, "", &bodyMap)
	if err != nil {
		log.Println("Elastic Insert err：" + bodyJson)
		//panic(err)
	}
}
