package audio

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "audio",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
