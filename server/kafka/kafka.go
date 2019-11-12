package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
)

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

}
