package bluetooth

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"bluetooth", Start, Stop, true})
}
