package screen_edges

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"screen_edges", Start, Stop, true})
}
