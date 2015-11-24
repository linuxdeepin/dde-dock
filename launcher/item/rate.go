package item

import (
	. "pkg.deepin.io/dde/daemon/launcher/utils"
	"pkg.deepin.io/lib/glib-2.0"
)

const (
	_RateRecordFile = "launcher/rate.ini"
	_RateRecordKey  = "rate"
)

// GetFrequencyRecordFile returns the file which records items' use frequency.
func GetFrequencyRecordFile() (*glib.KeyFile, error) {
	return ConfigFile(_RateRecordFile)
}
