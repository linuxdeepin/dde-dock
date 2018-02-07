/*
 * Copyright (C) 2015 ~ 2018 Deepin Technology Co., Ltd.
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

package dstore

import (
	"time"
)

// DUninstallTransaction is command object for uninstalling package.
// TODO: add Cancel
type DUninstallTransaction struct {
	pkgName         string
	purge           bool // this is useless for new interface.
	timeoutDuration time.Duration
	timeout         <-chan time.Time
	done            chan error
	disconnect      func()
}

// NewDUninstallTransaction creates a new DUninstallTransaction.
func NewDUninstallTransaction(pkgName string, purge bool, timeoutDuration time.Duration) *DUninstallTransaction {
	return &DUninstallTransaction{
		pkgName:         pkgName,
		purge:           purge,
		timeoutDuration: timeoutDuration,
		timeout:         nil,
		done:            make(chan error, 1),
	}
}

func (t *DUninstallTransaction) run() {
	proxy, err := newDStoreManager()
	if err != nil {
		t.done <- err
		return
	}
	defer destroyDStoreManager(proxy)

	t.timeout = time.After(t.timeoutDuration)
	jobPath, err := proxy.RemovePackage("", t.pkgName)
	if err != nil {
		t.done <- err
		return
	}

	go waitJobDone(jobPath, jobTypeRemove, t.timeout, &(t.done))
}

func (t *DUninstallTransaction) wait() error {
	err := <-t.done
	close(t.done)
	return err
}

// Exec executes this transaction.
func (t *DUninstallTransaction) Exec() error {
	t.run()
	return t.wait()
}
