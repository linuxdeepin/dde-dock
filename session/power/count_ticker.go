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
	exit     chan struct{}
}

func newCountTicker(interval time.Duration, action func(int)) *countTicker {
	t := &countTicker{
		interval: interval,
		action:   action,
	}
	t.Reset()
	return t
}

func (t *countTicker) Reset() {
	t.ticker = time.NewTicker(t.interval)
	t.count = 0
	t.action(0)
	t.exit = make(chan struct{})
	go func() {
		for {
			select {
			case <-t.ticker.C:
				t.count++
				logger.Debug("tick", t.count)
				t.action(t.count)
			case <-t.exit:
				t.exit = nil
				return
			}
		}
	}()
}

func (t *countTicker) Stop() {
	if t.ticker != nil {
		logger.Debug("Stop")
		t.ticker.Stop()
	}
	if t.exit != nil {
		logger.Debug("exit")
		close(t.exit)
	}
}
