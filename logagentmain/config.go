package logagentmain

import (
	"flag"

	"hank.com/logagent/log"
	"hank.com/logagent/server"
)

type config struct {
	lc *log.Config
	ec *server.Config
	cf *flag.FlagSet
}
