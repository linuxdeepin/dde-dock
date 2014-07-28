package bluetooth

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"bluetooth", Start, Stop, true})
}
