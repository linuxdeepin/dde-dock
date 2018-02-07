/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package appearance

import (
	"github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"testing"
)

func TestWriteEnvFile(t *testing.T) {
	convey.Convey("Test env file writer", t, func() {
		var target = `LANG=zh_CN.UTF-8
QT_SCALE_FACTOR=1.35
`
		var tmpFile = "testdata/env_file"
		err := writeKeyToEnvFile("QT_SCALE_FACTOR", "1.35", tmpFile)
		if err != nil {
			println("Failed to create file, skip...")
			return
		}
		convey.So(err, convey.ShouldBeNil)
		content, _ := ioutil.ReadFile(tmpFile)
		convey.So(string(content), convey.ShouldEqual, target)
		writeKeyToEnvFile("QT_SCALE_FACTOR", "1.8", tmpFile)
	})
}
