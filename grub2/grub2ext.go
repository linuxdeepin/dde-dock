/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

// This module is split from dde/daemon/grub2 to fix launch issue
// through dbus-daemon for that system bus in root couldn't access
// session bus interface.

package grub2

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/polkit"
	"pkg.deepin.io/lib/utils"
)

const (
	DbusGrub2ExtDest = "com.deepin.daemon.Grub2Ext"
	DbusGrub2ExtPath = "/com/deepin/daemon/Grub2Ext"
	DbusGrub2ExtIfs  = "com.deepin.daemon.Grub2Ext"
)

// Grub2Ext is a dbus object that split from dde/daemon/grub2 to fix
// issue that system bus in root permission couldn't access session
// bus's interface.
type Grub2Ext struct{}

// NewGrub2Ext create a Grub2Ext object.
func NewGrub2Ext() *Grub2Ext {
	polkit.Init()
	grub := &Grub2Ext{}
	return grub
}

// GetDBusInfo implement interface of dbus.DBusObject
func (ge *Grub2Ext) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       DbusGrub2ExtDest,
		ObjectPath: DbusGrub2ExtPath,
		Interface:  DbusGrub2ExtIfs,
	}
}

func checkAuthWithPid(pid uint32) (bool, error) {
	subject := polkit.NewSubject(polkit.SubjectKindUnixProcess)
	subject.SetDetail("pid", pid)
	subject.SetDetail("start-time", uint64(0))
	const actionId = DbusGrub2ExtDest
	details := make(map[string]string)
	details[""] = ""
	result, err := polkit.CheckAuthorization(subject, actionId, details,
		polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}

	return result.IsAuthorized, nil
}

var errAuthFailed = errors.New("authentication failed")

func checkAuth(dbusMsg dbus.DMessage) error {
	pid := dbusMsg.GetSenderPID()
	isAuthorized, err := checkAuthWithPid(pid)
	if err != nil {
		return err
	}
	if !isAuthorized {
		return errAuthFailed
	}
	return nil
}

// DoWriteConfig write file content to "/var/cache/deepin/grub2.json".
func (ge *Grub2Ext) DoWriteConfig(dbusMsg dbus.DMessage, fileContent string) (ok bool, err error) {
	logger.Debug("Grub2Ext.DoWriteConfig")
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	err = doWriteConfig([]byte(fileContent))
	if err != nil {
		logger.Error(err)
		return false, err
	}
	return true, nil
}

// DoWriteGrubSettings write file content to "/etc/default/grub".
func (ge *Grub2Ext) DoWriteGrubSettings(dbusMsg dbus.DMessage, fileContent string) (ok bool, err error) {
	logger.Debug("Grub2Ext.DoWriteGrubSettings")
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	err = doWriteGrubSettings(fileContent)
	if err != nil {
		logger.Error(err)
		return false, err
	}
	return true, nil
}

// DoGenerateGrubMenu execute command "/usr/sbin/update-grub" to
// generate a new grub configuration.
func (ge *Grub2Ext) DoGenerateGrubMenu(dbusMsg dbus.DMessage) (ok bool, err error) {
	logger.Debug("Grub2Ext.DoGenerateGrubMenu")
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}
	logger.Info("start to generate a new grub configuration file")
	// force use LANG=en_US.UTF-8 to make lsb-release/os-probe support
	// Unicode characters
	// FIXME: keep same with the current system language settings
	os.Setenv("LANG", "en_US.UTF-8")

	oldPathEnv := os.Getenv("PATH")
	os.Setenv("PATH", "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin")
	defer os.Setenv("PATH", oldPathEnv)

	cmd := exec.Command(grubMkconfigCmd, "-o", "/boot/grub/grub.cfg")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	logger.Infof("process error output: %s", stderr)
	if err != nil {
		logger.Errorf("generate grub configuration failed: %v", err)
		return false, err
	}
	logger.Info("generate grub configuration successful")
	return true, nil
}

// DoSetThemeBackgroundSourceFile setup a new background source file
// for deepin grub2 theme, and then generate the background depends on
// screen resolution.
func (ge *Grub2Ext) DoSetThemeBackgroundSourceFile(dbusMsg dbus.DMessage, imageFile string, screenWidth, screenHeight uint16) (ok bool, err error) {
	logger.Debug("Grub2Ext.DoSetThemeBackgroundSourceFile")
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}
	// if background source file is a symlink, just delete it
	if utils.IsSymlink(themeBgSrcFile) {
		os.Remove(themeBgSrcFile)
	}

	// backup background source file
	err = utils.CopyFile(imageFile, themeBgSrcFile)
	if err != nil {
		return false, err
	}

	// generate a new background
	return ge.DoGenerateThemeBackground(dbusMsg, screenWidth, screenHeight)
}

// DoGenerateThemeBackground generate the background for deepin grub2
// theme depends on screen resolution.
func (ge *Grub2Ext) DoGenerateThemeBackground(dbusMsg dbus.DMessage, screenWidth, screenHeight uint16) (ok bool, err error) {
	logger.Debug("Grub2Ext.DoGenerateThemeBackground")
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	err = doGenerateThemeBackground(screenWidth, screenHeight)
	if err != nil {
		logger.Error(err)
		return false, err
	}
	return true, nil
}

// DoWriteThemeMainFile write file content to "/boot/grub/themes/deepin/theme.txt".
func (ge *Grub2Ext) DoWriteThemeMainFile(dbusMsg dbus.DMessage, fileContent string) (ok bool, err error) {
	logger.Debug("Grub2Ext.DoWriteThemeMainFile")
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(themeMainFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err)
		return false, err
	}
	return true, nil
}

// DoWriteThemeTplFile write file content to "/boot/grub/themes/deepin/theme_tpl.json".
func (ge *Grub2Ext) DoWriteThemeTplFile(dbusMsg dbus.DMessage, fileContent string) (ok bool, err error) {
	logger.Debug("Grub2Ext.DoWriteThemeTplFile")
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(themeJSONFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err)
		return false, err
	}
	return true, nil
}

// DoResetThemeBackground link background_origin_source to background_source
func (ge *Grub2Ext) DoResetThemeBackground(dbusMsg dbus.DMessage) (ok bool, err error) {
	logger.Debug("Grub2Ext.DoResetThemeBackground")
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	os.Remove(themeBgSrcFile)
	err = os.Symlink(themeBgOrigSrcFile, themeBgSrcFile)
	if err != nil {
		logger.Error(err)
		return false, err
	}
	return true, nil
}
