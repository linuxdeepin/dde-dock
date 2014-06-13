package dock

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"dock", Start, Stop, true})
}
