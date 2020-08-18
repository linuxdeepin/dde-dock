package grub_gfx

import (
	"github.com/godbus/dbus"
	ofd "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	"pkg.deepin.io/dde/daemon/grub_common"
)

func detectChange() {
	if !grub_common.HasDeepinGfxmodeMod() {
		logger.Debug("not found grub module deepin_gfxmode")
		return
	}

	params, err := grub_common.LoadGrubParams()
	if err != nil {
		logger.Warning(err)
	}
	if grub_common.ShouldFinishGfxmodeDetect(params) {
		logger.Debug("finish gfxmode detect")
		err = startSysGrubService()
		if err != nil {
			logger.Warning("failed to start sys-grub service:", err)
		}
		return
	}
	if grub_common.InGfxmodeDetectionMode(params) {
		logger.Debug("in gfxmode detection mode")
		return
	}

	currentGfxmode, allGrubGfxmodes, err := grub_common.GetBootArgDeepinGfxmode()
	if err != nil {
		logger.Warning("failed to get boot arg DEEPIN_GFXMODE:", err)
		if !grub_common.IsGfxmodeDetectFailed(params) {
			err = prepareGfxmodeDetect()
			if err != nil {
				logger.Warning(err)
			}
		}
		return
	}
	logger.Debug("currentGfxmode:", currentGfxmode)
	logger.Debug("allGrubGfxmodes:", allGrubGfxmodes)

	randrGfxmodes, err := grub_common.GetGfxmodesFromXRandr()
	if err != nil {
		logger.Warning(err)
		return
	}

	logger.Debug("randrGfxmodes:", randrGfxmodes)

	maxGfxmode := randrGfxmodes.Intersection(allGrubGfxmodes).Max()
	logger.Debug("maxGfxmode:", maxGfxmode)

	var maxNotSupported bool
	if params[grub_common.DeepinGfxmodeNotSupported] == maxGfxmode.String() {
		maxNotSupported = true
	}
	logger.Debug("maxNotSupported:", maxNotSupported)

	cfgGfxmodeStr := grub_common.DecodeShellValue(params["GRUB_GFXMODE"])
	logger.Debug("cfgGfxmodeStr:", cfgGfxmodeStr)
	cfgGfxmode, cfgGfxmodeErr := grub_common.ParseGfxmode(cfgGfxmodeStr)
	if cfgGfxmodeErr != nil {
		logger.Warning("failed to parse cfgGfxmodeStr:", cfgGfxmodeErr)
	} else {
		logger.Debug("cfgGfxmode:", cfgGfxmode)
	}
	need := needDetect(cfgGfxmode, cfgGfxmodeErr, currentGfxmode, maxGfxmode, maxNotSupported)
	logger.Debug("need detect:", need)
	if need {
		err = prepareGfxmodeDetect()
		if err != nil {
			logger.Warning(err)
		}
	}
}

func needDetect(cfgGfxmode grub_common.Gfxmode, cfgGfxmodeErr error,
	currentGfxmode, maxGfxmode grub_common.Gfxmode, maxNotSupported bool) bool {

	return cfgGfxmodeErr != nil ||
		cfgGfxmode != currentGfxmode ||
		(currentGfxmode != maxGfxmode && !maxNotSupported)
}

func startSysGrubService() error {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	sysBusDaemon := ofd.NewDBus(sysBus)
	_, err = sysBusDaemon.StartServiceByName(dbus.FlagNoAutoStart,
		"com.deepin.daemon.Grub2", 0)
	return err
}

func getSysGrubObj() (dbus.BusObject, error) {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}

	obj := sysBus.Object("com.deepin.daemon.Grub2", "/com/deepin/daemon/Grub2")
	return obj, nil
}

func prepareGfxmodeDetect() error {
	logger.Debug("prepare gfxmode detect")
	sysGrubObj, err := getSysGrubObj()
	if err != nil {
		return err
	}

	return sysGrubObj.Call("com.deepin.daemon.Grub2.PrepareGfxmodeDetect", 0).Err
}
