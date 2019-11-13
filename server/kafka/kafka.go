package kafka

import (
	"log"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
)

var BodyJson = make(chan string, 0)

//NewKafkaProducer-
func NewKafkaProducer() (client sarama.SyncProducer, err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          //发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //新选出一个partition
	config.Producer.Return.Successes = true                   //// 成功交付的消息将在success channel返回
	addrs := beego.AppConfig.String("kafka.addrs")
	return sarama.NewSyncProducer([]string{addrs}, config)
}

func ProducerRun() {
	go WriteLoop()
}

func WriteLoop(){
	//TODO 出现错误了是否可以尝试重连
	product, err := NewKafkaProducer()
	if err != nil {
		log.Println("NewKafkaProducer Err：" + err.Error())
		panic(err)
	}

	for{

		bodyJson := <- BodyJson

		//发送数据
		msg := &sarama.ProducerMessage{}
		msg.Topic = beego.AppConfig.DefaultString("elastic.indexname", "weberr")
		msg.Value = sarama.StringEncoder(bodyJson)
		pid, offset, err := product.SendMessage(msg)
		if err != nil {
			//这里报错就打印一下先把
			log.Println("NewKafkaProducer Err：" + err.Error())
			continue
		}
		log.Printf("pid:%v offset:%v\n", pid, offset)
	}
}
