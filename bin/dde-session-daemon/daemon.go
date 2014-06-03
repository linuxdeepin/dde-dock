package main

import "dde-daemon/clipboard"
import "dde-daemon/audio"
import "dde-daemon/power"
import "dde-daemon/display"
import "dde-daemon/keybinding"
import "dde-daemon/datetime"
import "dde-daemon/mime"
import "dde-daemon/mounts"

//import "dde-daemon/screen_edges"
import "dde-daemon/themes"

//import "dde-daemon/dock"
//import "dde-daemon/launcher"

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

	interval, _ := objUtil.ReadKeyFromKeyFile(filename,
		"update", "interval", int32(0))
	isUpdate, _ := objUtil.ReadKeyFromKeyFile(filename,
		"update", "auto", false)
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
	go clipboard.Start()
	go audio.Start()
	go power.Start()
	go display.Start()
	<-time.After(time.Second * 3)

	//go dock.Start()
	//go launcher.Start()

	go keybinding.Start()
	go datetime.Start()
	go mime.Start()
	go mounts.Start()
	go themes.Start()

	go startMprisDaemon()

	dscAutoUpdate()

	<-time.After(time.Second)
	//go screen_edges.Start()
	glib.StartLoop()
}
