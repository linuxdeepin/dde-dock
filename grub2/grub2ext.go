/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
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

// This module is split from dde-daemon/grub2 to fix launch issue
// through dbus-daemon for that system bus in root couldn't access
// session bus interface.

package grub2

import (
	"io/ioutil"
	"os"
	"path"
	"pkg.linuxdeepin.com/lib/dbus"
	graphic "pkg.linuxdeepin.com/lib/gdkpixbuf"
	"pkg.linuxdeepin.com/lib/utils"
)

const (
	DbusGrub2ExtDest = "com.deepin.daemon.Grub2Ext"
	DbusGrub2ExtPath = "/com/deepin/daemon/Grub2Ext"
	DbusGrub2ExtIfs  = "com.deepin.daemon.Grub2Ext"
)

// Grub2Ext is a dbus object that split from dde-daemon/grub2 to fix
// issue that system bus in root permission couldn't access session
// bus's interface.
type Grub2Ext struct{}

// NewGrub2Ext create a Grub2Ext object.
func NewGrub2Ext() *Grub2Ext {
	grub := &Grub2Ext{}
	return grub
}

// GetDBusInfo implement interface of dbus.DBusObject
func (ge *Grub2Ext) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DbusGrub2ExtDest,
		DbusGrub2ExtPath,
		DbusGrub2ExtIfs,
	}
}

// TODO rename to DoWriteGrubSettings
// DoWriteSettings write file content to "/etc/default/grub".
func (ge *Grub2Ext) DoWriteSettings(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(grubConfigFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// TODO rename to DoWriteConfig
// DoWriteCacheConfig write file content to "/var/cache/deepin/grub2.json".
func (ge *Grub2Ext) DoWriteCacheConfig(fileContent string) (ok bool, err error) {
	// ensure parent directory exists
	if !utils.IsFileExist(configFile) {
		os.MkdirAll(path.Dir(configFile), 0755)
	}
	err = ioutil.WriteFile(configFile, []byte(fileContent), 0644)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoGenerateGrubConfig execute command "/usr/sbin/update-grub" to
// generate a new grub configuration.
func (ge *Grub2Ext) DoGenerateGrubConfig() (ok bool, err error) {
	logger.Info("start to generate a new grub configuration file")
	_, stderr, err := utils.ExecAndWait(30, grubUpdateExe)
	logger.Infof("process output: %s", stderr)
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
func (ge *Grub2Ext) DoSetThemeBackgroundSourceFile(imageFile string, screenWidth, screenHeight uint16) (ok bool, err error) {
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
	return ge.DoGenerateThemeBackground(screenWidth, screenHeight)
}

// DoGenerateThemeBackground generate the background for deepin grub2
// theme depends on screen resolution.
func (ge *Grub2Ext) DoGenerateThemeBackground(screenWidth, screenHeight uint16) (ok bool, err error) {
	imgWidth, imgHeight, err := graphic.GetImageSize(themeBgSrcFile)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	logger.Infof("source background size %dx%d", imgWidth, imgHeight)
	logger.Infof("background size %dx%d", screenWidth, screenHeight)
	err = graphic.ScaleImagePrefer(themeBgSrcFile, themeBgFile, int(screenWidth), int(screenHeight), graphic.GDK_INTERP_HYPER, graphic.FormatPng)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoCustomTheme write file content to "/boot/grub/themes/deepin/theme.txt".
func (ge *Grub2Ext) DoCustomTheme(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(themeMainFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoWriteThemeJSON write file content to "/boot/grub/themes/deepin/theme_tpl.json".
func (ge *Grub2Ext) DoWriteThemeJSON(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(themeJSONFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoResetThemeBackground link background_origin_source to background_source
func (ge *Grub2Ext) DoResetThemeBackground() (ok bool, err error) {
	os.Remove(themeBgSrcFile)
	err = os.Symlink(themeBgOrigSrcFile, themeBgSrcFile)
	if err != nil {
		logger.Error(err.Error()) // TODO
		return false, err
	}
	return true, nil
}
