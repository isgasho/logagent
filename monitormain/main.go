package monitormain

import (
	"context"
	"hank.com/web-monitor/elastic"
	"hank.com/web-monitor/log"
)

func Main(){
	//启动日志系统
	fl := log.NewFileLog()
	fl.StartServer()

	select{
		case commonLog := <- log.ChanLog:{
			//启动ElasticSearch
			el :=elastic.NewElastic(commonLog)
			el.Start(context.Background())
		}
	}
}
