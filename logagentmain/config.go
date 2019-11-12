package logagentmain

import (
	"flag"

	"hank.com/logagent/server/store"

	"hank.com/logagent/log"
)

type config struct {
	lc *log.Config
	ec *store.Config
	cf *flag.FlagSet
}
