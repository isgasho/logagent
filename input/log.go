package input

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/friendlyhank/logagent/server/storage"
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
	*storage.CommonLog
	Level               int
	loggerFuncCallDepth int
	EnableFileDepthType int
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
func NewLoggerByCommonLog(commonLog *storage.CommonLog) Log {
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
