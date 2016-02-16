/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
