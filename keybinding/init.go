package keybinding

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"keybinding", Start, Stop, true})
}
