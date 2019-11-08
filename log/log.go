package log

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"time"
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
	FileDepthNothing = iota
	FileDepthCommonLog
	FileDepthFuncCall
)

//大小为8的数组
var levelPrefix = [LevelDebug + 1]string{"[M] ", "[A] ", "[C] ", "[E] ", "[W] ", "[N] ", "[I] ", "[D] "}

type CommonLog struct {
	Module              string `json:"module"`              //出错的模块 应用的名称例如:xmiss
	ViewUrl             string `json:"viewurl"`             //请求的url
	LogLevel            int    `json:"loglevel"`            //错误等级 3err 4Warning 5Notice 7Debug
	FileName            string `json:"filename"`            //出错的文件
	Line                int64  `json:"line"`                //出错文件所在行
	Col                 int64  `json:"col"`                 //出错文件所在列
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

type Log interface {
	SetLevel(level int)
	WriteMsg(logLevel int, msg string, v ...interface{}) string
	Error(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Notice(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
}

type Logger struct {
	*CommonLog
	Level               int
	loggerFuncCallDepth int
	Switchlog           bool //日志开关
}

//NewLogger-
func NewLogger() *Logger {
	l := new(Logger)
	l.Level = LevelDebug
	l.loggerFuncCallDepth = 2
	l.EnableFileDepthType = FileDepthFuncCall
	l.Switchlog = true
	return l
}

//NewLoggerByCommonLog-
func NewLoggerByCommonLog(commonLog *CommonLog) Log {
	l := new(Logger)
	l.CommonLog = commonLog
	l.Level = LevelDebug
	l.loggerFuncCallDepth = 2
	l.Switchlog = true
	return l
}

//SetLevel-
func (l *Logger) SetLevel(level int) {
	l.Level = level
}

//WriteMsg-
func (l *Logger) WriteMsg(logLevel int, msg string, v ...interface{}) string {
	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v...)
	}

	if l.EnableFileDepthType == FileDepthFuncCall {
		_, file, line, ok := runtime.Caller(l.loggerFuncCallDepth)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		msg = "[" + filename + ":" + strconv.Itoa(line) + "] " + msg
	}

	if l.EnableFileDepthType == FileDepthCommonLog {
		var filename string
		var line, col int64
		if l.FileName == "" {
			filename = "???"
		}
		filename = l.CommonLog.FileName
		line = l.CommonLog.Line
		col = l.CommonLog.Col

		msg = "[" + filename + ":" + strconv.FormatInt(line, 10) + ":" + strconv.FormatInt(col, 10) + "] " + msg
	}

	//level
	var levelPre = ""
	if logLevel > 0 && logLevel < LevelDebug {
		levelPre = levelPrefix[logLevel]
	}
	msg = levelPre + msg

	msg = time.Unix(l.CommonLog.Timestamp, 0).Format("2006/01/02 15:04:05") + " " + msg

	return msg
}

//Error-
func (l *Logger) Error(format string, v ...interface{}) {
	if !l.Switchlog {
		return
	}
	if LevelError > l.Level {
		return
	}
	l.WriteMsg(LevelError, format, v...)
}

//Warn -
func (l *Logger) Warn(format string, v ...interface{}) {
	if !l.Switchlog {
		return
	}
	if LevelWarning > l.Level {
		return
	}
	l.WriteMsg(LevelWarning, format, v...)
}

//Notice-
func (l *Logger) Notice(format string, v ...interface{}) {
	if !l.Switchlog {
		return
	}
	if LevelNotice > l.Level {
		return
	}
	l.WriteMsg(LevelNotice, format, v...)
}

//Info-
func (l *Logger) Info(format string, v ...interface{}) {
	if !l.Switchlog {
		return
	}
	if LevelInformational > l.Level {
		return
	}
	l.WriteMsg(LevelInformational, format, v...)
}

//Debug-
func (l *Logger) Debug(format string, v ...interface{}) {
	if !l.Switchlog {
		return
	}
	if LevelDebug > l.Level {
		return
	}
	l.WriteMsg(LevelDebug, format, v...)
}
