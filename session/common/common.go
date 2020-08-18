package common

import (
	"errors"
	"fmt"
	"time"

	"github.com/godbus/dbus"
	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
)

func ActivateSysDaemonService(serviceName string) error {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	sysBusObj := ofdbus.NewDBus(sysBus)

	const (
		interval    = 100 * time.Millisecond
		max         = 50
		startErrMax = 10
	)
	startErrCount := 0

	for i := 0; i < max; i++ {
		if startErrCount > startErrMax {
			break
		}

		has, err := sysBusObj.NameHasOwner(0, serviceName)
		if err != nil {
			return err
		}

		if has {
			//fmt.Println("service activated", serviceName)
			return nil
		}

		has, err = sysBusObj.NameHasOwner(0, "com.deepin.daemon.Daemon")
		if err != nil {
			return err
		}

		if !has {
			// dde-system-daemon is not running yet
			//fmt.Println("call start service", serviceName)
			_, err = sysBusObj.StartServiceByName(0, serviceName, 0)
			if err != nil {
				startErrCount++
				fmt.Println(err)
			} else {
				continue
			}
		}

		time.Sleep(interval)
		//fmt.Println("sleep")
	}
	return errors.New("reach max number of retires")
}
