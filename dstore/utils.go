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
	"dbus/com/deepin/lastore"
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"sync"
	"time"
)

const (
	DStoreDBusDest = "com.deepin.lastore"
	DStoreDBusPath = "/com/deepin/lastore"

	JobStatusSucceed = "succeed"
	JobStatusFailed  = "failed"
	JobStatusEnd     = "end"
)

const (
	jobTypeInstall = "install"
	jobTypeRemove  = "remove"
)

func newDStoreManager() (*lastore.Manager, error) {
	return lastore.NewManager(DStoreDBusDest, DStoreDBusPath)
}

func destroyDStoreManager(manager *lastore.Manager) {
	if manager == nil {
		return
	}
	lastore.DestroyManager(manager)
}

func newDStoreJob(jobPath dbus.ObjectPath) (*lastore.Job, error) {
	return lastore.NewJob(DStoreDBusDest, jobPath)
}

func destroyDStoreJob(job *lastore.Job) {
	if job == nil {
		return
	}
	lastore.DestroyJob(job)
}

func waitJobDone(jobPath dbus.ObjectPath, jobType string, timeout <-chan time.Time, result *(chan error)) {
	job, err := newDStoreJob(jobPath)
	if err != nil {
		*result <- err
		return
	}
	defer destroyDStoreJob(job)

	isQuitFlag := false
	var quitLock sync.Mutex
	setQuit := func() {
		quitLock.Lock()
		defer quitLock.Unlock()
		isQuitFlag = true
	}
	isQuit := func() bool {
		quitLock.Lock()
		defer quitLock.Unlock()
		return isQuitFlag
	}
	quit := make(chan struct{})

	finishJob := func(e error) {
		setQuit()
		// nil must be used explicitly for interface value, otherwise `interfaceValue == nil` will be failed.
		if e == nil {
			*result <- nil
		} else {
			*result <- e
		}
		close(quit)
	}

	job.Status.ConnectChanged(func() {
		status := job.Status.Get()
		switch status {
		case JobStatusSucceed, JobStatusEnd:
			if isQuit() {
				return
			}

			finishJob(nil)
			return
		case JobStatusFailed:
			if isQuit() {
				return
			}

			finishJob(fmt.Errorf(job.Description.Get()))
			return
		default:
			// Only in the case of the installation or removal is successful,
			// the state it may be empty.
			if len(status) == 0 && !isQuit() {
				finishJob(nil)
				return
			}
		}
	})

	select {
	case <-quit:
		return
	case <-timeout:
		setQuit()
		*result <- fmt.Errorf("Do job '%v - %v' timeout",
			jobType, job.Packages.Get())
		return
	}
}

func IsInstalled(pkgName string) bool {
	proxy, err := newDStoreManager()
	if err != nil {
		return false
	}
	defer destroyDStoreManager(proxy)

	installed, _ := proxy.PackageExists(pkgName)
	return installed
}

func IsExists(pkgName string) bool {
	proxy, err := newDStoreManager()
	if err != nil {
		return false
	}
	defer destroyDStoreManager(proxy)

	exists, _ := proxy.PackageInstallable(pkgName)
	return exists
}
