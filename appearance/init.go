package appearance

import (
	"pkg.linuxdeepin.com/dde-daemon/loader"
	"pkg.linuxdeepin.com/lib/log"
)

var logger = log.NewLogger("dde-daemon/appearance")

func init() {
	loader.Register(NewAppearanceDaemon(logger))
	//loader.Register(&loader.Module{
	//Name:   "appearance",
	//Start:  Start,
	//Stop:   Stop,
	//Enable: true,
	//})
}
