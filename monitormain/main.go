package monitormain

import (
	"hank.com/web-monitor/kibanadiscover"
	"hank.com/web-monitor/log"
)

func Main() {
	//启动日志系统
	fl := log.NewFileLog()
	fl.StartServer()

	//kibanadiscover运行
	kibanadiscover.Run()
}
