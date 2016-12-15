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
	"encoding/json"
	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/pulse"
)

func isVolumeValid(v float64) bool {
	if v < 0 || v > pulse.VolumeUIMax {
		return false
	}
	return true
}

func playFeedback() {
	playFeedbackWithDevice("")
}

func playFeedbackWithDevice(device string) {
	soundutils.PlaySystemSound(soundutils.EventVolumeChanged, device, false)
}

func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}
