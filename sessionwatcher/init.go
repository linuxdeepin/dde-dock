package sessionwatcher

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "sessionwatcher",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
