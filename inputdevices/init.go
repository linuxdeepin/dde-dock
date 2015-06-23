package inputdevices

import (
	"pkg.linuxdeepin.com/dde-daemon/loader"
	"pkg.linuxdeepin.com/lib/log"
)

var logger = log.NewLogger("dde-daemon/inputdevices")

func init() {
	loader.Register(NewInputdevicesDaemon(logger))
	//loader.Register(&loader.Module{
	//Name:   "inputdevices",
	//Start:  Start,
	//Stop:   Stop,
	//Enable: true,
	//})
}
