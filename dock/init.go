package dock

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"dock", Start, Stop, true})
}
