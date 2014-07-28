package datetime

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"datetime", Start, Stop, true})
}
