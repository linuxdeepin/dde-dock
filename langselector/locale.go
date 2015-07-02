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

package langselector

import (
	"dbus/com/deepin/api/localehelper"
	libnetwork "dbus/com/deepin/daemon/network"
	"dbus/org/freedesktop/notifications"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"pkg.linuxdeepin.com/dde-daemon/langselector/language_info"
	"pkg.linuxdeepin.com/lib/log"
	"strings"
)

const (
	systemLocaleFile  = "/etc/default/locale"
	userLocaleFilePAM = ".pam_environment"

	defaultLocale = "en_US.UTF-8"
)

const (
	// Locale changed state: has been done
	//
	// Locale 更改状态：已经修改完成
	LocaleStateChanged = 0
	// Locale changed state: changing
	//
	// Locale 更改状态：正在修改中
	LocaleStateChanging = 1
)

var (
	// Error: not found the file
	//
	// 错误：没有此文件
	ErrFileNotExist = fmt.Errorf("File not exist")
	// Error: not found the locale
	//
	// 错误：无效的 Locale
	ErrLocaleNotFound = fmt.Errorf("Locale not found")
	// Error: changing locale failure
	//
	// 错误：修改 locale 失败
	ErrLocaleChangeFailed = fmt.Errorf("Changing locale failed")
)

type LangSelector struct {
	// The current locale
	CurrentLocale string
	// Signal: will be emited if locale changed
	Changed func(locale string)

	// Store locale changed state
	LocaleState int32

	logger  *log.Logger
	lhelper *localehelper.LocaleHelper
}

type envInfo struct {
	key   string
	value string
}
type envInfos []envInfo

func newLangSelector(l *log.Logger) *LangSelector {
	lang := LangSelector{LocaleState: LocaleStateChanged}

	if l != nil {
		lang.logger = l
	} else {
		lang.logger = log.NewLogger("dde-daemon/langselector")
	}

	var err error
	lang.lhelper, err = localehelper.NewLocaleHelper(
		"com.deepin.api.LocaleHelper",
		"/com/deepin/api/LocaleHelper")
	if err != nil {
		lang.logger.Warning("New LocaleHelper Failed:", err)
		return nil
	}

	lang.setPropCurrentLocale(getLocale())

	return &lang
}

func (lang *LangSelector) Destroy() {
	if lang.lhelper == nil {
		return
	}

	localehelper.DestroyLocaleHelper(lang.lhelper)
	lang.lhelper = nil
}

func sendNotify(icon, summary, body string) error {
	notifier, err := notifications.NewNotifier(
		"org.freedesktop.Notifications",
		"/org/freedesktop/Notifications")
	if err != nil {
		return err
	}

	_, err = notifier.Notify(dbusSender, 0,
		icon, summary, body,
		nil, nil, 0)

	return err
}

func isNetworkEnable() (bool, error) {
	network, err := libnetwork.NewNetworkManager(
		"com.deepin.daemon.Network",
		"/com/deepin/daemon/Network")
	if err != nil {
		return false, err
	}

	state := network.State.Get()
	// if state < 50, network disconnect
	if state < 50 {
		return false, nil
	}

	return true, nil
}

func getLocale() string {
	filename := path.Join(os.Getenv("HOME"), userLocaleFilePAM)
	locale, err := getLocaleFromFile(filename)
	if err != nil || len(locale) == 0 {
		locale, err = getLocaleFromFile(systemLocaleFile)
		if err != nil || len(locale) == 0 {
			locale = defaultLocale
		}

		writeUserLocale(locale)
	}

	if !language_info.IsLocaleValid(locale,
		language_info.LanguageListFile) {
		locale = defaultLocale
		writeUserLocale(locale)
	}

	return locale
}

func writeUserLocale(locale string) error {
	filename := path.Join(os.Getenv("HOME"), userLocaleFilePAM)
	return writeUserLocalePam(locale, filename)
}

/**
 * gnome locale config
 **/
func writeUserLocalePam(locale, filename string) error {
	var content = generatePamEnvFile(locale, filename)
	return ioutil.WriteFile(filename, []byte(content), 0644)
}

func generatePamEnvFile(locale, filename string) string {
	var (
		lFound   bool //LANG
		lgFound  bool //LANGUAGE
		content  string
		infos, _ = readEnvFile(filename)
		length   = len(infos)
		lang     = strings.Split(locale, ".")[0]
	)
	for i, info := range infos {
		if info.key == "LANG" {
			lFound = true
			info.value = locale
		} else if info.key == "LANGUAGE" {
			lgFound = true
			info.value = lang
		}
		content += fmt.Sprintf("%s=%s", info.key, info.value)
		if i != length-1 {
			content += "\n"
		}
	}
	if !lFound {
		content += fmt.Sprintf("LANG=%s", locale)
		if !lgFound {
			content += "\n"
		}
	}
	if !lgFound {
		content += fmt.Sprintf("LANGUAGE=%s", lang)
	}

	return content
}

func getLocaleFromFile(filename string) (string, error) {
	infos, err := readEnvFile(filename)
	if err != nil {
		return "", err
	}

	var locale string
	for _, info := range infos {
		if info.key != "LANG" {
			continue
		}
		locale = info.value
	}

	return locale, nil
}

func readEnvFile(file string) (envInfos, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var (
		infos envInfos
		lines = strings.Split(string(content), "\n")
	)
	for _, line := range lines {
		var array = strings.Split(line, "=")
		if len(array) != 2 {
			continue
		}

		infos = append(infos, envInfo{
			key:   array[0],
			value: array[1],
		})
	}

	return infos, nil
}
