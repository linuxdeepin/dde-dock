/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package langselector

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"dbus/com/deepin/api/localehelper"
	libnetwork "dbus/com/deepin/daemon/network"
	"dbus/org/freedesktop/notifications"

	"pkg.deepin.io/dde/api/lang_info"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

const (
	systemLocaleFile     = "/etc/default/locale"
	systemdLocaleFile    = "/etc/locale.conf"
	userLocaleFilePAM    = ".pam_environment"
	userLocaleConfigFile = ".config/locale.conf"

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

//go:generate dbusutil-gen -type LangSelector locale.go
type LangSelector struct {
	service      *dbusutil.Service
	logger       *log.Logger
	helper       *localehelper.LocaleHelper
	localesCache LocaleInfos

	PropsMaster dbusutil.PropsMaster
	// The current locale
	CurrentLocale string
	// Store locale changed state
	LocaleState int32

	methods *struct {
		SetLocale            func() `in:"locale"`
		GetLocaleList        func() `out:"locales"`
		GetLocaleDescription func() `in:"locale" out:"description"`
	}

	signals *struct {
		Changed struct {
			locale string
		}
	}
}

type envInfo struct {
	key   string
	value string
}
type envInfos []envInfo

type LocaleInfo struct {
	// Locale name
	Locale string
	// Locale description
	Desc string
}

type LocaleInfos []LocaleInfo

func (infos LocaleInfos) Get(locale string) (LocaleInfo, error) {
	for _, info := range infos {
		if info.Locale == locale {
			return info, nil
		}
	}
	return LocaleInfo{}, fmt.Errorf("invalid locale %q", locale)
}

func newLangSelector(service *dbusutil.Service) *LangSelector {
	lang := LangSelector{
		service:     service,
		LocaleState: LocaleStateChanged,
		logger:      logger,
	}

	var err error
	lang.helper, err = localehelper.NewLocaleHelper(
		"com.deepin.api.LocaleHelper",
		"/com/deepin/api/LocaleHelper")
	if err != nil {
		lang.logger.Warning("New LocaleHelper Failed:", err)
		return nil
	}

	locale := getCurrentUserLocale()
	if !lang.isSupportedLocale(locale) {
		logger.Warningf("newLangSelector: get invalid locale %q", locale)
		locale = defaultLocale
	}
	lang.CurrentLocale = locale
	return &lang
}

func getLocaleInfos() (LocaleInfos, error) {
	infos, err := lang_info.GetSupportedLangInfos()
	if err != nil {
		return nil, err
	}

	list := make(LocaleInfos, len(infos))
	for i, info := range infos {
		list[i] = LocaleInfo{
			Locale: info.Locale,
			Desc:   info.Description,
		}
	}
	return list, nil
}

func (ls *LangSelector) getCachedLocales() LocaleInfos {
	if ls.localesCache == nil {
		var err error
		ls.localesCache, err = getLocaleInfos()
		if err != nil {
			logger.Warning("getLocaleInfos failed:", err)
		}
	}
	return ls.localesCache
}

func (ls *LangSelector) isSupportedLocale(locale string) bool {
	infos := ls.getCachedLocales()
	_, err := infos.Get(locale)
	return err == nil
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

func getCurrentUserLocale() (locale string) {
	files := [3]string{
		path.Join(os.Getenv("HOME"), userLocaleFilePAM),
		systemLocaleFile,
		systemdLocaleFile, // It is used by systemd to store system-wide locale settings
	}

	var err error
	for _, file := range files {
		locale, err = getLocaleFromFile(file)
		if err == nil && locale != "" {
			// get locale success
			break
		}
	}
	if locale == "" {
		return defaultLocale
	}
	return locale
}

func writeUserLocale(locale string) error {
	homeDir := os.Getenv("HOME")
	pamEnvFile := path.Join(homeDir, userLocaleFilePAM)
	var err error
	// only for lightdm
	err = writeLocaleEnvFile(locale, pamEnvFile)
	if err != nil {
		return err
	}
	localeConfigFile := path.Join(homeDir, userLocaleConfigFile)
	err = writeLocaleEnvFile(locale, localeConfigFile)
	if err != nil {
		return err
	}
	return nil
}

func writeLocaleEnvFile(locale, filename string) error {
	var content = generateLocaleEnvFile(locale, filename)
	return ioutil.WriteFile(filename, content, 0644)
}

func generateLocaleEnvFile(locale, filename string) []byte {
	var (
		langFound     bool
		languageFound bool
		infos, _      = readEnvFile(filename)
		lang          = strings.Split(locale, ".")[0]
		buf           bytes.Buffer
	)
	for _, info := range infos {
		if info.key == "LANG" {
			langFound = true
			info.value = locale
		} else if info.key == "LANGUAGE" {
			languageFound = true
			info.value = lang
		}
		buf.WriteString(fmt.Sprintf("%s=%s\n", info.key, info.value))
	}
	if !langFound {
		buf.WriteString(fmt.Sprintf("LANG=%s\n", locale))
	}
	if !languageFound {
		buf.WriteString(fmt.Sprintf("LANGUAGE=%s\n", lang))
	}

	return buf.Bytes()
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

	locale = strings.Trim(locale, " '\"")
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
