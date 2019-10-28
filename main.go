package main

import (
	_ "hank.com/web-monitor/elastic"
	"hank.com/web-monitor/filelog"
)

const (
	//logpath = "/var/log/nginx/access.log"
	logpath = "data/access.log"
)

func main() {
	filelog.StartGetLogServer(logpath)
}
