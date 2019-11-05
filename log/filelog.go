package log

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"net/url"
	"strings"

	"github.com/hpcloud/tail"
)



const (
	MARKSTR   = "log.gif?log=" //日志标志位置
	ENDSTR    = ` HTTP/1.1" `  //日志尾部标记
)

var ChanLog = make(chan *CommonLog,1)

type CommonLog struct {
	Module    string `json:"module"`    //出错的模块 应用的名称例如:xmiss
	ViewUrl   string `json:"viewUrl"`   //请求的url
	LogLevel  int    `json:"loglevel"`  //错误等级 3err 4Warning 5Notice 7Debug
	File      string `json:"file"`      //出错的文件
	Line      int64  `json:"line"`      //出错文件所在行
	Col       int64  `json:"col"`       //出错文件所在列
	Message   string `json:"message"`   //自定义消息
	Platform  string `json:"platform"`  //系统架构
	Ua        string `json:"ua"`        //UserAgent浏览器信息
	Lang      string `json:"lang"`      //使用的语言
	Screen    string `json:"screen"`    //分辨率
	Carset    string `json:"carset"`    //浏览器编码环境
	Address   string `json:"address"`   //所在位置
	Date      string `json:"date"`      //发生的时间
	Timestamp int64  `json:"timestamp"` //发生的时间戳
}

type FileLog struct {
	Config *Config
}

var logMsg = make(chan string, 1)

func NewFileLog()*FileLog{
	config := &Config{
		Logpath:beego.AppConfig.String("filelog.logpath"),
	}
	return &FileLog{Config:config}
}

//StartGetLogServer- 启动日志服务 读写
func (fl *FileLog)StartServer() {
	go fl.ReadLogLoop()
	go fl.WriteLog2Ws() //TODO 异步
}

//ReadLogLoop- 读取日志  TODO channel
func (fl *FileLog)ReadLogLoop() {
	t, _ := tail.TailFile(fl.Config.Logpath, tail.Config{Follow: true})
	for line := range t.Lines {
		logMsg <-  line.Text
	}
}

func (fl *FileLog)WriteLog2Ws() {
	for {
		//读取信息
		l := <-logMsg

		//处理每行消息
		line := SplitLine(l)
		if line == "" {
			continue
		}

		fmt.Println(line)

		commonLog := &CommonLog{}
		err := json.Unmarshal([]byte(line), commonLog)
		if err != nil {
			panic(err)
		}

		ChanLog <-commonLog
	}
}

//SplitLine- 分解行
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

//LogFileMsg-文件格式生成
func LogFileMsg(file string, line, col int64) string {
	return fmt.Sprintf("[%v:%v:%v]", file, line, col)
}

//TODO func输出等级设置


