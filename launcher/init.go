package launcher

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"launcher", Start, Stop, true})
}
