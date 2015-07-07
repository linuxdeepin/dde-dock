package dsc

import (
	"pkg.deepin.io/dde-daemon/loader"
	"pkg.deepin.io/lib/log"
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
