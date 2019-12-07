package input

import (
	"fmt"
	"hank.com/logagent/server/kafka"
	"hank.com/logagent/util"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	// "hank.com/logagent/server/kafka"

	"github.com/astaxie/beego"
	"github.com/hpcloud/tail"
)

const (
	MARKSTR = "log.gif?log=" //日志标志位置
	ENDSTR  = ` HTTP/1.1" `  //日志尾部标记
)

var(
	offsetPath = fmt.Sprintf("%v%v",beego.AppConfig.DefaultString("filelog.dir","./data/"),
		"offset")
)

type FileLog struct {
	Config *Config
}

//StartGetLogServer- 启动日志服务 读写
func Run() {
	//初始化
	go ReadLogLoop()
}

//ReadLogLoop- 读取日志  TODO channel
func ReadLogLoop() {
	//文件名称
	logpath := beego.AppConfig.String("filelog.logpath")
	_,fileName := filepath.Split(logpath)

	var offset int64
	config := tail.Config{Follow:true,Location:&tail.SeekInfo{}}
	//获取文件记录的偏移量
	config.Location.Offset,_ = GetFileOffset(&offset)
	t, _ := tail.TailFile(logpath,config)

	for line := range t.Lines {
		//读取并设置最新文件偏移量信息
		offset,_ = t.Tell()

		//处理每行消息
		bodyJson := SplitLine(line.Text)
		if bodyJson == ""{
			SetFileOffset(fileName,offset)
			continue
		}

		kafka.BodyJson <- bodyJson

		//TODO 记住中断已经处理的读取位置
		SetFileOffset(fileName,offset)
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

//SetFileOffset -设置文件的偏移量
func SetFileOffset(name string,offset int64) (err error){
	content := fmt.Sprintf("%v %v",offset,name)
	err = util.WriteFile(offsetPath,content)
	return
}

//GetFileOffset -获取文件的偏移量
func GetFileOffset(offsetAdd *int64) (offset int64,err error){
	//从内存读取
	if *offsetAdd !=0{
		return  *offsetAdd,nil
	}

	//持久化硬盘读取
	oStr,err  := util.ReadFile(offsetPath)
	if err != nil{
		//一开始木有文件会抛出异常
		log.Println("ReadFile Err：" + err.Error())
		return 0,err
	}
	ot := strings.Split(oStr," ")
	offset, _ = strconv.ParseInt(ot[0], 10, 64)

	return offset,nil
}
