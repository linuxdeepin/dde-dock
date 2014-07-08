package search

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"search", Start, Stop, true})
}
