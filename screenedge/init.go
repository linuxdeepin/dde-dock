package screenedge

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"screenedge", Start, Stop, true})
}
