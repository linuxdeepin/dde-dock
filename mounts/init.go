package mounts

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "mounts",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
