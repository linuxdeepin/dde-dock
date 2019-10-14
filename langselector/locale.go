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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	// dbus services:
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.api.localehelper"
	libnetwork "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.network"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.lastore"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.notifications"

	"pkg.deepin.io/dde/api/lang_info"
	"pkg.deepin.io/dde/api/language_support"
	"pkg.deepin.io/dde/api/userenv"
	ddbus "pkg.deepin.io/dde/daemon/dbus"
	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/gsettings"
	"pkg.deepin.io/lib/strv"
	"pkg.deepin.io/lib/xdg/basedir"
)

const (
	systemLocaleFile     = "/etc/default/locale"
	systemdLocaleFile    = "/etc/locale.conf"
	userLocaleConfigFile = ".config/locale.conf"

	defaultLocale = "en_US.UTF-8"
)

var (
	// for locale-helper
	_ = Tr("Authentication is required to switch language")
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

	gsSchemaLocale = "com.deepin.dde.locale"
	gsKeyLocales   = "locales"
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
	systemBus    *dbus.Conn
	helper       *localehelper.LocaleHelper
	localesCache LocaleInfos

	PropsMu sync.RWMutex
	// The current locale
	CurrentLocale string
	// Store locale changed state
	LocaleState int32
	// dbusutil-gen: equal=nil
	Locales  []string
	settings *gio.Settings

	methods *struct {
		SetLocale                  func() `in:"locale"`
		GetLocaleList              func() `out:"locales"`
		GetLocaleDescription       func() `in:"locale" out:"description"`
		GetLanguageSupportPackages func() `in:"locale" out:"packages"`
		AddLocale                  func() `in:"locale"`
		DeleteLocale               func() `in:"locale"`
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

func newLangSelector(service *dbusutil.Service) (*LangSelector, error) {
	lang := LangSelector{
		service:     service,
		LocaleState: LocaleStateChanged,
	}

	systemBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	lang.systemBus = systemBus
	lang.helper = localehelper.NewLocaleHelper(systemBus)

	locale := getCurrentUserLocale()
	if !lang.isSupportedLocale(locale) {
		logger.Warningf("newLangSelector: get invalid locale %q", locale)
		locale = defaultLocale
	}
	lang.CurrentLocale = locale

	lang.settings = gio.NewSettings(gsSchemaLocale)
	locales := lang.settings.GetStrv(gsKeyLocales)
	if !strv.Strv(locales).Contains(locale) {
		locales = append(locales, locale)
		lang.settings.SetStrv(gsKeyLocales, locales)
	}
	lang.Locales = locales

	return &lang, nil
}

func (l *LangSelector) connectSettingsChanged() {
	gsettings.ConnectChanged(gsSchemaLocale, gsKeyLocales, func(key string) {
		locales := l.settings.GetStrv(gsKeyLocales)
		l.updateLocales(locales)
	})
}

func (l *LangSelector) updateLocales(locales []string) {
	l.PropsMu.Lock()
	defer l.PropsMu.Unlock()

	if !strv.Strv(locales).Contains(l.CurrentLocale) {
		locales = append(locales, l.CurrentLocale)
	}
	l.setPropLocales(locales)
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
	if locale == "" {
		return false
	}
	infos := ls.getCachedLocales()
	_, err := infos.Get(locale)
	return err == nil
}

func sendNotify(icon, summary, body string) {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		logger.Warning(err)
		return
	}
	n := notifications.NewNotifications(sessionBus)
	_, err = n.Notify(0, dbusServiceName, 0,
		icon, summary, body,
		nil, nil, 0)
	logger.Debugf("send notification icon: %q, summary: %q, body: %q",
		icon, summary, body)

	if err != nil {
		logger.Warning(err)
	}
	return
}

func isNetworkEnable() (bool, error) {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return false, err
	}
	network := libnetwork.NewNetwork(sessionBus)
	state, err := network.State().Get(0)
	if err != nil {
		return false, err
	}
	// if state < 50, network disconnect
	if state < 50 {
		return false, nil
	}

	return true, nil
}

func getCurrentUserLocale() (locale string) {
	locale, _ = userenv.Get("LANG")
	if locale == "" {
		files := []string{
			systemLocaleFile,
			systemdLocaleFile, // It is used by systemd to store system-wide locale settings
		}

		for _, file := range files {
			locale, _ = getLocaleFromFile(file)
			if locale != "" {
				// get locale success
				break
			}
		}
	}
	if locale == "" {
		return defaultLocale
	}
	return locale
}

func writeUserLocale(locale string) error {
	homeDir := basedir.GetUserHomeDir()
	err := userenv.Modify(func(m map[string]string) {
		m["LANG"] = locale
		m["LANGUAGE"] = strings.Split(locale, ".")[0]
	})
	if err != nil {
		return err
	}

	localeConfigFile := filepath.Join(homeDir, userLocaleConfigFile)
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

func (lang *LangSelector) setLocaleFailed(oldLocale string) {
	sendNotify(localeIconFailed, "",
		Tr("Failed to change system language, please try later"))
	// restore CurrentLocale
	lang.PropsMu.Lock()
	lang.setPropCurrentLocale(oldLocale)
	lang.setPropLocaleState(LocaleStateChanged)
	lang.PropsMu.Unlock()
}

func (lang *LangSelector) setLocale(locale string) {
	// begin
	lang.PropsMu.Lock()
	oldLocale := lang.CurrentLocale
	lang.setPropLocaleState(LocaleStateChanging)
	lang.setPropCurrentLocale(locale)
	lang.PropsMu.Unlock()

	// send notification
	networkEnabled, err := isNetworkEnable()
	if err != nil {
		logger.Warning(err)
	}

	if networkEnabled {
		sendNotify(localeIconStart, "",
			Tr("Changing system language and installing the required language packages, please wait..."))
	} else {
		sendNotify(localeIconStart, "",
			Tr("Changing system language, please wait..."))
	}

	// generate locale
	err = lang.generateLocale(locale)
	if err != nil {
		logger.Warning("failed to generate locale:", err)
		lang.setLocaleFailed(oldLocale)
		return
	} else {
		logger.Debug("generate locale success")
	}

	err = writeUserLocale(locale)
	if err != nil {
		logger.Warning("failed to write user locale:", err)
		lang.setLocaleFailed(oldLocale)
		return
	}

	// sync user locale to accounts daemon
	err = syncUserLocale(locale)
	if err != nil {
		logger.Warning("failed to sync user locale to accounts daemon:", err)
	}

	// remove font config file
	fontCfgFile := filepath.Join(basedir.GetUserConfigDir(),
		"fontconfig/conf.d/99-deepin.conf")
	err = os.Remove(fontCfgFile)
	if err != nil {
		// ignore not exist error
		if !os.IsNotExist(err) {
			logger.Warningf("failed to remove font config file: %v", err)
		}
	}

	// install language support packages
	if networkEnabled {
		err = lang.installLangSupportPackages(locale)
		if err != nil {
			logger.Warning("failed to install packages:", err)
		} else {
			logger.Debug("install packages success")
		}
	}

	sendNotify(localeIconFinished, "",
		Tr("System language changed, please log out and then log in"))

	// end
	lang.PropsMu.Lock()
	lang.setPropLocaleState(LocaleStateChanged)
	lang.PropsMu.Unlock()
}

func syncUserLocale(locale string) error {
	systemConn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	u, err := ddbus.NewUserByUid(systemConn, currentUser.Uid)
	if err != nil {
		return err
	}
	err = u.SetLocale(0, locale)
	return err
}

var errSignalBodyInvalid = errors.New("signal body is invalid")

func (lang *LangSelector) generateLocale(locale string) error {
	successMatchRule := dbusutil.NewMatchRuleBuilder().ExtSignal("/com/deepin/api/LocaleHelper", "com.deepin.api.LocaleHelper", "Success").Build()
	err := successMatchRule.AddTo(lang.systemBus)
	if err != nil {
		return err
	}
	sigChan := make(chan *dbus.Signal, 1)
	lang.systemBus.Signal(sigChan)

	defer func() {
		lang.systemBus.RemoveSignal(sigChan)
		err := successMatchRule.RemoveFrom(lang.systemBus)
		if err != nil {
			logger.Warning(err)
		}
	}()

	logger.Debug("generating locale")
	err = lang.helper.GenerateLocale(0, locale)
	if err != nil {
		return err
	}

	select {
	case <-time.NewTimer(10 * time.Minute).C:
		return errors.New("wait success signal timed out")
	case sig := <-sigChan:
		if len(sig.Body) != 2 {
			return errSignalBodyInvalid
		}
		genLocaleOk, ok := sig.Body[0].(bool)
		if !ok {
			return errSignalBodyInvalid
		}

		failReason, ok := sig.Body[1].(string)
		if !ok {
			return errSignalBodyInvalid
		}

		if genLocaleOk {
			return nil
		} else {
			return errors.New(failReason)
		}
	}
}

func (lang *LangSelector) installLangSupportPackages(locale string) error {
	logger.Debug("install language support packages for locale", locale)
	ls, err := language_support.NewLanguageSupport()
	if err != nil {
		return err
	}

	pkgs := ls.ByLocale(locale, false)
	ls.Destroy()
	logger.Info("need to install:", pkgs)
	return lang.installPackages(pkgs)
}

func (lang *LangSelector) installPackages(pkgs []string) error {
	if len(pkgs) == 0 {
		return nil
	}
	systemBus := lang.systemBus
	lastoreObj := lastore.NewLastore(systemBus)
	jobPath, err := lastoreObj.InstallPackage(0, "",
		strings.Join(pkgs, " "))
	if err != nil {
		return err
	}
	logger.Debug("install job path:", jobPath)

	jobMatchRule := dbusutil.NewMatchRuleBuilder().ExtPropertiesChanged(
		string(jobPath), "com.deepin.lastore.Job").Build()
	err = jobMatchRule.AddTo(systemBus)
	if err != nil {
		return err
	}

	sigChan := make(chan *dbus.Signal, 10)
	systemBus.Signal(sigChan)

	defer func() {
		systemBus.RemoveSignal(sigChan)
		err := jobMatchRule.RemoveFrom(systemBus)
		if err != nil {
			logger.Warning(err)
		}
	}()

	for sig := range sigChan {
		if sig.Path == jobPath &&
			sig.Name == "org.freedesktop.DBus.Properties.PropertiesChanged" {
			if len(sig.Body) != 3 {
				return errSignalBodyInvalid
			}

			props, ok := sig.Body[1].(map[string]dbus.Variant)
			if !ok {
				return errSignalBodyInvalid
			}
			v, ok := props["Progress"]
			if ok {
				progress, _ := v.Value().(float64)
				if progress == 1 {
					// install success
					return nil
				}
			}

			v, ok = props["Status"]
			if ok {
				status, _ := v.Value().(string)
				if status == "failed" {
					return errors.New("install failed")
				}
			}
		}
	}
	return nil
}

func (lang *LangSelector) addLocale(locale string) error {
	locales := lang.settings.GetStrv(gsKeyLocales)
	if !strv.Strv(locales).Contains(locale) {
		locales = append(locales, locale)
	}

	lang.settings.SetStrv(gsKeyLocales, locales)
	return nil
}

func (lang *LangSelector) deleteLocale(locale string) error {
	locales := lang.settings.GetStrv(gsKeyLocales)
	locales, isDeleted := strv.Strv(locales).Delete(locale)
	if isDeleted {
		lang.settings.SetStrv(gsKeyLocales, locales)
	}
	return nil
}
