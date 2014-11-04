package mpris

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "mpris",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
