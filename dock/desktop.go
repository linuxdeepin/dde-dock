package dock

import (
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/glib-2.0"
	"regexp"
)

const (
	ACTION_PATTERN = `(?P<actionGroup>.*) Shortcut Group` // |Desktop Action (.*)
)

var (
	actionReg, _ = regexp.Compile(ACTION_PATTERN)
)

type DesktopAppInfo struct {
	*gio.DesktopAppInfo
	*glib.KeyFile
	gioSupported bool
}

func NewDesktopAppInfo(desktopID string) *DesktopAppInfo {
	dai := &DesktopAppInfo{nil, nil, false}
	dai.DesktopAppInfo = gio.NewDesktopAppInfo(desktopID)
	if dai.DesktopAppInfo == nil {
		return nil
	}

	if len(dai.DesktopAppInfo.ListActions()) != 0 {
		dai.gioSupported = true
	}

	dai.KeyFile = glib.NewKeyFile()
	if ok, _ := dai.LoadFromFile(dai.GetFilename(), glib.KeyFileFlagsNone); !ok {
		dai.Unref()
		return nil
	}

	return dai
}

func NewDesktopAppInfoFromFilename(filename string) *DesktopAppInfo {
	dai := &DesktopAppInfo{nil, nil, false}
	dai.DesktopAppInfo = gio.NewDesktopAppInfoFromFilename(filename)
	if dai.DesktopAppInfo == nil {
		return nil
	}

	if len(dai.DesktopAppInfo.ListActions()) != 0 {
		dai.gioSupported = true
	}

	dai.KeyFile = glib.NewKeyFile()
	if ok, _ := dai.LoadFromFile(dai.GetFilename(), glib.KeyFileFlagsNone); !ok {
		dai.Unref()
		return nil
	}

	return dai
}

func (dai *DesktopAppInfo) ListActions() []string {
	logger.Debug(dai.GetFilename())
	if dai.gioSupported {
		return dai.DesktopAppInfo.ListActions()
	}

	logger.Debug("ListActions .* Shortcut Group")
	actions := make([]string, 0)
	_, groups := dai.GetGroups()
	for _, groupName := range groups {
		if tmp := actionReg.FindStringSubmatch(groupName); len(tmp) > 0 {
			actions = append(actions, tmp[1])
		}
	}

	return actions
}

func getGroupName(gioSupported bool, actionGropuName string) string {
	if gioSupported {
		return "Desktop Action " + actionGropuName
	}
	return actionGropuName + " Shortcut Group"
}

func (dai *DesktopAppInfo) GetActionName(actionGroup string) string {
	if dai.gioSupported {
		logger.Debug("[GetActionName]", dai.GetFilename(), "gio support")
		return dai.DesktopAppInfo.GetActionName(actionGroup)
	}

	logger.Debug("GetActionName")
	langs := GetLanguageNames()
	str := ""
	for _, lang := range langs {
		str, _ = dai.KeyFile.GetLocaleString(getGroupName(dai.gioSupported, actionGroup), "Name", lang)
		if str != "" {
			return str
		}
	}

	if str == "" {
		str, _ = dai.KeyFile.GetString(getGroupName(dai.gioSupported, actionGroup), "Name")
	}

	return str
}

func (dai *DesktopAppInfo) LaunchAction(actionGroup string, ctx gio.AppLaunchContextLike) {
	logger.Debug(dai.GetFilename())
	// LaunchAction won't work for new style desktop, fuck gio.
	// if dai.gioSupported {
	// 	logger.Info("[LaunchAction]", dai.GetFilename(), "gio support")
	// 	dai.DesktopAppInfo.LaunchAction(actionGroup, ctx)
	// 	return
	// }

	logger.Debug("LaunchAction")
	exec, _ := dai.KeyFile.GetString(getGroupName(dai.gioSupported, actionGroup), glib.KeyFileDesktopKeyExec)
	logger.Infof("exec: %q", exec)
	a, err := gio.AppInfoCreateFromCommandline(
		exec,
		"",
		gio.AppInfoCreateFlagsNone,
	)
	if err != nil {
		logger.Warning("Launch App Falied: ", err)
		return
	}

	defer a.Unref()
	_, err = a.Launch(make([]*gio.File, 0), ctx)
	if err != nil {
		logger.Warning("Launch App Failed: ", err)
	}
}

func (dai *DesktopAppInfo) Unref() {
	dai.DesktopAppInfo.Unref()
	dai.KeyFile.Free()
}
