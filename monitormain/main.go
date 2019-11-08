package monitormain

import (
	"context"

	"hank.com/web-monitor/kibanadiscover"
	"hank.com/web-monitor/log"
)

func Main() {
	//启动日志系统
	fl := log.NewFileLog()
	fl.StartServer()

	select {
	case commonLog := <-log.ChanLog:
		{
			//启动KibanaDiscover
			el := kibanadiscover.NewKibanaDiscover(commonLog)
			el.Start(context.Background())
		}
	}
}
