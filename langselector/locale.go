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
	"dbus/com/deepin/api/setdatetime"
	libnetwork "dbus/com/deepin/daemon/network"
	"dbus/org/freedesktop/notifications"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"pkg.linuxdeepin.com/dde-daemon/langselector/language_info"
	"pkg.linuxdeepin.com/lib/log"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"strings"
	"sync"
)

const (
	systemLocaleFile  = "/etc/default/locale"
	userLocaleFilePAM = ".pam_environment"

	defaultLocale = "en_US.UTF-8"
)

const (
	LocaleStateChanged  = 0
	LocaleStateChanging = 1
)

var (
	ErrFileNotExist       = fmt.Errorf("File not exist")
	ErrLocaleNotFound     = fmt.Errorf("Locale not found")
	ErrLocaleChangeFailed = fmt.Errorf("Changing locale failed")
)

type LangSelector struct {
	CurrentLocale string
	Changed       func(locale string)

	LocaleState int32
	logger      *log.Logger
	setDate     *setdatetime.SetDateTime
	lock        sync.Mutex
}

func newLangSelector(l *log.Logger) *LangSelector {
	lang := LangSelector{}

	if l != nil {
		lang.logger = l
	} else {
		lang.logger = log.NewLogger(dbusSender)
	}

	var err error
	lang.setDate, err = setdatetime.NewSetDateTime(
		"com.deepin.api.SetDateTime",
		"/com/deepin/api/SetDateTime")
	if err != nil {
		lang.logger.Warning("New SetDateTime Failed:", err)
		return nil
	}

	lang.setPropCurrentLocale(getLocale())

	return &lang
}

func (lang *LangSelector) Destroy() {
	if lang.setDate == nil {
		return
	}

	setdatetime.DestroySetDateTime(lang.setDate)
	lang.setDate = nil
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
	if err != nil {
		locale, err = getLocaleFromFile(systemLocaleFile)
		if err != nil {
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
	fp, err := os.Create(filename + "~")
	if err != nil {
		return err
	}
	defer fp.Close()

	contents := constructPamFile(locale, filename)
	if _, err = fp.WriteString(contents); err != nil {
		return err
	}
	fp.Sync()
	os.Rename(filename+"~", filename)

	return nil
}

func constructPamFile(locale, filename string) string {
	if !dutils.IsFileExist(filename) {
		return generatePamContents(locale)
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return generatePamContents(locale)
	}

	lines := strings.Split(string(contents), "\n")
	var tmp string
	for i, line := range lines {
		if i != 0 {
			tmp += "\n"
		}

		strs := strings.Split(line, "=")
		if strs[0] == "LANG" {
			tmp += strs[0] + "=" + locale
			continue
		} else if strs[0] == "LANGUAGE" {
			lcode := strings.Split(locale, ".")[0]
			tmp += strs[0] + "=" + lcode
			continue
		}

		tmp += line
	}

	return tmp
}

func generatePamContents(locale string) string {
	contents := ""
	str := "LANG=" + locale + "\n"
	contents += str
	tmp := strings.Split(locale, ".")
	str = "LANGUAGE=" + tmp[0] + "\n"
	contents += str

	return contents
}

func getLocaleFromFile(filename string) (string, error) {
	if !dutils.IsFileExist(filename) {
		return "", ErrFileNotExist
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	var locale string
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
		break
	}

	if len(locale) == 0 {
		return "", ErrLocaleNotFound
	}

	return locale, nil
}
