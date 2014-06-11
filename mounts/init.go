package mounts

import "dde-daemon"

func init() {
	loader.Register(&loader.Module{"mounts", Start, nil, true})
}
