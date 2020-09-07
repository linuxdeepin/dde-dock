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
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCPUInfo(t *testing.T) {
	Convey("Test cpu info", t, func(c C) {
		cpu, err := GetCPUInfo("testdata/cpuinfo")
		c.So(cpu, ShouldEqual,
			"Intel(R) Core(TM) i3 CPU M 330 @ 2.13GHz x 4")
		c.So(err, ShouldBeNil)

		cpu, err = GetCPUInfo("testdata/sw-cpuinfo")
		c.So(cpu, ShouldEqual, "sw 1.40GHz x 4")
		c.So(err, ShouldBeNil)

		cpu, err = GetCPUInfo("testdata/arm-cpuinfo")
		c.So(cpu, ShouldEqual, "NANOPI2 x 4")
		c.So(err, ShouldBeNil)

		cpu, err = GetCPUInfo("testdata/hw_kirin-cpuinfo")
		c.So(cpu, ShouldEqual, "HUAWEI Kirin 990 x 8")
		c.So(err, ShouldBeNil)
	})
}

func TestMemInfo(t *testing.T) {
	Convey("Test memory info", t, func(c C) {
		mem, err := getMemoryFromFile("testdata/meminfo")
		c.So(mem, ShouldEqual, uint64(4005441536))
		c.So(err, ShouldBeNil)
	})
}

func TestVersion(t *testing.T) {
	Convey("Test os version", t, func(c C) {
		lang := os.Getenv("LANGUAGE")
		os.Setenv("LANGUAGE", "en_US")
		defer os.Setenv("LANGUAGE", lang)

		deepin, err := getVersionFromDeepin("testdata/deepin-version")
		c.So(deepin, ShouldEqual, "2015 Desktop Alpha1")
		c.So(err, ShouldBeNil)

		lsb, err := getVersionFromLSB("testdata/lsb-release")
		c.So(lsb, ShouldEqual, "2014.3")
		c.So(err, ShouldBeNil)
	})
}

func TestDistro(t *testing.T) {
	Convey("Test os distro", t, func(c C) {
		lang := os.Getenv("LANGUAGE")
		os.Setenv("LANGUAGE", "en_US")
		defer os.Setenv("LANGUAGE", lang)

		distroId, distroDesc, distroVer, err := getDistroFromLSB("testdata/lsb-release")
		c.So(distroId, ShouldEqual, "Deepin")
		c.So(distroDesc, ShouldEqual, "Deepin 2014.3")
		c.So(distroVer, ShouldEqual, "2014.3")
		c.So(err, ShouldBeNil)
	})
}

func TestSystemBit(t *testing.T) {
	Convey("Test getconf", t, func(c C) {
		v := systemBit()
		if v != "32" {
			c.So(v, ShouldEqual, "64")
		}

		if v != "64" {
			c.So(v, ShouldEqual, "32")
		}
	})
}

func TestIsFloatEqual(t *testing.T) {
	Convey("Test memory info", t, func(c C) {
		c.So(isFloatEqual(0.001, 0.0), ShouldEqual, false)
		c.So(isFloatEqual(0.001, 0.001), ShouldEqual, true)
	})
}

func TestParseInfoFile(t *testing.T) {
	Convey("Test ParseInfoFile", t, func(c C) {
		v, err := parseInfoFile("testdata/lsb-release", "=")
		c.So(err, ShouldBeNil)
		c.So(v["DISTRIB_ID"], ShouldEqual, "Deepin")
		c.So(v["DISTRIB_RELEASE"], ShouldEqual, "2014.3")
		c.So(v["DISTRIB_DESCRIPTION"], ShouldEqual, strconv.Quote("Deepin 2014.3"))
	})
}

func TestGetCPUMaxMHzByLscpu(t *testing.T) {
	Convey("Test GetCPUMaxMHzByLscpu", t, func(c C) {
		ret, err := parseInfoFile("testdata/lsCPU", ":")
		c.So(err, ShouldBeNil)
		v, err := getCPUMaxMHzByLscpu(ret)
		c.So(err, ShouldBeNil)
		c.So(v, ShouldEqual, 3600.0000)
	})
}

func TestGetProcessorByLscpuu(t *testing.T) {
	Convey("Test GetProcessorByLscpu", t, func(c C) {
		ret, err := parseInfoFile("testdata/lsCPU", ":")
		c.So(err, ShouldBeNil)
		v, err := getProcessorByLscpu(ret)
		c.So(err, ShouldBeNil)
		c.So(v, ShouldEqual,"Intel(R) Core(TM) i5-4570 CPU @ 3.20GHz x 4")
	})
}

func TestDoReadCache(t *testing.T) {
	Convey("Test DoReadCache", t, func(c C) {
		ret, err := doReadCache("testdata/systeminfo.cache")
		c.So(err, ShouldBeNil)
		c.So(ret.Version, ShouldEqual, "20 专业版")
		c.So(ret.DistroID, ShouldEqual, "uos")
		c.So(ret.DistroDesc, ShouldEqual, "UnionTech OS 20")
		c.So(ret.DistroVer, ShouldEqual, "20")
		c.So(ret.Processor, ShouldEqual, "Intel(R) Core(TM) i5-4570 CPU @ 3.20GHz x 4")
		c.So(ret.DiskCap, ShouldEqual, uint64(500107862016))
		c.So(ret.MemoryCap, ShouldEqual, uint64(8280711168))
		c.So(ret.SystemType, ShouldEqual, 64)
		c.So(ret.CPUMaxMHz, ShouldEqual, 3600)
		c.So(ret.CurrentSpeed, ShouldEqual, 0)
	})
}
