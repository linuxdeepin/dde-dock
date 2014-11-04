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

package langselect

import (
	"dbus/org/freedesktop/notifications"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"pkg.linuxdeepin.com/lib/dbus"
	. "pkg.linuxdeepin.com/lib/gettext"
	"pkg.linuxdeepin.com/lib/glib-2.0"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"strings"
)

const (
	_PAM_ENV_FILE   = ".pam_environment"
	_DMRC_FILE      = ".dmrc"
	_DMRC_KEY_GROUP = "Desktop"

	_DEFAULT_LOCALE_FILE = "/etc/default/locale"
	_ALL_LANG_FILE       = "/usr/share/dde-daemon/lang/all_languages"
)

type localeInfo struct {
	Locale string
	Desc   string
}

func (ls *LangSelect) sendNotify(icon, summary, body string) {
	notifier, err := notifications.NewNotifier("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	if err != nil {
		Logger.Warning("New Notifier Failed:", err)
		return
	}

	notifier.Notify(DEST, 0, icon, summary, body, nil, nil, 0)
}

func (ls *LangSelect) getDefaultLocale() (locale string, ok bool) {
	locale = ""
	ok = false

	if !dutils.IsFileExist(_DEFAULT_LOCALE_FILE) {
		Logger.Warningf("'%s' not exist", _DEFAULT_LOCALE_FILE)
		return
	}

	contents, err := ioutil.ReadFile(_DEFAULT_LOCALE_FILE)
	if err != nil {
		Logger.Warningf("ReadFile '%s' failed: %v", _DEFAULT_LOCALE_FILE, err)
		return
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		strs := strings.Split(line, "=")
		if len(strs) != 2 {
			continue
		}

		if strs[0] != "LANG" {
			continue
		}

		locale = strings.Trim(strs[1], "\"")
		ok = true
		break
	}

	return
}

func (ls *LangSelect) getUserLocale() (locale string, ok bool) {
	locale = ""
	ok = false

	homeDir := dutils.GetHomeDir()
	filePath := path.Join(homeDir, _PAM_ENV_FILE)
	if !dutils.IsFileExist(filePath) {
		Logger.Warningf("'%s' not exist", filePath)
		return
	}

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		Logger.Warningf("ReadFile '%s' failed: %v", filePath, err)
		return
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		strs := strings.Split(line, "=")
		if len(strs) != 2 {
			continue
		}

		if strs[0] != "LANG" {
			continue
		}

		locale = strs[1]
		ok = true
		break
	}

	return
}

func (ls *LangSelect) writeUserDmrc(locale string) error {
	homeDir := dutils.GetHomeDir()
	if len(homeDir) < 1 {
		return errors.New("Get home dir failed")
	}

	filePath := path.Join(homeDir, _DMRC_FILE)
	dutils.WriteKeyToKeyFile(filePath, _DMRC_KEY_GROUP,
		"LANG", locale)
	dutils.WriteKeyToKeyFile(filePath, _DMRC_KEY_GROUP,
		"LANGUAGE", locale)

	return nil
}

func (ls *LangSelect) genUserPamContent(locale string) (contents string) {
	contents = ""

	tmp := "LANG=" + locale + "\n"
	contents += tmp
	tmp = "LANGUAGE=" + locale + "\n"
	contents += tmp

	return
}

func (ls *LangSelect) writeUserPamEnv(locale string) error {
	homeDir := dutils.GetHomeDir()
	if len(homeDir) < 1 {
		return errors.New("Get home dir failed")
	}

	filePath := path.Join(homeDir, _PAM_ENV_FILE)
	fp, err := os.Create(filePath + "~")
	if err != nil {
		Logger.Warningf("Create '%s' failed: %v", filePath+"~", err)
		return err
	}
	defer fp.Close()

	if _, err = fp.WriteString(ls.genUserPamContent(locale)); err != nil {
		Logger.Warningf("Write '%s' failed: %v", filePath+"~", err)
		return err
	}
	fp.Sync()
	os.Rename(filePath+"~", filePath)

	return nil
}

func (ls *LangSelect) getLocaleInfo(kf *glib.KeyFile, l, locale string) (info localeInfo, ok bool) {
	if kf == nil {
		return
	}

	lang := os.Getenv("LANG")
	if len(lang) < 1 {
		lang = "en_GB"
	} else {
		strs := strings.Split(lang, ".")
		if dutils.IsElementInList(strs[0], localeList) {
			lang = strs[0]
		} else {
			tmps := strings.Split(strs[0], "_")
			if dutils.IsElementInList(tmps[0], localeList) {
				lang = tmps[0]
			} else {
				lang = "en_GB"
			}
		}
	}

	v, err := kf.GetLocaleString(l, "Name", lang)
	if err != nil {
		return
	}
	info.Locale = locale
	info.Desc = v
	ok = true

	return
}

func (ls *LangSelect) getLocaleInfoList() (list []localeInfo) {
	keyFile := glib.NewKeyFile()
	defer keyFile.Free()
	if _, err := keyFile.LoadFromFile(_ALL_LANG_FILE,
		glib.KeyFileFlagsKeepComments|
			glib.KeyFileFlagsKeepTranslations); err != nil {
		return
	}

	for n, v := range localeListMap {
		info, ok := ls.getLocaleInfo(keyFile, n, v)
		if !ok {
			continue
		}

		list = append(list, info)
	}

	return
}

func (ls *LangSelect) listenLocaleChange() {
	setDate.ConnectGenLocaleStatus(func(ok bool, locale string) {
		if ok && ls.changeLocaleFlag {
			//setLocaleDmrc(locale)
			ls.writeUserPamEnv(locale)
			ls.changeLocaleFlag = false
		}
		ls.setPropCurrentLocale(ls.getCurrentLocale())
		dbus.Emit(ls, "LocaleStatus", ok, locale)
		if ok {
			ls.sendNotify("", "", Tr("Language has been changed successfully and will be effective after logged out."))
		} else {
			ls.sendNotify("", "", Tr("Language failed to change, please try later."))
		}
	})
}
