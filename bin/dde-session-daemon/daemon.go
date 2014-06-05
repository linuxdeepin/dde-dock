package main

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

import _ "net/http/pprof"
import "net/http"

//import "time"

type methodInfo struct {
	f      func()
	filter bool
}

var methodFilterMap = map[string]*methodInfo{
	"network":      &methodInfo{network.Start, false},
	"clipboard":    &methodInfo{clipboard.Start, false},
	"audio":        &methodInfo{audio.Start, false},
	"power":        &methodInfo{power.Start, false},
	"dock":         &methodInfo{dock.Start, false},
	"launcher":     &methodInfo{launcher.Start, false},
	"keybinding":   &methodInfo{keybinding.Start, false},
	"mounts":       &methodInfo{mounts.Start, false},
	"datetime":     &methodInfo{datetime.Start, false},
	"mime":         &methodInfo{mime.Start, false},
	"themes":       &methodInfo{themes.Start, false},
	"bluetooth":    &methodInfo{bluetooth.Start, false},
	"screen_edges": &methodInfo{screen_edges.Start, false},
	"grub2":        &methodInfo{grub2.Start, false},
	"dsc":          &methodInfo{dscAutoUpdate, false},
	"mpris":        &methodInfo{startMprisDaemon, false},
}

var Logger = logger.NewLogger("com.deepin.daemon")

func init() {
	go http.ListenAndServe("localhost:6060", nil)
}

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

	l := len(os.Args)
	if l >= 2 {
		for i := 1; i < l; i++ {
			if v, ok := methodFilterMap[os.Args[i]]; ok {
				v.filter = true
			}
		}
	}

	for _, v := range methodFilterMap {
		if !v.filter {
			go v.f()
			//<-time.After(time.Millisecond * 300)
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

	glib.StartLoop()
}
