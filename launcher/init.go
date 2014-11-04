package launcher

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "launcher",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
