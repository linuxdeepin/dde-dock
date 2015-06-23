package keybinding

import (
	"pkg.linuxdeepin.com/dde-daemon/loader"
	"pkg.linuxdeepin.com/lib/log"
)

var logger = log.NewLogger("dde-daemon/keybinding")

func init() {
	loader.Register(NewKeybindingDaemon(logger))
	//loader.Register(&loader.Module{
	//Name:   "keybinding",
	//Start:  Start,
	//Stop:   Stop,
	//Enable: true,
	//})
}
