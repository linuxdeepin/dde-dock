package common

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	dbus "github.com/godbus/dbus"
	huawei_fprint "github.com/linuxdeepin/go-dbus-factory/com.huawei.fingerprint"
	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
)

const (
	fprintdDir      = "/var/lib/fprint"
	HuaweiFprintDir = "/var/lib/dde-daemon/fingerprint/huawei"

	HuaweiDeleteTypeOne = 0
	HuaweiDeleteTypeAll = 1
)

func DeleteEnrolledFingers(username, userUuid string) error {
	// remove fprintd fingers
	err := os.RemoveAll(filepath.Join(fprintdDir, username))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// remove huawei fingers
	huaweiFprintUserDir := filepath.Join(HuaweiFprintDir, userUuid)
	fileInfoList, err := ioutil.ReadDir(huaweiFprintUserDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	err = os.RemoveAll(huaweiFprintUserDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if len(fileInfoList) == 0 {
		return nil
	}

	// call reload
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	huaweiDev := huawei_fprint.NewFingerprint(sysBus)

	sysBusDaemon := ofdbus.NewDBus(sysBus)
	hasHuaweiService, err := sysBusDaemon.NameHasOwner(0, huaweiDev.ServiceName_())
	if err != nil {
		return err
	}

	if !hasHuaweiService {
		return nil
	}

	reloadRet, err := huaweiDev.Reload(0, HuaweiDeleteTypeAll)
	if err != nil {
		return err
	}
	if reloadRet == -1 {
		return errors.New("failed to reload")
	}

	return nil
}
