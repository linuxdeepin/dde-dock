package timedate

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "timedate",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
