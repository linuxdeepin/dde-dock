package bluetooth

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/bluetooth")

func init() {
	loader.Register(newBluetoothDaemon(logger))
	//loader.Register(&loader.Module{
	//Name:   "bluetooth",
	//Start:  Start,
	//Stop:   Stop,
	//Enable: true,
	//})
}
