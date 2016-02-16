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
