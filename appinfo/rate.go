package appinfo

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

func GetFrequency(id string, f *glib.KeyFile) uint64 {
	rate, _ := f.GetUint64(id, _RateRecordKey)
	return rate
}

func SetFrequency(id string, freq uint64, f *glib.KeyFile) {
	f.SetUint64(id, _RateRecordKey, freq)
	SaveKeyFile(f, ConfigFilePath(_RateRecordFile))
}
