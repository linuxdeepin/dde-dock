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

package systeminfo

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCPUInfo(t *testing.T) {
	Convey("Test cpu info", t, func() {
		cpu, err := GetCPUInfo("testdata/cpuinfo")
		So(cpu, ShouldEqual,
			"Intel(R) Core(TM) i3 CPU M 330 @ 2.13GHz x 4")
		So(err, ShouldBeNil)

		cpu, err = GetCPUInfo("testdata/sw-cpuinfo")
		So(cpu, ShouldEqual, "sw 1.40GHz x 4")
		So(err, ShouldBeNil)

		cpu, err = GetCPUInfo("testdata/loonson3-cpuinfo")
		So(cpu, ShouldEqual, "Loongson-3B V0.7 FPU V0.1 x 6")
		So(err, ShouldBeNil)

		cpu, err = GetCPUInfo("testdata/arm-cpuinfo")
		So(cpu, ShouldEqual, "NANOPI2 x 4")
		So(err, ShouldBeNil)

		cpu, err = GetCPUInfo("testdata/hw_kirin-cpuinfo")
		So(cpu, ShouldEqual, "HUAWEI Kirin 990 x 8")
		So(err, ShouldBeNil)
	})
}

func TestMemInfo(t *testing.T) {
	Convey("Test memory info", t, func() {
		mem, err := getMemoryFromFile("testdata/meminfo")
		So(mem, ShouldEqual, uint64(4005441536))
		So(err, ShouldBeNil)
	})
}

func TestVersion(t *testing.T) {
	Convey("Test os version", t, func() {
		lang := os.Getenv("LANGUAGE")
		os.Setenv("LANGUAGE", "en_US")
		defer os.Setenv("LANGUAGE", lang)

		deepin, err := getVersionFromDeepin("testdata/deepin-version")
		So(deepin, ShouldEqual, "2015 Desktop Alpha1")
		So(err, ShouldBeNil)

		lsb, err := getVersionFromLSB("testdata/lsb-release")
		So(lsb, ShouldEqual, "2014.3")
		So(err, ShouldBeNil)
	})
}

func TestDistro(t *testing.T) {
	Convey("Test os distro", t, func() {
		lang := os.Getenv("LANGUAGE")
		os.Setenv("LANGUAGE", "en_US")
		defer os.Setenv("LANGUAGE", lang)

		distroId, distroDesc, distroVer, err := getDistroFromLSB("testdata/lsb-release")
		So(distroId, ShouldEqual, "Deepin")
		So(distroDesc, ShouldEqual, "Deepin 2014.3")
		So(distroVer, ShouldEqual, "2014.3")
		So(err, ShouldBeNil)
	})
}

func TestSystemBit(t *testing.T) {
	Convey("Test getconf", t, func() {
		v := systemBit()
		if v != "32" {
			So(v, ShouldEqual, "64")
		}

		if v != "64" {
			So(v, ShouldEqual, "32")
		}
	})
}
