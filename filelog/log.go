package filelog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"hank.com/web-monitor/elastic"

	"github.com/hpcloud/tail"
)

var (
	MARKSTR   = "log.gif?data=" //日志标志位置
	ENDSTR    = ` HTTP/1.1" `   //日志尾部标记
	INDEXNAME = "weberr"
)

var logMsgs = make(chan *WsMsg, 1000)

type WsMsg struct {
	Msg string `json:msg`
}

func StartGetLogServer(logpath string) {
	go ReadLogLoop(logpath)
	WriteLog2Ws() //TODO 异步
}

//TODO channel
func ReadLogLoop(logpath string) {
	t, _ := tail.TailFile(logpath, tail.Config{Follow: true})
	for line := range t.Lines {
		logMsgs <- &WsMsg{Msg: line.Text}
	}
}

func WriteLog2Ws() {
	var i int
	for {
		//读取信息
		l := <-logMsgs

		//处理每行消息
		line := SplitLine(l.Msg)
		if line == "" {
			continue
		}

		fmt.Println(line)

		//计数器
		i++

		errMonitor := &elastic.ErrMonitor{}
		err := json.Unmarshal([]byte(line), errMonitor)
		if err != nil {
			fmt.Println(err)
			return
		}

		lineTime := time.Unix(errMonitor.Timestamp, 0)
		indexName := INDEXNAME + "-" + lineTime.Format("2006-01-02")

		ctx := context.Background()
		//创建elastic索引
		exists, err := elastic.ElasticClient().IndexExists(indexName).Do(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		if !exists {
			createIndex, err := elastic.ElasticClient().CreateIndex(indexName).Do(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}
			if createIndex.Acknowledged {
				//Not acknowledged
			}
		}

		//设置数据
		put1, err := elastic.ElasticClient().Index().
			Index(indexName).
			Type("errmonitor").
			//Id("1").
			BodyJson(errMonitor).
			Do(ctx)

		if err != nil {
			// Handle error
			panic(err)
		}

		fmt.Printf("Indexed errmonitor %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	}
}

//SplitLine-
func SplitLine(msg string) (line string) {
	comma := strings.Index(msg, MARKSTR)
	if comma == -1 {
		return
	}

	//这里要改成正则匹配
	endComma := strings.Index(msg, ENDSTR)

	index := comma + len(MARKSTR)
	endindex := endComma

	//头尾匹配去除得到数据
	line = msg[index:endindex]

	//url解码一下,nginx默认url编码
	if line != "" {
		line, _ = url.QueryUnescape(line)
	}
	return
}
