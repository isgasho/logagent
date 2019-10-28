package main

import "hank.com/web-monitor/filelog"

func main(){
	filelog.StartGetLogServer("data/access.log")
}
