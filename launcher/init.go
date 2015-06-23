package launcher

import (
	"pkg.linuxdeepin.com/dde-daemon/loader"
	"pkg.linuxdeepin.com/lib/log"
)

var logger = log.NewLogger("dde-daemon/launcher-daemon")

func init() {
	loader.Register(NewLauncherDaemon(logger))
	// loader.Register(&loader.Module{
	// 	Name:   "launcher",
	// 	Start:  Start,
	// 	Stop:   Stop,
	// 	Enable: true,
	// })
}
