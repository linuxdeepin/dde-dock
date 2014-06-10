package main

import "dde-daemon/screensaver"
import "dde-daemon/network"
import "dde-daemon/clipboard"
import "dde-daemon/audio"
import "dde-daemon/power"

import "dde-daemon/keybinding"
import "dde-daemon/datetime"
import "dde-daemon/mime"

import "dde-daemon/mounts"
import "dde-daemon/bluetooth"

import "dde-daemon/screen_edges"
import "dde-daemon/themes"

import "dde-daemon/dock"
import "dde-daemon/launcher"
import "dde-daemon/grub2"

import "dlib/glib-2.0"

//#cgo pkg-config:gtk+-3.0
//#include <gtk/gtk.h>
//void init(){gtk_init(0,0);}
import "C"
import . "dlib/gettext"
import "dlib"
import "dlib/logger"
import "os"
import "dlib/dbus"

type Module struct {
	Name   string
	Start  func()
	Stop   func()
	Enable bool
}

var modules = []Module{
	Module{"keybinding", keybinding.Start, nil, true},
	Module{"screensaver", screensaver.Start, nil, true},
	Module{"power", power.Start, nil, true},
	Module{"themes", themes.Start, nil, true},
	Module{"screen_edges", screen_edges.Start, nil, true},
	Module{"launcher", launcher.Start, nil, true},
	Module{"dock", dock.Start, nil, true},

	Module{"network", network.Start, nil, true},
	Module{"audio", audio.Start, nil, true},
	Module{"mounts", mounts.Start, nil, true},
	Module{"datetime", datetime.Start, nil, true},
	Module{"mime", mime.Start, nil, true},
	Module{"bluetooth", bluetooth.Start, nil, true},
	Module{"grub2", grub2.Start, nil, true},
	Module{"dsc", dscAutoUpdate, nil, true},
	Module{"mpris", startMprisDaemon, nil, true},
	Module{"clipboard", clipboard.Start, nil, true},
}

var Logger = logger.NewLogger("com.deepin.daemon")

func main() {
	if !dlib.UniqueOnSession("com.deepin.daemon") {
		Logger.Warning("There already has an dde-daemon running.")
		return
	}

	Logger.BeginTracing()
	defer Logger.EndTracing()
	Logger.SetRestartCommand("/usr/lib/deepin-daemon/dde-session-daemon")

	InitI18n()
	Textdomain("dde-daemon")

	C.init()

	if len(os.Args) >= 2 {
		for _, disabledModuleName := range os.Args[1:] {
			for _, m := range modules {
				if disabledModuleName == m.Name {
					m.Enable = false
					break
				}
			}
		}
	}

	for _, module := range modules {
		if module.Enable {
			go module.Start()
		}
	}

	go func() {
		if err := dbus.Wait(); err != nil {
			Logger.Errorf("Lost dbus: %v", err)
			os.Exit(-1)
		} else {
			os.Exit(0)
		}
	}()

	ddeSessionRegister()
	glib.StartLoop()
}
