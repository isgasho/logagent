package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"
	"github.com/astaxie/beego"
	"hank.com/logagent/server/esc"
	"github.com/olivere/elastic/config"
)

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

type ElasticMessage struct {
	IndexName string
	Value     []byte
}

//NewKibanaDiscover -
func NewElasticMessage(indexName string) *ElasticMessage {
	return &ElasticMessage{IndexName: indexName}
}

func GetIndexName(indexName string) string {
	return indexName + "-" + time.Now().Format("2006-01-02")
}

func Run() {
	var addrs []string
	if addrStr := beego.AppConfig.String("kafka.addrs");addrStr != ""{
		json.Unmarshal([]byte(addrStr),&addrs)
	}
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		panic(err)
		return
	}

	//TODO indexName写死
	topic := beego.AppConfig.DefaultString("elastic.indexname", "weberr")

	//开始运行
	for {
		partitionList, err := consumer.Partitions(topic) // 根据topic取到所有的分区
		if err != nil {
			fmt.Printf("fail to get list of partition:err%v\n", err)
			panic(err)
			return
		}

		for partition := range partitionList { // 遍历所有的分区
			// 针对每个分区创建一个对应的分区消费者
			pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
			if err != nil {
				fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
				return
			}
			defer pc.AsyncClose()
			// 异步从每个分区消费信息
			for msg := range pc.Messages() {
				elasticMessage := &ElasticMessage{IndexName:topic,Value:msg.Value}
				StoringData(elasticMessage)
			}
		}
	}
}

//StoringData -
func StoringData(elasticMessage *ElasticMessage) {

	ctx := context.Background()

	//indexName 获取索引
	indexName := GetIndexName(elasticMessage.IndexName)

	sniff := false
	elasticSource,err := esc.NewClient(&config.Config{URL:beego.AppConfig.String("elastic.url"),Sniff:&sniff})
	if err != nil{
		log.Println("Elastic NewClient err：" + beego.AppConfig.String("elastic.url"))
		panic(err)
	}

	//CreateTable 获取Table
	err = elasticSource.CreateTable(ctx, indexName)
	if err != nil {
		panic(err)
	}

	var bodyMap map[string]interface{}
	err =json.Unmarshal(elasticMessage.Value, &bodyMap)
	if err != nil{
		log.Println("Elastic bodyMap Unmarshal err：" + string(elasticMessage.Value))
		return
	}

	err = elasticSource.Insert(ctx, indexName, "", &bodyMap)
	if err != nil {
		log.Println("Elastic Insert err：" + string(elasticMessage.Value))
		//panic(err) 会报解析错误
	}
}
