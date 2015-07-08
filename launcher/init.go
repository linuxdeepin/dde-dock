package launcher

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/launcher-daemon")

func init() {
	loader.Register(NewLauncherDaemon(logger))
	// loader.Register(&loader.Module{
	// 	Name:   "launcher",
	// 	Start:  Start,
	// 	Stop:   Stop,
	// 	Enable: true,
	// })
}
