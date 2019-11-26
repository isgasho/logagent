package logagentmain

import (
	"hank.com/logagent/input"
	"hank.com/logagent/server/kafka"
	"hank.com/logagent/server/storage"
)

func Main() {
	//客户端日志系统启动
	input.Run()

	//启动中间件-生产者
	kafka.ProducerRun()

	//服务端启动
	storage.Run()
}
