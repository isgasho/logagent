package monitormain

import (
	"flag"

	"hank.com/web-monitor/kibanadiscover"
	"hank.com/web-monitor/log"
)

type config struct {
	lc *log.Config
	ec *kibanadiscover.Config
	cf *flag.FlagSet
}
