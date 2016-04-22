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

type countTicker struct {
	action   func(int)
	interval time.Duration
	count    int
	ticker   *time.Ticker
}

func newCountTicker(interval time.Duration, action func(int)) *countTicker {
	st := &countTicker{
		interval: interval,
		action:   action,
	}
	st.Reset()
	return st
}

func (st *countTicker) Reset() {
	st.ticker = time.NewTicker(st.interval)
	st.count = 0
	st.action(0)
	go func() {
		for range st.ticker.C {
			st.count++
			st.action(st.count)
		}
	}()
}

func (st *countTicker) Stop() {
	logger.Debug("Stop")
	if st.ticker != nil {
		st.ticker.Stop()
		st.ticker = nil
	}
}
