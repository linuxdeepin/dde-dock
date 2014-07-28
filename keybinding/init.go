package keybinding

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"keybinding", Start, Stop, true})
}
