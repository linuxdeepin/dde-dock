/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package audio

import (
	"pkg.deepin.io/lib/pulse"
)

type Port struct {
	Name        string
	Description string
	Available   byte // Unknow:0, No:1, Yes:2
}

func toPort(v pulse.PortInfo) Port {
	return Port{
		Name:        v.Name,
		Description: v.Description,
		Available:   byte(v.Available),
	}
}

func portsEqual(a, b []Port) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
