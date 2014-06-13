package datetime

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"datetime", Start, Stop, true})
}
