package appearance

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "appearance",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
