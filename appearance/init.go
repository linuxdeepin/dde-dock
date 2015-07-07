package appearance

import (
	"pkg.deepin.io/dde-daemon/loader"
	"pkg.deepin.io/lib/log"
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
