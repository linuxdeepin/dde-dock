/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	"time"
)

type secondTicker struct {
	action func(int)
	count  int
	ticker *time.Ticker
}

func newSecondTicker(action func(int)) *secondTicker {
	st := &secondTicker{
		ticker: time.NewTicker(time.Second),
		action: action,
		count:  0,
	}
	st.action(0)
	go func() {
		for range st.ticker.C {
			st.count++
			logger.Debug("count", st.count)
			st.action(st.count)
		}
	}()
	return st
}

func (st *secondTicker) Stop() {
	logger.Debug("Stop")
	if st.ticker != nil {
		st.ticker.Stop()
		st.ticker = nil
	}
	if st.action != nil {
		st.action = nil
	}
}
