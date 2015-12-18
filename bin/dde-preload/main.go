package main

//#cgo pkg-config: x11 gtk+-3.0
//#include <X11/Xlib.h>
//#include <gtk/gtk.h>
//void init(){XInitThreads();gtk_init(0,0);}
import "C"
import (
	"gir/glib-2.0"
	_ "pkg.deepin.io/dde/daemon/dock"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/proxy"

	"os"
	"time"
)

var logger = log.NewLogger("daemon/preload")

func runMainLoop() {
	logger.Info("register session")
	startTime := time.Now()
	ddeSessionRegister()
	logger.Info("register session done, cost", time.Now().Sub(startTime))

	logger.Info("DealWithUnhandledMessage")
	startTime = time.Now()
	dbus.DealWithUnhandledMessage()
	logger.Info("DealWithUnhandledMessage done, cost", time.Now().Sub(startTime))
	go glib.StartLoop()

	logger.Info("initialize done")
	if err := dbus.Wait(); err != nil {
		logger.Errorf("Lost dbus: %v", err)
		os.Exit(-1)
	}

	logger.Info("dbus connection is closed by user")
	os.Exit(0)
}

type Preload struct {
}

func (*Preload) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Preload",
		ObjectPath: "/com/deepin/daemon/Preload",
		Interface:  "com.deepin.daemon.Preload",
	}
}

func main() {
	preload := new(Preload)
	if !lib.UniqueOnSession(preload.GetDBusInfo().Dest) {
		logger.Warning("There already has a dde preload running.")
		os.Exit(0)
	}

	err := dbus.InstallOnSession(preload)
	if err != nil {
		logger.Fatal(err)
	}

	cmd := InitCMD()
	cmd.Parse(os.Args[1:])

	cfg := &Config{
		CPUProfile: cmd.CpuProf(),
		MemProfile: cmd.MemProf(),
	}
	cfg.Start()
	defer cfg.Stop()

	InitI18n()
	Textdomain("dde-daemon")

	C.init()
	proxy.SetupProxy()

	loader.SetLogLevel(cmd.LogLevel())
	loader.EnableModules([]string{"dock"}, nil, loader.EnableFlagIgnoreMissingModule)

	runMainLoop()
}
