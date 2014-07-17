package power

//this should only use org.freedesktop.ScreenSaver interface with SimulateUserActivity

import (
	"dbus/org/freedesktop/screensaver"
	//"pkg.linuxdeepin.com/lib/logger"
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
	"io/ioutil"
	"strings"
	"time"
)

type fullScreenWorkaround struct {
	xu               *xgbutil.XUtil
	targets          []string
	activeWindowAtom xproto.Atom
	isHintingTarget  bool
}

func newFullScreenWorkaround() *fullScreenWorkaround {
	XU, _ := xgbutil.NewConn()
	ACTIVE_WINDOW, _ := xprop.Atm(XU, "_NET_ACTIVE_WINDOW")
	w := &fullScreenWorkaround{
		XU, []string{"libflash", "chrome", "mplayer", "operaplugin", "soffice", "wpp", "evince"}, ACTIVE_WINDOW, false,
	}
	return w
}

func (wa *fullScreenWorkaround) detectTarget(w xproto.Window) {
	pid, _ := xprop.PropValNum(xprop.GetProperty(wa.xu, w, "_NET_WM_PID"))

	contents, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return
	}

	if wa.isFullScreen(w) {
		for _, target := range wa.targets {
			if strings.Contains(string(contents), target) {
				wa.inhibit(target, string(contents))
				return
			}
		}
	}
	wa.isHintingTarget = false
}

func (wa *fullScreenWorkaround) inhibit(target, cmdline string) {
	wa.isHintingTarget = true
	if ss, err := screensaver.NewScreenSaver("org.freedesktop.ScreenSaver", "/org/freedesktop/ScreenSaver"); err == nil {
		var hit func()
		hit = func() {
			time.AfterFunc(time.Second*2, func() {
				if wa.isHintingTarget {
					ss.SimulateUserActivity()
					hit()
				}
			})
		}
		hit()
		Logger.Debug("Inhibit Hight Performance :", "TARGET:", target, "CMDLINE:", cmdline)
	} else {
		Logger.Error("ERRRR:", err)
	}
}

func (wa *fullScreenWorkaround) isFullScreen(xid xproto.Window) bool {
	states, _ := ewmh.WmStateGet(wa.xu, xid)
	found := 0
	for _, s := range states {
		if s == "_NET_WM_STATE_FULLSCREEN" {
			found++
		}
		if s == "_NET_WM_STATE_FOCUSED" {
			found++
		}
	}
	if found == 2 {
		Logger.Debug("HAHAH:::::", states)
	}
	return found == 2
}

func (wa *fullScreenWorkaround) start() {
	var runner func()
	runner = func() {
		w, _ := ewmh.ActiveWindowGet(wa.xu)
		wa.detectTarget(w)
		time.AfterFunc(time.Second*5, runner)
	}
	runner()

	root := xwindow.New(wa.xu, wa.xu.RootWin())
	root.Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		if wa.activeWindowAtom == ev.Atom {
			w, _ := ewmh.ActiveWindowGet(XU)
			wa.detectTarget(w)
		}
	}).Connect(wa.xu, root.Id)
	xevent.Main(wa.xu)
}
