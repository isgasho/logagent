package logagentmain

import (
	"flag"

	"hank.com/logagent/server/storage"

	"hank.com/logagent/input"
)

type config struct {
	lc *input.Config
	ec *storage.Config
	cf *flag.FlagSet
}
