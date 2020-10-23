/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package dock

import (
	"path/filepath"
	"strconv"
	"strings"

	"pkg.deepin.io/lib/procfs"
)

type IdentifyWindowFunc struct {
	Name string
	Fn   _IdentifyWindowFunc
}

type _IdentifyWindowFunc func(*Manager, *WindowInfo) (string, *AppInfo)

func (m *Manager) registerIdentifyWindowFuncs() {
	m.registerIdentifyWindowFunc("PidEnv", identifyWindowByPidEnv)
	m.registerIdentifyWindowFunc("CmdlineTurboBooster", identifyWindowByCmdlineTurboBooster)
	m.registerIdentifyWindowFunc("Cmdline-XWalk", identifyWindowByCmdlineXWalk)
	m.registerIdentifyWindowFunc("FlatpakAppID", identifyWindowByFlatpakAppID)
	m.registerIdentifyWindowFunc("CrxId", identifyWindowByCrxId)
	m.registerIdentifyWindowFunc("Rule", identifyWindowByRule)
	m.registerIdentifyWindowFunc("Bamf", identifyWindowByBamf)
	m.registerIdentifyWindowFunc("Pid", identifyWindowByPid)
	m.registerIdentifyWindowFunc("Scratch", identifyWindowByScratch)
	m.registerIdentifyWindowFunc("GtkAppId", identifyWindowByGtkAppId)
	m.registerIdentifyWindowFunc("WmClass", identifyWindowByWmClass)
}

func (m *Manager) registerIdentifyWindowFunc(name string, fn _IdentifyWindowFunc) {
	m.identifyWindowFuns = append(m.identifyWindowFuns, &IdentifyWindowFunc{
		Name: name,
		Fn:   fn,
	})
}

func (m *Manager) identifyWindow(winInfo *WindowInfo) (innerId string, appInfo *AppInfo) {
	logger.Debugf("identifyWindow: window id: %v, window innerId: %q",
		winInfo.window, winInfo.innerId)
	if winInfo.innerId == "" {
		logger.Debug("identifyWindow: winInfo.innerId is empty")
		return
	}

	for idx, item := range m.identifyWindowFuns {
		name := item.Name
		logger.Debugf("identifyWindow: try %s:%d", name, idx)
		innerId, appInfo = item.Fn(m, winInfo)
		if innerId != "" {
			// success
			logger.Debugf("identifyWindow by %s success, innerId: %q, appInfo: %v",
				name, innerId, appInfo)
			// NOTE: if name == "Pid", appInfo may be nil
			if appInfo != nil {
				fixedAppInfo := fixAutostartAppInfo(appInfo)
				if fixedAppInfo != nil {
					appInfo = fixedAppInfo
					appInfo.identifyMethod = name + "+FixAutostart"
					innerId = fixedAppInfo.innerId
				} else {
					appInfo.identifyMethod = name
				}
			}
			return
		}
	}
	// fail
	logger.Debugf("identifyWindow: failed")
	return winInfo.innerId, nil
}

func fixAutostartAppInfo(appInfo *AppInfo) *AppInfo {
	file := appInfo.GetFileName()
	if isInAutostartDir(file) {
		logger.Debug("file is in autostart dir")
		base := filepath.Base(file)
		return NewAppInfo(base)
	}
	return nil
}

func identifyWindowByScratch(m *Manager, winInfo *WindowInfo) (string, *AppInfo) {
	desktopFile := filepath.Join(scratchDir, addDesktopExt(winInfo.innerId))
	logger.Debugf("try scratch desktop file: %q", desktopFile)
	appInfo := NewAppInfoFromFile(desktopFile)
	if appInfo != nil {
		// success
		return appInfo.innerId, appInfo
	}
	// fail
	return "", nil
}

func identifyWindowByPid(m *Manager, winInfo *WindowInfo) (string, *AppInfo) {
	if winInfo.pid > 10 {
		logger.Debugf("identifyWindowByPid: pid: %d", winInfo.pid)
		entry := m.Entries.GetByWindowPid(winInfo.pid)
		if entry != nil {
			// success
			return entry.innerId, entry.appInfo
		}
	}
	// fail
	return "", nil
}

func identifyWindowByGtkAppId(m *Manager, winInfo *WindowInfo) (string, *AppInfo) {
	gtkAppId := winInfo.gtkAppId
	logger.Debugf("identifyWindowByGtkAppId: gtkAppId: %q", gtkAppId)
	if gtkAppId != "" {
		appInfo := NewAppInfo(gtkAppId)
		if appInfo != nil {
			// success
			return appInfo.innerId, appInfo
		}
	}
	// fail
	return "", nil
}

func identifyWindowByFlatpakAppID(m *Manager, winInfo *WindowInfo) (string, *AppInfo) {
	flatpakRef := winInfo.flatpakAppID
	logger.Debugf("identifyWindowByFlatpakAppID: flatpak ref is %q", flatpakRef)
	if strings.HasPrefix(flatpakRef, "app/") {
		parts := strings.Split(flatpakRef, "/")
		if len(parts) > 1 {
			appID := parts[1]
			appInfo := NewAppInfo(appID)
			if appInfo != nil {
				// success
				return appInfo.innerId, appInfo
			}
		}
	}
	// fail
	return "", nil
}

var crxAppIdMap = map[string]string{
	"crx_onfalgmmmaighfmjgegnamdjmhpjpgpi": "apps.com.aiqiyi",
	"crx_gfhkopakpiiaeocgofdpcpjpdiglpkjl": "apps.cn.kugou.hd",
	"crx_gaoopbnflngfkoobibfgbhobdeiipcgh": "apps.cn.kuwo.kwmusic",
	"crx_jajaphleehpmpblokgighfjneejapnok": "apps.com.evernote",
	"crx_ebhffdbfjilfhahiinoijchmlceailfn": "apps.com.letv",
	"crx_almpoflgiciaanepplakjdkiaijmklld": "apps.com.tongyong.xxbox",
	"crx_heaphplipeblmpflpdcedfllmbehonfo": "apps.com.peashooter",
	"crx_dbngidmdhcooejaggjiochbafiaefndn": "apps.com.rovio.angrybirdsseasons",
	"crx_chfeacahlaknlmjhiagghobhkollfhip": "apps.com.sina.weibo",
	"crx_cpbmecbkmjjfemjiekledmejoakfkpec": "apps.com.openapp",
	"crx_lalomppgkdieklbppclocckjpibnlpjc": "apps.com.baidutieba",
	"crx_gejbkhjjmicgnhcdpgpggboldigfhgli": "apps.com.zhuishushenqi",
	"crx_gglenfcpioacendmikabbkecnfpanegk": "apps.com.duokan",
	"crx_nkmmgdfgabhefacpfdabadjfnpffhpio": "apps.com.zhihu.daily",
	"crx_ajkogonhhcighbinfgcgnjiadodpdicb": "apps.com.netease.newsreader",
	"crx_hgggjnaaklhemplabjhgpodlcnndhppo": "apps.com.baidu.music.pad",
	"crx_ebmgfebnlgilhandilnbmgadajhkkmob": "apps.cn.ibuka",
	"crx_nolebplcbgieabkblgiaacdpgehlopag": "apps.com.tianqitong",
	"crx_maghncnmccfbmkekccpmkjjfcmdnnlip": "apps.com.youjoy.strugglelandlord",
	"crx_heliimhfjgfabpgfecgdhackhelmocic": "apps.cn.emoney",
	"crx_jkgmneeafmgjillhgmjbaipnakfiidpm": "apps.com.instagram",
	"crx_cdbkhmfmikobpndfhiphdbkjklbmnakg": "apps.com.easymindmap",
	"crx_djflcciklfljleibeinjmjdnmenkciab": "apps.com.lgj.thunderbattle",
	"crx_ffdgbolnndgeflkapnmoefhjhkeilfff": "apps.com.qianlong",
	"crx_fmpniepgiofckbfgahajldgoelogdoap": "apps.com.windhd",
	"crx_dokjmallmkihbgefmladclcdcinjlnpj": "apps.com.youdao.hanyu",
	"crx_dicimeimfmbfcklbjdpnlmjgegcfilhm": "apps.com.ibookstar",
	"crx_cokkcjnpjfffianjbpjbcgjefningkjm": "apps.com.yidianzixun",
	"crx_ehflkacdpmeehailmcknlnkmjalehdah": "apps.com.xplane",
	"crx_iedokjbbjejfinokgifgecmboncmkbhb": "apps.com.wedevote",
	"crx_eaefcagiihjpndconigdpdmcbpcamaok": "apps.com.tongwei.blockbreaker",
	"crx_mkjjfibpccammnliaalefmlekiiikikj": "apps.com.dayima",
	"crx_gflkpppiigdigkemnjlonilmglokliol": "apps.com.cookpad",
	"crx_jfhpkchgedddadekfeganigbenbfaohe": "apps.com.issuu",
	"crx_ggkmfnbkldhmkehabgcbnmlccfbnoldo": "apps.bible.cbol",
	"crx_phlhkholfcljapmcidanddmhpcphlfng": "apps.com.kanjian.radio",
	"crx_bjgfcighhaahojkibojkdmpdihhcehfm": "apps.de.danoeh.antennapod",
	"crx_kldipknjommdfkifomkmcpbcnpmcnbfi": "apps.com.asoftmurmur",
	"crx_jfhlegimcipljdcionjbipealofoncmd": "apps.com.tencentnews",
	"crx_aikgmfkpmmclmpooohngmcdimgcocoaj": "apps.com.tonghuashun",
	"crx_ifimglalpdeoaffjmmihoofapmpflkad": "apps.com.letv.lecloud.disk",
	"crx_pllcekmbablpiogkinogefpdjkmgicbp": "apps.com.hwadzanebook",
	"crx_ohcknkkbjmgdfcejpbmhjbohnepcagkc": "apps.com.douban.radio",
}

func identifyWindowByCrxId(m *Manager, winInfo *WindowInfo) (string, *AppInfo) {
	if winInfo.wmClass != nil &&
		strings.EqualFold(winInfo.wmClass.Class, "chromium-browser") &&
		strings.HasPrefix(winInfo.wmClass.Instance, "crx_") {

		appId, ok := crxAppIdMap[winInfo.wmClass.Instance]
		logger.Debug("identifyWindowByCrxId", appId)
		if ok {
			appInfo := NewAppInfo(appId)
			if appInfo != nil {
				// success
				return appInfo.innerId, appInfo
			}
		}
	}
	// fail
	return "", nil
}

func identifyWindowByCmdlineTurboBooster(m *Manager, winInfo *WindowInfo) (string, *AppInfo) {
	pid := winInfo.pid
	process := winInfo.process
	if process != nil && pid != 0 {
		if len(process.cmdline) >= 0 {
			var desktopFile string
			if strings.HasSuffix(process.cmdline[0], desktopExt) {
				desktopFile = process.cmdline[0]
			} else if strings.Contains(process.cmdline[0], "/applications/") {
				matches, err := filepath.Glob(process.cmdline[0] + "*")
				if err != nil {
					logger.Warning("filepath.Glob err:", err)
					return "", nil
				}
				if len(matches) > 0 && strings.HasSuffix(matches[0], desktopExt) {
					desktopFile = matches[0]
				}
			}

			if desktopFile != "" {
				logger.Debugf("identifyWindowByCmdlineTurboBooster: desktopFile: %s", desktopFile)
				appInfo := NewAppInfoFromFile(desktopFile)
				if appInfo != nil {
					// success
					return appInfo.innerId, appInfo
				}
			}
		}
	}

	// fail
	return "", nil
}

func identifyWindowByPidEnv(m *Manager, winInfo *WindowInfo) (string, *AppInfo) {
	pid := winInfo.pid
	process := winInfo.process
	if process != nil && pid != 0 {
		launchedDesktopFile := process.environ.Get("GIO_LAUNCHED_DESKTOP_FILE")
		launchedDesktopFilePid, _ := strconv.ParseUint(
			process.environ.Get("GIO_LAUNCHED_DESKTOP_FILE_PID"), 10, 32)

		logger.Debugf("identifyWindowByPidEnv: launchedDesktopFile: %q, pid: %d",
			launchedDesktopFile, launchedDesktopFilePid)

		var try bool
		if uint(launchedDesktopFilePid) == pid {
			try = true
		} else if uint(launchedDesktopFilePid) == process.ppid && process.ppid != 0 {
			logger.Debug("ppid equal")
			parentProcess := procfs.Process(process.ppid)
			cmdline, err := parentProcess.Cmdline()
			if err == nil && len(cmdline) > 0 {
				logger.Debugf("parent process cmdline: %#v", cmdline)
				base := filepath.Base(cmdline[0])
				if base == "sh" || base == "bash" {
					try = true
				}
			}
		}

		if try {
			appInfo := NewAppInfoFromFile(launchedDesktopFile)
			if appInfo != nil {
				// success
				return appInfo.innerId, appInfo
			}
		}
	}
	// fail
	return "", nil
}

func identifyWindowByRule(m *Manager, winInfo *WindowInfo) (string, *AppInfo) {
	ret := m.windowPatterns.Match(winInfo)
	if ret == "" {
		return "", nil
	}
	logger.Debug("identifyWindowByRule ret:", ret)
	// parse ret
	// id=$appId or env
	var appInfo *AppInfo
	if len(ret) > 4 && strings.HasPrefix(ret, "id=") {
		appInfo = NewAppInfo(ret[3:])
	} else if ret == "env" {
		process := winInfo.process
		if process != nil {
			launchedDesktopFile := process.environ.Get("GIO_LAUNCHED_DESKTOP_FILE")
			if launchedDesktopFile != "" {
				appInfo = NewAppInfoFromFile(launchedDesktopFile)
			}
		}
	} else {
		logger.Warningf("bad ret: %q", ret)
	}

	if appInfo != nil {
		return appInfo.innerId, appInfo
	}
	return "", nil
}

func identifyWindowByWmClass(m *Manager, winInfo *WindowInfo) (string, *AppInfo) {
	if winInfo.wmClass != nil {
		instance := winInfo.wmClass.Instance
		if instance != "" {
			// example:
			// WM_CLASS(STRING) = "Brackets", "Brackets"
			// wm class instance is Brackets
			// try app id org.deepin.flatdeb.brackets
			appInfo := NewAppInfo("org.deepin.flatdeb." + strings.ToLower(instance))
			if appInfo != nil {
				return appInfo.innerId, appInfo
			}

			appInfo = NewAppInfo(instance)
			if appInfo != nil {
				return appInfo.innerId, appInfo
			}
		}

		class := winInfo.wmClass.Class
		if class != "" {
			appInfo := NewAppInfo(class)
			if appInfo != nil {
				return appInfo.innerId, appInfo
			}
		}
	}
	// fail
	return "", nil
}

func identifyWindowByBamf(m *Manager, winInfo *WindowInfo) (string, *AppInfo) {
	// bamf
	win := winInfo.window
	desktop, err := getDesktopFromWindowByBamf(win)
	if err != nil {
		logger.Warning(err)
		return "", nil
	}

	if desktop != "" {
		appInfo := NewAppInfoFromFile(desktop)
		if appInfo != nil {
			// success
			return appInfo.innerId, appInfo
		}
	}
	return "", nil
}

func identifyWindowByCmdlineXWalk(m *Manager, winInfo *WindowInfo) (string, *AppInfo) {
	process := winInfo.process
	if process == nil || winInfo.pid == 0 {
		return "", nil
	}

	exeBase := filepath.Base(process.exe)
	args := process.args
	if exeBase != "xwalk" || len(args) == 0 {
		return "", nil
	}
	lastArg := args[len(args)-1]
	logger.Debugf("lastArg: %q", lastArg)

	if filepath.Base(lastArg) == "manifest.json" {
		appId := filepath.Base(filepath.Dir(lastArg))
		appInfo := NewAppInfo(appId)
		if appInfo != nil {
			// success
			return appInfo.innerId, appInfo
		}
	}
	// failed
	return "", nil
}
