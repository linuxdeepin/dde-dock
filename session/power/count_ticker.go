/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
