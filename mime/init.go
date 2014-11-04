package mime

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "mime",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
