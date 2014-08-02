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
	graphic "pkg.linuxdeepin.com/lib/gdkpixbuf"
	"pkg.linuxdeepin.com/lib/utils"
)

// SetupWrapper is a copy of dde-api/grub2ext, for to remove the dependency
// of dbus when setup grub.
// TODO
type SetupWrapper struct{}

// Setup grub2 environment, regenerate configure and theme if need, don't depends on dbus
func (grub *Grub2) Setup(gfxmode string) {
	setup := &SetupWrapper{}
	runWithoutDBus = true

	grub.config = newConfig()
	grub.config.save()

	// do not call grub.readEntries() here, for that
	// "/boot/grub/grub.cfg" may not exists
	grub.readSettings()
	grub.fixSettings()
	grub.fixSettingDistro()

	// setup gfxmode
	if len(gfxmode) == 0 {
		grub.setSettingGfxmode(grub.config.Resolution)
	} else {
		grub.setSettingGfxmode(gfxmode)
	}

	// write to setting file
	settingFileContent := grub.getSettingContentToSave()
	setup.DoWriteSettings(settingFileContent)

	// setup theme and generate theme background
	grub.SetupTheme(gfxmode)
}

func (grub *Grub2) SetupTheme(gfxmode string) {
	setup := &SetupWrapper{}
	runWithoutDBus = true
	grub.loadConfig()
	if len(gfxmode) == 0 {
		gfxmode = grub.config.Resolution
	}
	w, h := parseGfxmode(gfxmode)
	setup.DoGenerateThemeBackground(w, h)
}

func writeSettingsWithoutDBus(fileContent string) {
	setup := &SetupWrapper{}
	setup.DoWriteSettings(fileContent)
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
func (setup *SetupWrapper) DoGenerateGrubConfig() (ok bool, err error) {
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
func (setup *SetupWrapper) DoSetThemeBackgroundSourceFile(imageFile string, screenWidth, screenHeight uint16) (ok bool, err error) {
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
	err = graphic.ScaleImagePrefer(themeBgSrcFile, themeBgFile, int(screenWidth), int(screenHeight), graphic.GDK_INTERP_HYPER, graphic.FormatPng)
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
