/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package langselector

import (
	"fmt"
	"pkg.deepin.io/dde/api/lang_info"
	. "pkg.deepin.io/lib/gettext"
)

const (
	localeIconStart    = "notification-change-language-start"
	localeIconFailed   = "notification-change-language-failed"
	localeIconFinished = "notification-change-language-finished"
)

type LocaleInfo struct {
	// Locale name
	Locale string
	// Locale description
	Desc string
}

// Set user desktop environment locale, the new locale will work after relogin.
// (Notice: this locale is only for the current user.)
//
// 设置用户会话的 locale，注销后生效，此改变只对当前用户生效。
//
// locale: see '/etc/locale.gen'
func (lang *LangSelector) SetLocale(locale string) error {
	if lang.LocaleState == LocaleStateChanging {
		return nil
	}

	if len(locale) == 0 || !lang_info.IsSupportedLocale(locale) {
		return fmt.Errorf("Invalid locale: %v", locale)
	}
	if lang.lhelper == nil {
		return fmt.Errorf("LocaleHelper object is nil")
	}

	if lang.CurrentLocale == locale {
		return nil
	}

	go func() {
		lang.LocaleState = LocaleStateChanging
		lang.setPropCurrentLocale(locale)
		if ok, _ := isNetworkEnable(); !ok {
			err := sendNotify(localeIconStart, "",
				Tr("System language is being changed, please wait..."))
			if err != nil {
				lang.logger.Warning("sendNotify failed:", err)
			}
		} else {
			err := sendNotify(localeIconStart, "",
				Tr("System language is being changed with an installation of lacked language packages, please wait..."))
			if err != nil {
				lang.logger.Warning("sendNotify failed:", err)
			}
		}
		err := lang.lhelper.GenerateLocale(locale)
		if err != nil {
			lang.logger.Warning("GenerateLocale failed:", err)
			lang.LocaleState = LocaleStateChanged
		}
	}()

	return nil
}

// Get locale info list that deepin supported
//
// 得到系统支持的 locale 信息列表
func (lang *LangSelector) GetLocaleList() []LocaleInfo {
	list, err := getLocaleInfoList()
	if err != nil {
		lang.logger.Warning(err)
		return nil
	}

	return list
}

func getLocaleInfoList() ([]LocaleInfo, error) {
	infos, err := lang_info.GetSupportedLangInfos()
	if err != nil {
		return nil, err
	}

	var list []LocaleInfo
	for _, info := range infos {
		tmp := LocaleInfo{info.Locale, info.Description}
		list = append(list, tmp)
	}

	return list, nil
}

// Reset set user desktop environment locale to system default locale
func (lang *LangSelector) Reset() error {
	locale, err := getLocaleFromFile(systemLocaleFile)
	if err != nil {
		return err
	}
	return lang.SetLocale(locale)
}
