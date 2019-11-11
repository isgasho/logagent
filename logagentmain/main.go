package logagentmain

import (
	"hank.com/logagent/log"
	"hank.com/logagent/server"
)

func Main() {
	//客户端日志系统启动
	fl := log.NewFileLog()
	fl.StartServer()

	//服务端启动
	server.Run()

}
