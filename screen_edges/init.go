package screen_edges

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{"screen_edges", Start, Stop, true})
}
