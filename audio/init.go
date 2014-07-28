package audio

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"audio", Start, Stop, true})
}
