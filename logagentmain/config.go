package logagentmain

import (
	"flag"

	"hank.com/logagent/server/storage"

	"hank.com/logagent/log"
)

type config struct {
	lc *log.Config
	ec *storage.Config
	cf *flag.FlagSet
}
