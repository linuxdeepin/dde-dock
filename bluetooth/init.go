package bluetooth

import (
	"pkg.linuxdeepin.com/dde-daemon/loader"
	"pkg.linuxdeepin.com/lib/log"
)

var logger = log.NewLogger("dde-daemon/bluetooth")

func init() {
	loader.Register(NewBluetoothDaemon(logger))
	//loader.Register(&loader.Module{
	//Name:   "bluetooth",
	//Start:  Start,
	//Stop:   Stop,
	//Enable: true,
	//})
}
