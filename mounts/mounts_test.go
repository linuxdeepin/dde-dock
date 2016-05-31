package mounts

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetIdByMountPoint(t *testing.T) {
	convey.Convey("Test get id by mountPoint", t, func() {
		var value = "egcbhjcd123"
		convey.So(getIdByMountPoint("afc://"+value+"/"), convey.ShouldEqual, value)
		convey.So(getIdByMountPoint("mtp://"+value+"/"), convey.ShouldEqual, "mtp://"+value+"/")
	})
}

func TestDiskType(t *testing.T) {
	convey.Convey("Test disk type", t, func() {
		var info = DiskInfo{CanEject: true, Icon: "disk-usb-device"}
		info.correctDiskType()
		convey.So(info.Type, convey.ShouldEqual, DiskTypeRemovable)

		info.Icon = "iphone"
		info.correctDiskType()
		convey.So(info.Type, convey.ShouldEqual, DiskTypeIPhone)

		info.Icon = "phone"
		info.correctDiskType()
		convey.So(info.Type, convey.ShouldEqual, DiskTypePhone)

		info.Icon = "camera-photo"
		info.correctDiskType()
		convey.So(info.Type, convey.ShouldEqual, DiskTypeCamera)

		info.Icon = "disk-dvd-device"
		info.correctDiskType()
		convey.So(info.Type, convey.ShouldEqual, DiskTypeDVD)

		info.MountPoint = "smb://share"
		info.correctDiskType()
		convey.So(info.Type, convey.ShouldEqual, DiskTypeNetwork)
	})
}

func TestStringStartWith(t *testing.T) {
	convey.Convey("Test stringStartWith", t, func() {
		convey.So(stringStartWith("afc://edadsfseweds12dsc", "afc://"),
			convey.ShouldEqual, true)
		convey.So(stringStartWith("afc://edadsfseweds12dsc", "afcc://"),
			convey.ShouldEqual, false)
		convey.So(stringStartWith("afc://", "afcc://"),
			convey.ShouldEqual, false)
	})
}
