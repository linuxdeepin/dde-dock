package item

import (
	. "pkg.linuxdeepin.com/dde-daemon/launcher/interfaces"
	. "pkg.linuxdeepin.com/dde-daemon/launcher/utils"
)

const (
	_RateRecordFile = "launcher/rate.ini"
	_RateRecordKey  = "rate"
)

func GetFrequencyRecordFile() (RateConfigFileInterface, error) {
	return ConfigFile(_RateRecordFile, "")
}
