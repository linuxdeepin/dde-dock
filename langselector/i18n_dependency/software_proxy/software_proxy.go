/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package software_proxy

import (
	"dbus/com/linuxdeepin/softwarecenter"
	"pkg.linuxdeepin.com/lib/dbus"
	"time"
)

var (
	actionEndReasons = []string{
		"pkg-installed",
		"pkg-not-in-cache",
		"parse-download-error",
		"download-failed",
		"download-stop",
		"action-finish",
		"action-failed",
	}
)

type SoftwareProxy struct {
	softCenter    *softwarecenter.SoftwareCenter
	listenPkgs    []string
	actionEndChan chan struct{}
}

func NewSoftwareProxy() (*SoftwareProxy, error) {
	softProxy := &SoftwareProxy{}

	var err error
	softProxy.softCenter, err = softwarecenter.NewSoftwareCenter(
		"com.linuxdeepin.softwarecenter",
		"/com/linuxdeepin/softwarecenter",
	)
	if err != nil {
		return nil, err
	}

	softProxy.listenPackageChangeSignal()

	return softProxy, nil
}

func (softProxy *SoftwareProxy) Destroy() {
	softwarecenter.DestroySoftwareCenter(softProxy.softCenter)
}

func (softProxy *SoftwareProxy) IsPackageInstalled(pkg string) bool {
	if len(pkg) == 0 {
		return true
	}

	ret, _ := softProxy.softCenter.GetPkgInstalled(pkg)
	if ret == 1 {
		return true
	}

	return false
}

func (softProxy *SoftwareProxy) IsPackageExist(pkg string) bool {
	if len(pkg) == 0 {
		return false
	}

	list, _ := softProxy.softCenter.IsPkgInCache(pkg)
	if isStrInList(pkg, list) {
		return true
	}

	return false
}

func (softProxy *SoftwareProxy) InstallPackage(packages []string) error {
	for _, pkg := range packages {
		err := softProxy.softCenter.InstallPkg([]string{pkg})
		if err != nil {
			return err
		}
	}

	return nil
}

func (softProxy *SoftwareProxy) UninstallPackage(packages []string) error {
	var purge = true
	for _, pkg := range packages {
		err := softProxy.softCenter.UninstallPkg(pkg, purge)
		return err
	}

	return nil
}

func (softProxy *SoftwareProxy) SetListenPackages(packages []string) {
	softProxy.actionEndChan = make(chan struct{})
	softProxy.listenPkgs = packages
}

func (softProxy *SoftwareProxy) WaitActionEnd() {
	select {
	// Install software timeout
	case <-time.After(time.Minute * 30):
		return
	case <-softProxy.actionEndChan:
		return
	}
}

func (softProxy *SoftwareProxy) EndAction() {
	close(softProxy.actionEndChan)
	softProxy.listenPkgs = nil
}

// pkg install or uninstall
func (softProxy *SoftwareProxy) listenPackageChangeSignal() {
	var cnt int

	softProxy.softCenter.Connectupdate_signal(func(messages [][]interface{}) {
		defer func() {
			err := recover()
			if err != nil {
				softProxy.EndAction()
			}
		}()

		if len(softProxy.listenPkgs) == 0 {
			return
		}

		for _, msg := range messages {
			if msg == nil {
				continue
			}

			var action = msg[0].(string)
			if !isStrInList(action, actionEndReasons) {
				continue
			}

			detail := msg[1].(dbus.Variant).Value().([]interface{})
			pkgName := detail[0].(string)

			if !isStrInList(pkgName, softProxy.listenPkgs) {
				continue
			}

			cnt += 1
		}

		if cnt == len(softProxy.listenPkgs) {
			cnt = 0
			softProxy.EndAction()
		}
	})
}

func isStrInList(str string, list []string) bool {
	for _, v := range list {
		if str == v {
			return true
		}
	}

	return false
}
