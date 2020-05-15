/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package mime

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAppInfos(t *testing.T) {
	Convey("Delete info", t, func(c C) {
		var infos = AppInfos{
			&AppInfo{
				Id:   "gvim.desktop",
				Name: "gvim",
				Exec: "gvim",
			},
			&AppInfo{
				Id:   "firefox.desktop",
				Name: "Firefox",
				Exec: "firefox",
			}}
		c.So(len(infos.Delete("gvim.desktop")), ShouldEqual, 1)
		c.So(len(infos.Delete("vim.desktop")), ShouldEqual, 2)
	})
}

func TestUnmarshal(t *testing.T) {
	Convey("Test unmarsal", t, func(c C) {
		table, err := unmarshal("testdata/data.json")
		c.So(err, ShouldBeNil)
		c.So(len(table.Apps), ShouldEqual, 2)

		c.So(table.Apps[0].AppId, ShouldResemble, []string{"org.gnome.Nautilus.desktop"})
		c.So(table.Apps[0].AppType, ShouldEqual, "file-manager")
		c.So(table.Apps[0].Types, ShouldResemble, []string{
			"inode/directory",
			"application/x-gnome-saved-search",
		})

		c.So(table.Apps[1].AppId, ShouldResemble, []string{"org.gnome.gedit.desktop"})
		c.So(table.Apps[1].AppType, ShouldEqual, "editor")
		c.So(table.Apps[1].Types, ShouldResemble, []string{
			"text/plain",
		})

	})
}

func TestIsStrInList(t *testing.T) {
	Convey("Test str whether in list", t, func(c C) {
		var list = []string{"abc", "abs"}
		c.So(isStrInList("abs", list), ShouldEqual, true)
		c.So(isStrInList("abd", list), ShouldEqual, false)
	})
}

func TestUserAppInfo(t *testing.T) {
	Convey("User appinfo test", t, func(c C) {
		var infos = userAppInfos{
			{
				DesktopId: "test-web.desktop",
				SupportedMime: []string{
					"application/test.xml",
					"application/test.html",
				},
			},
			{
				DesktopId: "test-doc.desktop",
				SupportedMime: []string{
					"application/test.doc",
					"application/test.xls",
				},
			},
		}
		var file = "testdata/tmp_user_mime.json"
		var manager = &userAppManager{
			appInfos: infos,
			filename: file,
		}
		c.So(manager.Get("application/test.xml")[0].DesktopId, ShouldEqual, "test-web.desktop")
		c.So(manager.Get("application/test.ppt"), ShouldBeNil)
		c.So(manager.Add([]string{"application/test.xml"}, "test-web.desktop"), ShouldEqual, false)
		c.So(manager.Add([]string{"application/test.ppt"}, "test-doc.desktop"), ShouldEqual, true)
		c.So(manager.Get("application/test.ppt")[0].DesktopId, ShouldEqual, "test-doc.desktop")
		c.So(manager.Delete("test-web.desktop"), ShouldBeNil)
		c.So(manager.Delete("test-xxx.desktop"), ShouldNotBeNil)
		c.So(manager.Get("application/test.xml"), ShouldBeNil)
		c.So(manager.Write(), ShouldBeNil)
		tmp, err := newUserAppManager(file)
		c.So(err, ShouldBeNil)
		c.So(tmp.Get("application/test.xml"), ShouldBeNil)
		c.So(tmp.Get("application/test.ppt")[0].DesktopId, ShouldEqual, "test-doc.desktop")
		os.Remove(file)
	})
}
