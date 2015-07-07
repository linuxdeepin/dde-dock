package audio

import (
	"pkg.deepin.io/dde-daemon/loader"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("dde-daemon/audio")

func init() {
	loader.Register(NewAudioDaemon(logger))
	//loader.Register(&loader.Module{
	//Name:   "audio",
	//Start:  Start,
	//Stop:   Stop,
	//Enable: true,
	//})
}
