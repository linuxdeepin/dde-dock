package themes

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "themes",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
