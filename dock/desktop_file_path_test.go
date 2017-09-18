/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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

package dock

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var deskotpFilePathTestMap = map[string]string{
	"/usr/share/applications/deepin-screenshot.desktop":                               "/S@deepin-screenshot",
	"/usr/local/share/applications/wps-office-et.desktop":                             "/L@wps-office-et",
	"/home/tp/.config/dock/scratch/docked:w:42f9e4a33162e38b2febbad0d9e39a3f.desktop": "/D@docked:w:42f9e4a33162e38b2febbad0d9e39a3f",
	"/home/tp/.local/share/applications/webtorrent-desktop.desktop":                   "/H@webtorrent-desktop",
}

func init() {
	homeDir = "/home/tp/"
	scratchDir = homeDir + ".config/dock/scratch/"
	initPathDirCodeMap()
}

func Test_addDesktopExt(t *testing.T) {
	Convey("addDesktopExt", t, func() {
		So(addDesktopExt("0ad"), ShouldEqual, "0ad.desktop")
		So(addDesktopExt("0ad.desktop"), ShouldEqual, "0ad.desktop")
		So(addDesktopExt("0ad.desktop-x"), ShouldEqual, "0ad.desktop-x.desktop")
	})
}

func Test_trimDesktopExt(t *testing.T) {
	Convey("trimDesktopExt", t, func() {
		So(trimDesktopExt("deepin-movie"), ShouldEqual, "deepin-movie")
		So(trimDesktopExt("deepin-movie.desktop"), ShouldEqual, "deepin-movie")
		So(trimDesktopExt("deepin-movie.desktop-x"), ShouldEqual, "deepin-movie.desktop-x")
	})
}

func Test_zipDesktopPath(t *testing.T) {
	Convey("zipDesktopPath", t, func() {
		for path, zipped := range deskotpFilePathTestMap {
			So(zipped, ShouldEqual, zipDesktopPath(path))
		}
	})
}

func Test_unzipDesktopPath(t *testing.T) {
	Convey("unzipDesktopPath", t, func() {
		for path, zipped := range deskotpFilePathTestMap {
			So(path, ShouldEqual, unzipDesktopPath(zipped))
		}
	})
}

func Test_getDesktopIdByFilePath(t *testing.T) {
	Convey("getDesktopIdByFilePath", t, func() {
		path := "/usr/share/applications/deepin-screenshot.desktop"
		desktopId := getDesktopIdByFilePath(path)
		So(desktopId, ShouldEqual, "deepin-screenshot.desktop")

		path = "/usr/share/applications/kde4/krita.desktop"
		desktopId = getDesktopIdByFilePath(path)
		So(desktopId, ShouldEqual, "kde4-krita.desktop")

		path = "/home/tp/.local/share/applications/telegramdesktop.desktop"
		desktopId = getDesktopIdByFilePath(path)
		So(desktopId, ShouldEqual, "telegramdesktop.desktop")

		path = "/home/tp/.local/share/applications/dirfortest/dir2/space test.desktop"
		desktopId = getDesktopIdByFilePath(path)
		So(desktopId, ShouldEqual, "dirfortest-dir2-space test.desktop")
	})
}
