/*
 * Copyright (C) 2015 ~ 2017 Deepin Technology Co., Ltd.
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

type DInstallTransaction struct {
	pkgNames        string
	desc            string
	timeoutDuration time.Duration
	timeout         <-chan time.Time
	done            chan error
	disconnect      func()
}

func NewDInstallTransaction(pkgs string, desc string, timeout time.Duration) *DInstallTransaction {
	transaction := &DInstallTransaction{
		pkgNames:        pkgs,
		desc:            desc,
		timeoutDuration: timeout,
		timeout:         nil,
		done:            make(chan error, 1),
	}
	return transaction
}

func (t *DInstallTransaction) run() {
	proxy, err := newDStoreManager()
	if err != nil {
		t.done <- err
		return
	}
	defer destroyDStoreManager(proxy)

	t.timeout = time.After(t.timeoutDuration)
	jobPath, err := proxy.InstallPackage(t.desc, t.pkgNames)
	if err != nil {
		t.done <- err
		return
	}

	go waitJobDone(jobPath, jobTypeInstall, t.timeout, &(t.done))
}

func (t *DInstallTransaction) wait() error {
	err := <-t.done
	close(t.done)
	return err
}

func (t *DInstallTransaction) Exec() error {
	t.run()
	return t.wait()
}
