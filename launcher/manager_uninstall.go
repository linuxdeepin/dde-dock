/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package launcher

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
)

// Simply remove the desktop file
func (m *Manager) uninstallDesktopFile(item *Item) error {
	err := os.Remove(item.Path)
	go m.notifyUninstallDone(item, err == nil)
	return err
}

func (m *Manager) notifyUninstallDone(item *Item, succeed bool) {
	//icon := item.Icon
	icon := "deepin-appstore"
	var msgBody string
	if succeed {
		msgBody = fmt.Sprintf(Tr("%q removed successfully."), item.Name)
	} else {
		msgBody = fmt.Sprintf(Tr("Failed to uninstall %q."), item.Name)
	}
	m.notify(icon, "", msgBody)
}

func (m *Manager) uninstallDeepinWineApp(item *Item) error {
	logger.Debug("uninstallDeepinWineApp", item.Path)
	cmd := exec.Command("/opt/deepinwine/tools/uninstall.sh", item.Path)
	err := cmd.Run()
	go m.notifyUninstallDone(item, err == nil)
	return err
}

func (m *Manager) uninstall(id string) error {
	item := m.getItemById(id)
	if item == nil {
		logger.Warning("RequestUninstall failed", errorInvalidID)
		return errorInvalidID
	}

	pkg := m.queryPkgName(item)
	if pkg != "" {
		return m.uninstallSystemPackage(pkg)
	} else {
		// uninstall chrome shortcut
		if item.isChromeShortcut() {
			logger.Debug("item is chrome shortcut")
			return m.uninstallDesktopFile(item)
		}

		// uninstall wine app
		isWineApp, err := item.isWineApp()
		if err != nil {
			return err
		}

		if isWineApp {
			logger.Debug("item is wine app")
			return m.uninstallDeepinWineApp(item)
		}

		return m.uninstallDesktopFile(item)
	}
}

const (
	JobStatusSucceed = "succeed"
	JobStatusFailed  = "failed"
	JobStatusEnd     = "end"
)

func (m *Manager) uninstallSystemPackage(pkg string) error {
	if m.lastore == nil {
		return errors.New("Manager.lastore is nil")
	}

	jobPath, err := m.lastore.RemovePackage("", pkg)
	logger.Debugf("uninstallSystemPackage pkg: %q jobPath: %q", pkg, jobPath)
	if err != nil {
		return err
	}
	return m.waitJobDone(jobPath)
}

func (m *Manager) waitJobDone(jobPath dbus.ObjectPath) error {
	logger.Debug("waitJobDone", jobPath)
	defer logger.Debug("waitJobDone end")
	done := make(chan error)
	go m.monitorJobStatusChange(jobPath, done)
	return <-done
}

func (m *Manager) monitorJobStatusChange(jobPath dbus.ObjectPath, done chan error) {
	con := m.systemDBusConn
	if con == nil {
		err := errors.New("SystemDBusConn is nil")
		logger.Warning(err)
		done <- err
		return
	}
	ch := con.Signal()

	// add match rule
	matchRule := fmt.Sprintf(
		"type='signal',interface='org.freedesktop.DBus.Properties',sender='%s',member='PropertiesChanged',path='%s'", lastoreDBusDest, jobPath)
	logger.Debug("AddMatch", matchRule)
	err := con.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule).Store()
	if err != nil {
		logger.Warning("AddMatch failed:", err)
		done <- err
		return
	}

	// remove match rule
	defer func() {
		err := con.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, matchRule).Store()
		if err != nil {
			logger.Warning("RemoveMatch failed:", err)
		}
		logger.Debug("monitorJobStatusChange return")
	}()

	for v := range ch {
		if v.Name == "org.freedesktop.DBus.Properties.PropertiesChanged" &&
			v.Path == jobPath {
			if len(v.Body) != 3 {
				continue
			}

			ifc, _ := v.Body[0].(string)
			if ifc != "com.deepin.lastore.Job" {
				continue
			}
			props, _ := v.Body[1].(map[string]dbus.Variant)
			status, ok := props["Status"]
			if !ok {
				continue
			}
			statusStr, _ := status.Value().(string)
			logger.Debug("job status changed", statusStr)
			switch statusStr {
			case JobStatusSucceed:
				done <- nil
				return
			case JobStatusFailed:
				done <- errors.New("Job Failed")
				return
			case JobStatusEnd:
				return
			}
		}
	}
}
