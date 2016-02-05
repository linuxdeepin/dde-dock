/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package screenedge

import (
	"time"
)

type edgeTimer struct {
	timer *time.Timer
}

func (eTimer *edgeTimer) Start(timeout int32, action func()) {
	eTimer.timer = time.NewTimer(time.Millisecond * time.Duration(timeout))
	go func() {
		if eTimer.timer == nil {
			return
		}
		<-eTimer.timer.C
		action()
		eTimer.timer = nil
	}()
}

func (eTimer *edgeTimer) Stop() {
	if eTimer.timer == nil {
		return
	}

	eTimer.timer.Stop()
	eTimer.timer = nil
}
