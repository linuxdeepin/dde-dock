package dsc

import (
	"pkg.linuxdeepin.com/dde-daemon/loader"
	"pkg.linuxdeepin.com/lib/log"
)

var logger = log.NewLogger("dde-daemon/dsc")

func init() {
	loader.Register(NewDSCDaemon(logger))
	//loader.Register(&loader.Module{
	//Name:   "dsc",
	//Start:  Start,
	//Stop:   Stop,
	//Enable: true,
	//})
}
