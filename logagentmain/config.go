package logagentmain

import (
	"flag"

	"github.com/friendlyhank/logagent/server/storage"

	"github.com/friendlyhank/logagent/input"
)

type config struct {
	lc *input.Config
	ec *storage.Config
	cf *flag.FlagSet
}
