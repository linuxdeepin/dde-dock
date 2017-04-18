package main

import (
	"flag"
	"fmt"
	"os"
	"pkg.deepin.io/dde/daemon/grub2"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"time"
)

var logger = log.NewLogger("daemon/grub2")

func init() {
	grub2.SetLogger(logger)
}

var (
	argDebug      bool
	argSetup      bool
	argSetupTheme bool
	//argGrubSettingFile string
	//argThemeDir        string
	argResolution string
)

func main() {
	flag.BoolVar(&argDebug, "d", false, "debug mode")
	flag.BoolVar(&argDebug, "debug", false, "debug mode")
	flag.BoolVar(&argSetup, "setup", false, "setup grub and exit")
	flag.BoolVar(&argSetupTheme, "setup-theme", false, "setup grub theme only and exit")
	//flag.StringVar(&argGrubSettingFile, "setting-file", "", "specify an alternative setting file instead of /etc/default/grub when setup grub")
	// TODO --grub-config, --backend, [grub, efi]
	//flag.StringVar(&argThemeDir, "theme-dir", "", "specify an alternative theme directory instead of /boot/grub/themes/deepin when setup grub")
	flag.StringVar(&argResolution, "gfxmode", "auto", "specify gfxmode when setup grub")
	flag.Parse()
	if argDebug {
		logger.SetLogLevel(log.LevelDebug)
	}

	if argSetupTheme {
		fmt.Println("mode: setup theme")
		err := grub2.SetupTheme()
		if err != nil {
			logger.Warning(err)
			os.Exit(1)
		}
	} else if argSetup {
		fmt.Println("mode: setup")
		err := grub2.Setup(argResolution)
		if err != nil {
			logger.Warning(err)
			os.Exit(2)
		}
	} else {
		fmt.Println("mode: daemon")
		runAsDaemon()
	}
}

func runAsDaemon() {
	if !lib.UniqueOnSystem(grub2.DBusDest) {
		logger.Warning("There already has an Grub2 daemon running.")
		return
	}

	err := grub2.Start()
	if err != nil {
		logger.Warning("start failed:", err)
		os.Exit(3)
	}

	dbus.SetAutoDestroyHandler(5*time.Minute, grub2.CanSafelyExit)
	dbus.DealWithUnhandledMessage()

	if err := dbus.Wait(); err != nil {
		logger.Error("lost dbus session:", err)
		os.Exit(4)
	}
}
