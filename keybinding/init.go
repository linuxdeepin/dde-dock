package keybinding

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "keybinding",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
