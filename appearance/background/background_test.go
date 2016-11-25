/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package background

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestScanner(t *testing.T) {
	Convey("getBgFilesInDir", t, func() {
		So(getBgFilesInDir("testdata/Theme1/wallpapers"), ShouldResemble,
			[]string{
				"testdata/Theme1/wallpapers/desktop.jpg",
			})
		So(getBgFilesInDir("testdata/Theme2/wallpapers"), ShouldBeNil)
	})
}

func TestFileInDirs(t *testing.T) {
	Convey("Test file whether in dirs", t, func() {
		var dirs = []string{
			"/tmp/backgrounds",
			"/tmp/wallpapers",
		}

		So(isFileInSpecialDir("/tmp/backgrounds/1.jpg", dirs),
			ShouldEqual, true)
		So(isFileInSpecialDir("/tmp/wallpapers/1.jpg", dirs),
			ShouldEqual, true)
		So(isFileInSpecialDir("/tmp/background/1.jpg", dirs),
			ShouldEqual, false)
	})
}

func TestGetBgFiles(t *testing.T) {
	files := getBgFiles()
	t.Log(files)
}
