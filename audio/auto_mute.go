/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package audio

import (
	"pkg.deepin.io/lib/pulse"
	"strings"
)

func (a *Audio) autoMuteDetect() {
	// TODO: bluetooth, headset, dev.name
	if a.DefaultSink != nil {
		var state = pulse.AvailableTypeUnknow
		if isHeadphonePort(a.prevActivePort) {
			for _, port := range a.DefaultSink.Ports {
				if port.Name == a.prevActivePort {
					state = int(port.Available)
					break
				}
			}
		}
		a.prevActivePort = a.DefaultSink.ActivePort.Name
		if state == pulse.AvailableTypeNo {
			a.DefaultSink.SetMute(true)
		}
	}
}

func isHeadphonePort(portName string) bool {
	return strings.Contains(strings.ToLower(portName), "headphone")
}
