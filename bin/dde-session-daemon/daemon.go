package main

import "dde-daemon/network"
import "dde-daemon/clipboard"
import "dde-daemon/audio"
import "dde-daemon/power"

import "dde-daemon/display"
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
import "time"
import . "dlib/gettext"
import "dlib"
import "dlib/logger"
import "dlib/utils"
import "path"
import "os/exec"
import "os"
import "dlib/dbus"

import _ "net/http/pprof"
import "net/http"

func init() {
	go http.ListenAndServe("localhost:6060", nil)
}

var Logger = logger.NewLogger("com.deepin.daemon")
var objUtil = utils.NewUtils()

const (
	DSC_CONFIG_PATH = ".config/deepin-software-center/config_info.ini"
)

func setDSCAutoUpdate(interval time.Duration) {
	if interval <= 0 {
		return
	}

	for {
		timer := time.After(time.Hour * interval)
		select {
		case <-timer:
			go exec.Command("/usr/bin/dsc-daemon", []string{"--no-daemon"}...).Run()
		}
	}
}

func dscAutoUpdate() {
	homeDir, ok := objUtil.GetHomeDir()
	if !ok {
		return
	}
	filename := path.Join(homeDir, DSC_CONFIG_PATH)
	if !objUtil.IsFileExist(filename) {
		return
	}

	interval, ok1 := objUtil.ReadKeyFromKeyFile(filename,
		"update", "interval", int32(0))
	if !ok1 {
		interval = 3
	}
	isUpdate, ok2 := objUtil.ReadKeyFromKeyFile(filename,
		"update", "auto", false)
	if !ok2 {
		isUpdate = true
	}
	if v, ok := isUpdate.(bool); ok && v {
		if i, ok := interval.(int32); ok {
			go setDSCAutoUpdate(time.Duration(i))
		}
	}
}

func main() {
	if !dlib.UniqueOnSession("com.deepin.daemon") {
		Logger.Warning("There already has an dde-daemon running.")
		return
	}
	InitI18n()
	Textdomain("dde-daemon")

	C.init()

	netFlag := false
	clipFlag := false
	audioFlag := false
	powerFlag := false
	dpyFlag := false
	dockFlag := false
	launFlag := false
	keyFlag := false
	mtFlag := false
	dtFlag := false
	mimeFlag := false
	themeFlag := false
	blueFlag := false
	edgeFlag := false
	mpFlag := false
	grubFlag := false
	l := len(os.Args)
	if l >= 2 {
		for i := 1; i < l; i++ {
			switch os.Args[i] {
			case "network":
				netFlag = true
			case "clipboard":
				clipFlag = true
			case "audio":
				audioFlag = true
			case "power":
				powerFlag = true
			case "display":
				dpyFlag = true
			case "dock":
				dockFlag = true
			case "launcher":
				launFlag = true
			case "keybinding":
				keyFlag = true
			case "mounts":
				mtFlag = true
			case "datetime":
				dtFlag = true
			case "mime":
				mimeFlag = true
			case "themes":
				themeFlag = true
			case "bluetooth":
				blueFlag = true
			case "screen_edges":
				edgeFlag = true
			case "mpris":
				mpFlag = true
			case "grub2":
				grubFlag = true
			}
		}
	}

	if !netFlag {
		go network.Start()
	}
	if !clipFlag {
		go clipboard.Start()
	}

	if !audioFlag {
		go audio.Start()
	}
	if !powerFlag {
		go power.Start()
	}
	if !dpyFlag {
		go display.Start()
	}
	<-time.After(time.Second * 3)

	if !dockFlag || !dpyFlag {
		go dock.Start()
	}
	if !launFlag {
		go launcher.Start()
	}

	if !keyFlag {
		go keybinding.Start()
	}
	if !dtFlag {
		go datetime.Start()
	}
	if !mimeFlag {
		go mime.Start()
	}
	if !mtFlag {
		go mounts.Start()
	}
	if !themeFlag {
		go themes.Start()
	}
	if !blueFlag {
		go bluetooth.Start()
	}

	if !mpFlag {
		go startMprisDaemon()
	}

	dscAutoUpdate()

	<-time.After(time.Second)
	if !dpyFlag || !edgeFlag {
		go screen_edges.Start()
	}
	if !grubFlag {
		go grub2.Start()
	}
	glib.StartLoop()

	if err := dbus.Wait(); err != nil {
		Logger.Errorf("Lost dbus: %v", err)
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
