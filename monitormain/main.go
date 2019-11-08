package monitormain

import (
	"context"

	"github.com/astaxie/beego"

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
			//启动KibanaDiscover TODO indexName写死
			indexName := beego.AppConfig.DefaultString("elastic.indexname", "weberr")

			kd := kibanadiscover.NewKibanaDiscover(indexName, commonLog)
			kd.Start(context.Background())
		}
	}
}
