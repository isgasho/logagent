package log

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/astaxie/beego"

	"github.com/hpcloud/tail"
)

const (
	MARKSTR = "log.gif?log=" //日志标志位置
	ENDSTR  = ` HTTP/1.1" `  //日志尾部标记
)

var ChanLog = make(chan *CommonLog, 1)

type FileLog struct {
	Config *Config
}

var logMsg = make(chan string, 1)

func NewFileLog() *FileLog {
	config := &Config{
		Logpath: beego.AppConfig.String("filelog.logpath"),
	}
	return &FileLog{Config: config}
}

//StartGetLogServer- 启动日志服务 读写
func (fl *FileLog) StartServer() {
	go fl.ReadLogLoop()
	go fl.WriteLog2Ws() //TODO 异步
}

//ReadLogLoop- 读取日志  TODO channel
func (fl *FileLog) ReadLogLoop() {
	t, _ := tail.TailFile(fl.Config.Logpath, tail.Config{Follow: true})
	for line := range t.Lines {
		logMsg <- line.Text
	}
}

func (fl *FileLog) WriteLog2Ws() {
	for {
		//读取信息
		l := <-logMsg

		//处理每行消息
		commonLog := SplitLine(l)
		if commonLog == nil || commonLog.Message == "" {
			continue
		}

		fmt.Println(commonLog)
		ChanLog <- commonLog
	}
}

//SplitLine- 分解行
func SplitLine(msg string) *CommonLog {
	//TODO 这里要改成正则匹配
	comma := strings.Index(msg, MARKSTR)
	endComma := strings.Index(msg, ENDSTR)
	if comma < 0 || endComma < 0 {
		return nil
	}

	index := comma + len(MARKSTR)
	endindex := endComma

	if index > endindex || index > len(msg) || endindex > len(msg) {
		return nil
	}

	//头尾匹配去除得到数据
	line := msg[index:endindex]

	//url解码一下,nginx默认url编码
	if line != "" {
		line, _ = url.QueryUnescape(line)
	}

	fmt.Println(line)

	//解析参数
	commonLog := &CommonLog{}
	err := json.Unmarshal([]byte(line), commonLog)
	if err != nil {
		return nil
	}
	return commonLog
}
