package keybinding

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/keybinding")

func init() {
	loader.Register(NewKeybindingDaemon(logger))
	//loader.Register(&loader.Module{
	//Name:   "keybinding",
	//Start:  Start,
	//Stop:   Stop,
	//Enable: true,
	//})
}
