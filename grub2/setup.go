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

package main

import (
	"dlib/graphic"
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// Setup is a copy of dde-api/grub2ext, for to remove the dependency
// of dbus when setup grub.
type Setup struct{}

// setup grub2 environment, regenerate configure and theme if need, don't depends on dbus
func (grub *Grub2) setup() {
	setup := &Setup{}

	grub.readSettings()
	grub.fixSettings()

	// grub.resetGfxmodeIfNeed()
	grub.resetGfxmode()

	// grub.writeSettings()
	settingFileContent := grub.getSettingContentToSave()
	setup.DoWriteSettings(settingFileContent)

	// grub.writeCacheConfig()
	grub.config.NeedUpdate = true
	cacheFileContent, _ := json.Marshal(grub.config)
	setup.DoWriteCacheConfig(string(cacheFileContent))

	setup.DoGenerateGrubConfig()

	// grub.writeCacheConfig()
	grub.config.NeedUpdate = false
	cacheFileContent, _ = json.Marshal(grub.config)
	setup.DoWriteCacheConfig(string(cacheFileContent))

	// generate theme background
	screenWidth, screenHeight := getPrimaryScreenBestResolution()
	setup.DoGenerateThemeBackground(screenWidth, screenHeight)
}

// DoWriteSettings write file content to "/etc/default/grub".
func (setup *Setup) DoWriteSettings(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(grubConfigFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoWriteCacheConfig write file content to "/var/cache/dde-daemon/grub2.json".
func (setup *Setup) DoWriteCacheConfig(fileContent string) (ok bool, err error) {
	// ensure parent directory exists
	if !isFileExists(grubCacheFile) {
		os.MkdirAll(path.Dir(grubCacheFile), 0755)
	}
	err = ioutil.WriteFile(grubCacheFile, []byte(fileContent), 0644)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoGenerateGrubConfig execute command "/usr/sbin/update-grub" to
// generate a new grub configuration.
func (setup *Setup) DoGenerateGrubConfig() (ok bool, err error) {
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
func (setup *Setup) DoSetThemeBackgroundSourceFile(imageFile string, screenWidth, screenHeight uint16) (ok bool, err error) {
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
func (setup *Setup) DoGenerateThemeBackground(screenWidth, screenHeight uint16) (ok bool, err error) {
	imgWidth, imgHeight, err := graphic.GetImageSize(themeBgSrcFile)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	logger.Infof("source background size %dx%d", imgWidth, imgHeight)

	x0, y0, x1, y1 := getImgClipRectByResolution(screenWidth, screenHeight, imgWidth, imgHeight)
	logger.Infof("background clip rect (%d,%d), (%d,%d)", x0, y0, x1, y1)
	logger.Infof("background size %dx%d", x1-x0, y1-y0)
	err = graphic.ClipImage(themeBgSrcFile, themeBgFile, x0, y0, x1, y1, graphic.PNG)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoCustomTheme write file content to "/boot/grub/themes/deepin/theme.txt".
func (setup *Setup) DoCustomTheme(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(themeMainFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoWriteThemeJSON write file content to "/boot/grub/themes/deepin/theme_tpl.json".
func (setup *Setup) DoWriteThemeJSON(fileContent string) (ok bool, err error) {
	err = ioutil.WriteFile(themeJSONFile, []byte(fileContent), 0664)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// DoResetThemeBackground link background_origin_source to background_source
func (setup *Setup) DoResetThemeBackground() (ok bool, err error) {
	os.Remove(themeBgSrcFile)
	err = os.Symlink(themeBgOrigSrcFile, themeBgSrcFile)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}
