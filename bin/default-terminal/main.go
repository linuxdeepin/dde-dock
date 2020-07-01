package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/godbus/dbus"
	sessionmanager "github.com/linuxdeepin/go-dbus-factory/com.deepin.sessionmanager"
	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/appinfo/desktopappinfo"
)

const (
	gsSchemaDefaultTerminal = "com.deepin.desktop.default-applications.terminal"
	gsKeyAppId              = "app-id"
	gsKeyExec               = "exec"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	settings := gio.NewSettings(gsSchemaDefaultTerminal)
	defer settings.Unref()

	appId := settings.GetString(gsKeyAppId)
	appInfo := desktopappinfo.NewDesktopAppInfo(appId)

	if len(os.Args) == 1 {
		if appInfo != nil {
			sessionBus, err := dbus.SessionBus()
			if err != nil {
				log.Fatal(err)
			}
			startManager := sessionmanager.NewStartManager(sessionBus)
			filename := appInfo.GetFileName()
			workDir, err := os.Getwd()
			if err != nil {
				log.Println("warning: failed to get work dir:", err)
			}
			options := map[string]dbus.Variant{
				"path": dbus.MakeVariant(workDir),
			}
			err = startManager.LaunchAppWithOptions(0, filename, 0,
				nil, options)
			if err != nil {
				log.Println(err)
				runFallbackTerm()
			}
		} else {
			runFallbackTerm()
		}
	} else {
		// define -e option
		termExec := settings.GetString(gsKeyExec)
		termPath, _ := exec.LookPath(termExec)
		if termPath == "" {
			// try again
			termPath = getTerminalPath()
			if termPath == "" {
				log.Fatal("failed to get terminal path")
			}
		}

		args := os.Args[1:]
		cmd := exec.Command(termPath, args...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func runFallbackTerm() {
	termPath := getTerminalPath()
	if termPath == "" {
		log.Println("failed to get terminal path")
		return
	}
	cmd := exec.Command(termPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

var terms = []string{
	"deepin-terminal",
	"gnome-terminal",
	"terminator",
	"xfce4-terminal",
	"rxvt",
	"xterm",
}

func getTerminalPath() string {
	for _, exe := range terms {
		file, _ := exec.LookPath(exe)
		if file != "" {
			return file
		}
	}
	return ""
}
