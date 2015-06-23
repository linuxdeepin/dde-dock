package audio

import (
	"pkg.linuxdeepin.com/dde-daemon/loader"
	"pkg.linuxdeepin.com/lib/log"
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
