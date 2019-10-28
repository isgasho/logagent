package filelog

import (
	"encoding/json"
	"fmt"
	"strings"

	"hank.com/web-monitor/elastic"

	"github.com/hpcloud/tail"
)

const (
	markStr = "log.gif?" //日志标志位置
	endStr  = ` HTTP/1.1" `
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

		//计数器
		i++

		errMonitor := &elastic.ErrMonitor{}
		err := json.Unmarshal([]byte(line), errMonitor)
		if err != nil {
			fmt.Println(err)
			return
		}

		//发送数据

	}
}

//SplitLine-
func SplitLine(msg string) (line string) {
	comma := strings.Index(msg, markStr)
	if comma == -1 {
		return
	}

	//这里要改成正则匹配
	endComma := strings.Index(msg, endStr)

	index := comma + len(markStr)
	endindex := endComma

	//头尾匹配去除得到数据
	line = msg[index:endindex]

	return
}
