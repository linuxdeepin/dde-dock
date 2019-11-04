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
	Convey("Delete info", t, func() {
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
		So(len(infos.Delete("gvim.desktop")), ShouldEqual, 1)
		So(len(infos.Delete("vim.desktop")), ShouldEqual, 2)
	})
}

func TestUnmarshal(t *testing.T) {
	Convey("Test unmarsal", t, func() {
		table, err := unmarshal("testdata/data.json")
		So(err, ShouldBeNil)
		So(len(table.Apps), ShouldEqual, 2)

		So(table.Apps[0].AppId, ShouldResemble, []string{"org.gnome.Nautilus.desktop"})
		So(table.Apps[0].AppType, ShouldEqual, "file-manager")
		So(table.Apps[0].Types, ShouldResemble, []string{
			"inode/directory",
			"application/x-gnome-saved-search",
		})

		So(table.Apps[1].AppId, ShouldResemble, []string{"org.gnome.gedit.desktop"})
		So(table.Apps[1].AppType, ShouldEqual, "editor")
		So(table.Apps[1].Types, ShouldResemble, []string{
			"text/plain",
		})

	})
}

func TestMarshal(t *testing.T) {
	Convey("Marshal info", t, func() {
		content, err := toJSON(&AppInfo{
			Id:   "gvim.desktop",
			Name: "gvim",
			Exec: "gvim",
		})
		So(err, ShouldBeNil)
		So(content, ShouldEqual,
			"{\"Id\":\"gvim.desktop\","+
				"\"Name\":\"gvim\","+
				"\"DisplayName\":\"\","+
				"\"Description\":\"\","+
				"\"Icon\":\"\","+
				"\"Exec\":\"gvim\"}")
	})

	Convey("Marshal info list", t, func() {
		content, err := toJSON(AppInfos{
			&AppInfo{
				Id:   "gvim.desktop",
				Name: "gvim",
				Exec: "gvim",
			},
			&AppInfo{
				Id:   "firefox.desktop",
				Name: "Firefox",
				Exec: "firefox",
			},
		})
		So(err, ShouldBeNil)
		So(content, ShouldEqual, "["+
			"{\"Id\":\"gvim.desktop\","+
			"\"Name\":\"gvim\","+
			"\"DisplayName\":\"\","+
			"\"Description\":\"\","+
			"\"Icon\":\"\","+
			"\"Exec\":\"gvim\"},"+
			"{\"Id\":\"firefox.desktop\","+
			"\"Name\":\"Firefox\","+
			"\"DisplayName\":\"\","+
			"\"Description\":\"\","+
			"\"Icon\":\"\","+
			"\"Exec\":\"firefox\"}"+
			"]")
	})

	Convey("Marshal nil", t, func() {
		content, err := toJSON(nil)
		So(content, ShouldEqual, "null")
		So(err, ShouldBeNil)
	})
}

func TestIsStrInList(t *testing.T) {
	Convey("Test str whether in list", t, func() {
		var list = []string{"abc", "abs"}
		So(isStrInList("abs", list), ShouldEqual, true)
		So(isStrInList("abd", list), ShouldEqual, false)
	})
}

func TestUserAppInfo(t *testing.T) {
	Convey("User appinfo test", t, func() {
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
		So(manager.Get("application/test.xml")[0].DesktopId, ShouldEqual, "test-web.desktop")
		So(manager.Get("application/test.ppt"), ShouldBeNil)
		So(manager.Add([]string{"application/test.xml"}, "test-web.desktop"), ShouldEqual, false)
		So(manager.Add([]string{"application/test.ppt"}, "test-doc.desktop"), ShouldEqual, true)
		So(manager.Get("application/test.ppt")[0].DesktopId, ShouldEqual, "test-doc.desktop")
		So(manager.Delete("test-web.desktop"), ShouldBeNil)
		So(manager.Delete("test-xxx.desktop"), ShouldNotBeNil)
		So(manager.Get("application/test.xml"), ShouldBeNil)
		So(manager.Write(), ShouldBeNil)
		tmp, err := newUserAppManager(file)
		So(err, ShouldBeNil)
		So(tmp.Get("application/test.xml"), ShouldBeNil)
		So(tmp.Get("application/test.ppt")[0].DesktopId, ShouldEqual, "test-doc.desktop")
		os.Remove(file)
	})
}
