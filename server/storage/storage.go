package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"github.com/olivere/elastic/config"
	"hank.com/logagent/server/esc"
)

type CommonLog struct {
	Module    string `json:"module"`    //出错的模块 应用的名称例如:xmiss
	ViewUrl   string `json:"viewurl"`   //请求的url
	LogLevel  int    `json:"loglevel"`  //错误等级 3err 4Warning 5Notice 7Debug
	FileName  string `json:"filename"`  //文件名称
	Line      int64  `json:"line"`      //文件所在的行
	Col       int64  `json:"col"`       //文件所在的列
	Platform  string `json:"platform"`  //系统架构
	Ua        string `json:"ua"`        //UserAgent浏览器信息
	Lang      string `json:"lang"`      //使用的语言
	Screen    string `json:"screen"`    //分辨率
	Carset    string `json:"carset"`    //浏览器编码环境
	Address   string `json:"address"`   //所在位置
	Date      string `json:"date"`      //发生的时间
	Timestamp int64  `json:"timestamp"` //发生的时间戳
}

type ElasticMessage struct {
	IndexName string
	CommonLog *CommonLog
	Value     []byte

	UpTime time.Time
}

//NewElasticMessage-
func NewElasticMessage(indexName string,value []byte)(elasticMessage *ElasticMessage,err error){
	e := &ElasticMessage{IndexName:indexName,Value:value}

	//解析公共参数
	e.CommonLog = &CommonLog{}
	err = json.Unmarshal(value, e.CommonLog)
	if err != nil {
		log.Printf("commonLog Unmarshal： %v|Err|%v",string(value),err)
		e.CommonLog = nil
	}

	//记录上传时间
	e.UpTime = time.Unix(e.CommonLog.Timestamp, 0)
	return e,err
}

func (e *ElasticMessage) GetIndexNameWithTime(indexName string) string {
	return indexName + "-" + e.UpTime.Format("2006-01-02")
}

func Run() {
	var addrs []string
	if addrStr := beego.AppConfig.String("kafka.addrs"); addrStr != "" {
		json.Unmarshal([]byte(addrStr), &addrs)
	}

	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		log.Printf("fail to start consumer, err:%v\n", err)
		panic(err)
		return
	}

	//TODO indexName写死
	topic := beego.AppConfig.DefaultString("elastic.indexname", "weberr")
	partition,_ := beego.AppConfig.Int("kafka.partition")

	partitionConsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
	if err != nil {
		fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
		return
	}
	defer partitionConsumer.AsyncClose()
	for{
		msg := <- partitionConsumer.Messages()

		log.Println("msg reciver："+string(msg.Value))
		elasticMessage,_ := NewElasticMessage(topic,msg.Value)
		StoringData(elasticMessage)
	}
}

//StoringData -
func StoringData(elasticMessage *ElasticMessage) {

	var bodyMap map[string]interface{}
	err := json.Unmarshal(elasticMessage.Value, &bodyMap)
	if err != nil {
		log.Printf("Elastic bodyMap Unmarshal： %v|Err|%v",string(elasticMessage.Value),err)
		return
	}

	ctx := context.Background()

	//indexName 获取索引
	indexName := elasticMessage.GetIndexNameWithTime(elasticMessage.IndexName)

	sniff := false
	elasticSource, err := esc.NewClient(&config.Config{URL: beego.AppConfig.String("elastic.url"), Sniff: &sniff})
	if err != nil {
		log.Println("Elastic NewClient err：" + beego.AppConfig.String("elastic.url"))
		panic(err)
	}

	//CreateTable 获取Table
	err = elasticSource.CreateTable(ctx, indexName)
	if err != nil {
		panic(err)
	}

	bodyMap["time"]=elasticMessage.UpTime

	err = elasticSource.Insert(ctx, indexName, "", &bodyMap)
	if err != nil {
		log.Println("Elastic Insert err：" + string(elasticMessage.Value))
		//panic(err) 会报解析错误
	}
}
