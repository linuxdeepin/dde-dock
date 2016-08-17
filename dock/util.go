/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"encoding/base64"
	"gir/gio-2.0"
	"io"
	"os"
	"pkg.deepin.io/dde/daemon/appinfo"
	"time"
)

func trimDesktop(desktopID string) string {
	desktopIDLen := len(desktopID)
	if desktopIDLen == 0 {
		return ""
	}

	if desktopIDLen > 8 {
		return strings.TrimSuffix(desktopID, ".desktop")
	}

	panic(fmt.Sprintf("%q is not a desktop id", desktopID))
}

func normalizeAppID(candidateID string) string {
	return strings.Replace(appinfo.NormalizeAppID(candidateID), " ", "-", -1)
}

var _DesktopAppIdReg = regexp.MustCompile(`(?:[^.]+\.)*(?P<desktopID>[^.]+)\.desktop`)

func getAppIDFromDesktopID(candidateID string) string {
	desktopID := guess_desktop_id(candidateID)
	logger.Debug(fmt.Sprintf("get desktop id: %q", desktopID))
	if desktopID == "" {
		return ""
	}

	appID := normalizeAppID(trimDesktop(desktopID))
	return appID
}

// the key is appID
// the value is desktopID
var _appIDCache map[string]string = make(map[string]string)

func guess_desktop_id(appId string) string {
	logger.Debugf("guess_desktop_id %q", appId)
	if desktopID, ok := _appIDCache[appId]; ok {
		logger.Debug(appId, "is in cache")
		return desktopID
	}

	desktopID := appId + ".desktop"
	allApp := gio.AppInfoGetAll()

	defer func() {
		for _, app := range allApp {
			app.Unref()
		}
	}()

	for _, app := range allApp {
		_appInfo := gio.ToDesktopAppInfo(app)

		if _appInfo == nil {
			continue
		}

		_desktopID := _appInfo.GetId()
		normalizedDesktopID := normalizeAppID(_desktopID)
		if strings.HasSuffix(normalizedDesktopID, desktopID) {
			_appIDCache[appId] = _desktopID
			return _desktopID
		}

		// TODO: this is not a silver bullet, fix it later.
		appIDs := _DesktopAppIdReg.FindStringSubmatch(normalizedDesktopID)
		if len(appIDs) == 2 && appIDs[1] == appId {
			_appIDCache[appId] = _desktopID
			return _desktopID
		}
	}

	return ""
}

func dataUriToFile(dataUri, path string) (string, error) {
	commaIndex := strings.Index(dataUri, ",")
	img, err := base64.StdEncoding.DecodeString(dataUri[commaIndex+1:])
	if err != nil {
		return path, err
	}

	return path, ioutil.WriteFile(path, img, 0744)
}

func strSliceEqual(sa, sb []string) bool {
	if len(sa) != len(sb) {
		return false
	}
	for i, va := range sa {
		vb := sb[i]
		if va != vb {
			return false
		}
	}
	return true
}

func uniqStrSlice(slice []string) []string {
	newSlice := make([]string, 0)
	for _, e := range slice {
		if !strSliceContains(newSlice, e) {
			newSlice = append(newSlice, e)
		}
	}
	return newSlice
}

func strSliceContains(slice []string, v string) bool {
	for _, e := range slice {
		if e == v {
			return true
		}
	}
	return false
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func getCurrentTimestamp() uint32 {
	return uint32(time.Now().Unix())
}

func recordFrequency(appId string) {
	f, err := appinfo.GetFrequencyRecordFile()
	if err == nil {
		appinfo.SetFrequency(appId, appinfo.GetFrequency(appId, f)+1, f) // FIXME: DesktopID???
		f.Free()
	}
}
