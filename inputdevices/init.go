package inputdevices

import (
	"pkg.deepin.io/dde-daemon/loader"
	"pkg.deepin.io/lib/log"
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
