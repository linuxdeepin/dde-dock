package screensaver

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "screensaver",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
