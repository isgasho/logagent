package main

import (
	_ "hank.com/web-monitor/elastic"
	"hank.com/web-monitor/filelog"
)

const (
	//logpath = "/var/log/nginx/access.log"//linux测试路径
	logpath = "/data/logs/nginx/st.yunlaimi.com_access.log"

//logpath = "E:/phpstudy_pro/Extensions/Nginx1.15.11/logs/access.log"
//logpath = "data/access.log" //本地数据
)

func main() {
	filelog.StartGetLogServer(logpath)
}
