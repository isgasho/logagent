package log

import (
	"log"
	"net/url"
	"strings"

	"hank.com/logagent/server/kafka"

	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
	"github.com/hpcloud/tail"
)

const (
	MARKSTR = "log.gif?log=" //日志标志位置
	ENDSTR  = ` HTTP/1.1" `  //日志尾部标记
)

type FileLog struct {
	Config *Config
}

var bodyJson = make(chan string, 0)

//StartGetLogServer- 启动日志服务 读写
func Run() {
	//初始化
	go ReadLogLoop()
}

//ReadLogLoop- 读取日志  TODO channel
func ReadLogLoop() {
	t, _ := tail.TailFile(beego.AppConfig.String("filelog.logpath"), tail.Config{Follow: true})
	for line := range t.Lines {
		//处理每行消息
		bodyJson := SplitLine(line.Text)

		//TODO 封装
		product, err := kafka.NewKafkaProducer()
		if err != nil {
			log.Println("NewKafkaProducer Err：" + err.Error())
			panic(err)
		}

		//发送数据
		msg := &sarama.ProducerMessage{}
		msg.Topic = beego.AppConfig.DefaultString("elastic.indexname", "weberr")
		msg.Value = sarama.StringEncoder(bodyJson)
		pid, offset, err := product.SendMessage(msg)
		if err != nil {
			log.Println("NewKafkaProducer Err：" + err.Error())
			panic(err)
		}
		log.Printf("pid:%v offset:%v\n", pid, offset)
	}
}

//SplitLine- 分解行
func SplitLine(msg string) string {
	//TODO 这里要改成正则匹配
	comma := strings.Index(msg, MARKSTR)
	endComma := strings.Index(msg, ENDSTR)
	if comma < 0 || endComma < 0 {
		return ""
	}

	index := comma + len(MARKSTR)
	endindex := endComma

	if index > endindex || index > len(msg) || endindex > len(msg) {
		//打印下日志 TODO 优化
		log.Println("Index Err：" + msg)
		return ""
	}

	//头尾匹配去除得到数据
	line := msg[index:endindex]

	//url解码一下,nginx默认url编码
	if line != "" {
		line, _ = url.QueryUnescape(line)
	}

	return line
}
