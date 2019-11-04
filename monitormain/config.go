package monitormain

import (
	"flag"
	"hank.com/web-monitor/elastic"
	"hank.com/web-monitor/log"
)

type config struct{
	lc *log.Config
	ec *elastic.Config
	cf *flag.FlagSet
}

