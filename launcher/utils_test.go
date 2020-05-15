/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package launcher

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_getAppIdByFilePath(t *testing.T) {
	Convey("getAppIdByFilePath", t, func(c C) {
		appDirs := []string{"/usr/share/applications", "/usr/local/share/applications", "/home/test_user/.local/share/applications"}

		id := getAppIdByFilePath("/usr/share/applications/d-feet.desktop", appDirs)
		c.So(id, ShouldEqual, "d-feet")

		id = getAppIdByFilePath("/usr/share/applications/kde4/krita.desktop", appDirs)
		c.So(id, ShouldEqual, "kde4/krita")

		id = getAppIdByFilePath("/usr/local/share/applications/deepin-screenshot.desktop", appDirs)
		c.So(id, ShouldEqual, "deepin-screenshot")

		id = getAppIdByFilePath("/home/test_user/.local/share/applications/space test.desktop", appDirs)
		c.So(id, ShouldEqual, "space test")

		id = getAppIdByFilePath("/other/dir/a.desktop", appDirs)
		c.So(id, ShouldEqual, "")
	})
}

func Test_getUserAppDir(t *testing.T) {
	Convey("getUserAppDir", t, func(c C) {
		home := os.Getenv("HOME")
		c.So(getUserAppDir(), ShouldEqual, filepath.Join(home, ".local/share/applications"))
	})
}

func Test_runeSliceDiff(t *testing.T) {
	Convey("runeSliceDiff", t, func(c C) {
		// pop
		popCount, runesPush := runeSliceDiff([]rune("abc"), []rune("abc"))
		c.So(popCount, ShouldEqual, 0)
		c.So(len(runesPush), ShouldEqual, 0)

		popCount, runesPush = runeSliceDiff([]rune("abc"), []rune("abcd"))
		c.So(popCount, ShouldEqual, 1)
		c.So(len(runesPush), ShouldEqual, 0)

		popCount, runesPush = runeSliceDiff([]rune("abc"), []rune("abcde"))
		c.So(popCount, ShouldEqual, 2)
		c.So(len(runesPush), ShouldEqual, 0)

		// push
		popCount, runesPush = runeSliceDiff([]rune("abcd"), []rune("abc"))
		c.So(popCount, ShouldEqual, 0)
		c.So(len(runesPush), ShouldEqual, 1)
		c.So(runesPush[0], ShouldEqual, 'd')

		popCount, runesPush = runeSliceDiff([]rune("abcde"), []rune("abc"))
		c.So(popCount, ShouldEqual, 0)
		c.So(len(runesPush), ShouldEqual, 2)
		c.So(runesPush[0], ShouldEqual, 'd')
		c.So(runesPush[1], ShouldEqual, 'e')

		// pop and push
		popCount, runesPush = runeSliceDiff([]rune("abcd"), []rune("abce"))
		c.So(popCount, ShouldEqual, 1)
		c.So(len(runesPush), ShouldEqual, 1)
		c.So(runesPush[0], ShouldEqual, 'd')

		popCount, runesPush = runeSliceDiff([]rune("deepin"), []rune("deeinp"))
		c.So(popCount, ShouldEqual, 3)
		c.So(len(runesPush), ShouldEqual, 3)
		c.So(runesPush[0], ShouldEqual, 'p')
		c.So(runesPush[1], ShouldEqual, 'i')
		c.So(runesPush[2], ShouldEqual, 'n')
	})
}

func Test_parseFlatpakAppCmdline(t *testing.T) {
	Convey("test parseFlatpakAppCmdline", t, func(c C) {
		info, err := parseFlatpakAppCmdline(`/usr/bin/flatpak run --branch=master --arch=x86_64 --command=blender --file-forwarding org.blender.Blender @@ %f @@`)
		c.So(err, ShouldBeNil)
		c.So(info, ShouldResemble, &flatpakAppInfo{
			name:   "org.blender.Blender",
			arch:   "x86_64",
			branch: "master",
		})
	})
}
