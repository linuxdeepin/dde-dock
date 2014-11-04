package dsc

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "dsc",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
