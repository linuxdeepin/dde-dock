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

package grub2

import (
	"io/ioutil"
	"os"
	"path"
	"pkg.linuxdeepin.com/lib/graphic"
	"strconv"
	"strings"
)

// SetupWrapper is a copy of dde-api/grub2ext, for to remove the dependency
// of dbus when setup grub.
type SetupWrapper struct{}

// setup grub2 environment, regenerate configure and theme if need, don't depends on dbus
func (grub *Grub2) Setup(gfxmode string) {
	setup := &SetupWrapper{}
	grub.loadConfig()

	// do not call grub.readEntries() here, for /boot/grub/grub.cfg file maybe
	// not exists
	grub.readSettings()
	grub.fixSettings()

	if len(gfxmode) == 0 {
		grub.setSettingGfxmode(grub.config.Resolution)
	} else {
		grub.setSettingGfxmode(gfxmode)
	}

	settingFileContent := grub.getSettingContentToSave()
	setup.DoWriteSettings(settingFileContent)

	// generate theme background
	grub.SetupTheme(gfxmode)
}

func (grub *Grub2) SetupTheme(gfxmode string) {
	setup := &SetupWrapper{}
	grub.loadConfig()
	if len(gfxmode) == 0 {
		gfxmode = grub.config.Resolution
	}
	w, h := parseGfxmode(gfxmode)
	setup.DoGenerateThemeBackground(w, h)
}

func parseGfxmode(gfxmode string) (w, h uint16) {
	w, h = getPrimaryScreenBestResolution() // default value
	if gfxmode == "auto" {
		return
	}
	a := strings.Split(gfxmode, "x")
	if len(a) != 2 {
		logger.Error("gfxmode format error", gfxmode)
		return
	}

	// parse width
	tmpw, err := strconv.ParseUint(a[0], 10, 16)
	if err != nil {
		logger.Error(err)
		return
	}

	// parse height
	tmph, err := strconv.ParseUint(a[1], 10, 16)
	if err != nil {
		logger.Error(err)
		return
	}

	w = uint16(tmpw)
	h = uint16(tmph)
	return
}

// DoWriteSettings write file content to "/etc/default/grub".
func (setup *SetupWrapper) DoWriteSettings(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(grubConfigFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoWriteCacheConfig write file content to "/var/cache/deepin/grub2.json".
func (setup *SetupWrapper) DoWriteCacheConfig(fileContent string) (ok bool, err error) {
	// ensure parent directory exists
	if !isFileExists(configFile) {
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
func (setup *SetupWrapper) DoGenerateGrubConfig() (ok bool, err error) {
	logger.Info("start to generate a new grub configuration file")
	_, stderr, err := execAndWait(30, grubUpdateExe)
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
func (setup *SetupWrapper) DoSetThemeBackgroundSourceFile(imageFile string, screenWidth, screenHeight uint16) (ok bool, err error) {
	// if background source file is a symlink, just delete it
	if isSymlink(themeBgSrcFile) {
		os.Remove(themeBgSrcFile)
	}

	// backup background source file
	_, err = copyFile(imageFile, themeBgSrcFile)
	if err != nil {
		return false, err
	}

	// generate a new background
	return setup.DoGenerateThemeBackground(screenWidth, screenHeight)
}

// DoGenerateThemeBackground generate the background for deepin grub2
// theme depends on screen resolution.
func (setup *SetupWrapper) DoGenerateThemeBackground(screenWidth, screenHeight uint16) (ok bool, err error) {
	imgWidth, imgHeight, err := graphic.GetImageSize(themeBgSrcFile)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	logger.Infof("source background size %dx%d", imgWidth, imgHeight)
	logger.Infof("background size %dx%d", screenWidth, screenHeight)
	err = graphic.FillImage(themeBgSrcFile, themeBgFile, int(screenWidth), int(screenHeight),
		graphic.FillProportionCenterScale, graphic.PNG)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoCustomTheme write file content to "/boot/grub/themes/deepin/theme.txt".
func (setup *SetupWrapper) DoCustomTheme(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(themeMainFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoWriteThemeJSON write file content to "/boot/grub/themes/deepin/theme_tpl.json".
func (setup *SetupWrapper) DoWriteThemeJSON(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(themeJSONFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoResetThemeBackground link background_origin_source to background_source
func (setup *SetupWrapper) DoResetThemeBackground() (ok bool, err error) {
	os.Remove(themeBgSrcFile)
	err = os.Symlink(themeBgOrigSrcFile, themeBgSrcFile)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}
