package filelog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"hank.com/web-monitor/elastic"

	"github.com/hpcloud/tail"
)

const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

const (
	MARKSTR   = "log.gif?data=" //日志标志位置
	ENDSTR    = ` HTTP/1.1" `   //日志尾部标记
	INDEXNAME = "weberr"
)

//大小为8的数组
var levelPrefix = [LevelDebug + 1]string{"[M] ", "[A] ", "[C] ", "[E] ", "[W] ", "[N] ", "[I] ", "[D] "}

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
	for {
		//读取信息
		l := <-logMsgs

		//处理每行消息
		line := SplitLine(l.Msg)
		if line == "" {
			continue
		}

		fmt.Println(line)

		//KibanaDiscover
		kibanaDiscover := &elastic.KibanaDiscover{Date: time.Now().Unix(), FieldsTag: "js", Message: "你好"}

		errMonitor := &elastic.ErrMonitor{}
		err := json.Unmarshal([]byte(line), errMonitor)
		if err != nil {
			fmt.Println(err)
			return
		}

		//时间处理
		date := time.Unix(errMonitor.Timestamp, 0)
		errMonitor.Date = date.Format("2006-01-02 15:04:05")

		kibanaDiscover.FieldsTag = errMonitor.Module
		kibanaDiscover.Message = FormatErrMonitorMessage(errMonitor)

		ctx := context.Background()
		indexName := INDEXNAME + "-" + date.Format("2006-01-02")
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
			//Type("errmonitor").
			//Id("1").
			BodyJson(kibanaDiscover).
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

	//TODO 这里要改成正则匹配
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

//LogLevelMsg- 日志等级输出
func LogLevelMsg(level int) string {
	if level < 0 || level > LevelDebug {
		return ""
	}
	return levelPrefix[level]
}

func LogFileMsg(url string, line, col int64) string {
	_, file := filepath.Split(url)
	return fmt.Sprintf("[%v:%v:%v]", file, line, col)
}

//TODO func输出等级设置

//FormatErrMonitorMessage-转化为Kibana指定格式输出
func FormatErrMonitorMessage(errMonitor *elastic.ErrMonitor) (message string) {
	date := errMonitor.Date
	logLeve := LogLevelMsg(errMonitor.Errlevel)
	msgByte, _ := json.Marshal(errMonitor)
	fileMsg := LogFileMsg(errMonitor.ViewUrl, errMonitor.Line, errMonitor.Col)

	message = fmt.Sprintf("%v %v %v %v", date, logLeve, fileMsg, string(msgByte))
	return message
}
