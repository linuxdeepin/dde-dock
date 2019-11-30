/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package main

//#cgo pkg-config: x11 gtk+-3.0
//#include <X11/Xlib.h>
//#include <gtk/gtk.h>
//void init(){XInitThreads();gtk_init(0,0);}
import "C"
import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/linuxdeepin/go-dbus-factory/com.deepin.api.soundthemeplayer"
	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/dde/api/userenv"
	"pkg.deepin.io/dde/daemon/loader"
	dapp "pkg.deepin.io/lib/app"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/proxy"
	"pkg.deepin.io/lib/utils"
	"pkg.deepin.io/lib/xdg/basedir"
)

var logger = log.NewLogger("daemon/dde-session-daemon")
var hasDDECookie bool

func allowRun() bool {
	if os.Getenv("DDE_SESSION_PROCESS_COOKIE_ID") != "" {
		hasDDECookie = true
		return true
	}

	systemBus, err := dbus.SessionBus()
	if err != nil {
		logger.Warning(err)
		os.Exit(1)
	}
	sessionManagerObj := systemBus.Object("com.deepin.SessionManager",
		"/com/deepin/SessionManager")
	var allowRun bool
	err = sessionManagerObj.Call("com.deepin.SessionManager.AllowSessionDaemonRun",
		dbus.FlagNoAutoStart).Store(&allowRun)
	if err != nil {
		logger.Warning(err)
		return true
	}

	return allowRun
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	if !allowRun() {
		logger.Warning("session manager does not allow me to run")
		os.Exit(1)
	}

	InitI18n()
	Textdomain("dde-daemon")

	cmd := dapp.New("dde-session-dameon", "dde session daemon", "version "+__VERSION__)

	usr, err := user.Current()
	if err == nil {
		os.Chdir(usr.HomeDir)
	}

	flags := new(Flags)
	flags.IgnoreMissingModules = cmd.Flag("Ignore", "ignore missing modules, --no-ignore to revert it.").Short('i').Default("true").Bool()
	flags.ForceStart = cmd.Flag("force", "Force start disabled module.").Short('f').Bool()

	cmd.Command("auto", "Automatically get enabled and disabled modules from settings.").Default()
	enablingModules := cmd.Command("enable", "Enable modules and their dependencies, ignore settings.").Arg("module", "module names.").Required().Strings()
	disableModules := cmd.Command("disable", "Disable modules, ignore settings.").Arg("module", "module names.").Required().Strings()
	listModule := cmd.Command("list", "List all the modules or the dependencies of one module.").Arg("module", "module name.").String()

	subCmd := cmd.ParseCommandLine(os.Args[1:])

	cmd.StartProfile()

	C.init()
	proxy.SetupProxy()

	app := NewSessionDaemon(flags, daemonSettings, logger)

	service, err := dbusutil.NewSessionService()
	if err != nil {
		logger.Fatal(err)
	}

	if err = app.register(service); err != nil {
		logger.Info(err)
		os.Exit(0)
	}

	loader.SetService(service)

	logger.SetLogLevel(log.LevelInfo)
	if cmd.IsLogLevelNone() &&
		(utils.IsEnvExists(log.DebugLevelEnv) || utils.IsEnvExists(log.DebugMatchEnv)) {
		logger.Info("Log level is none and debug env exists, so ignore cmd.loglevel")
	} else {
		appLogLevel := cmd.LogLevel()
		logger.Info("App log level:", appLogLevel)
		// set all modules log level to appLogLevel
		loader.SetLogLevel(appLogLevel)
	}

	needRunMainLoop := true

	// Ensure each module and mainloop in the same thread
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	switch subCmd {
	case "auto":
		app.execDefaultAction()
	case "enable":
		err = app.enableModules(*enablingModules)
	case "disable":
		err = app.disableModules(*disableModules)
	case "list":
		err = app.listModule(*listModule)
		needRunMainLoop = false
	}

	if err != nil {
		logger.Warning(err)
		os.Exit(1)
	}

	err = migrateUserEnv()
	if err != nil {
		logger.Warning("failed to migrate user env:", err)
	}

	err = syncConfigToSoundThemePlayer()
	if err != nil {
		logger.Warning(err)
	}

	if needRunMainLoop {
		runMainLoop()
	}
}

// migrate user env from ~/.pam_environment to ~/.dde-env
func migrateUserEnv() error {
	_, err := os.Stat(userenv.DefaultFile())
	if os.IsNotExist(err) {
		// when ~/.dde-env does not exist, read ~/.pam_environment,
		// remove the key we set before, and save it back.
		pamEnvFile := filepath.Join(basedir.GetUserHomeDir(), ".pam_environment")
		pamEnv, err := loadPamEnv(pamEnvFile)
		if os.IsNotExist(err) {
			return nil
		} else if err != nil {
			return err
		}

		var reservedPamEnv []pamEnvKeyValue
		for _, kv := range pamEnv {
			switch kv.key {
			case "LANG", "LANGUAGE", "QT_SCALE_FACTOR", "_JAVA_OPTIONS":
			// ignore it
			default:
				reservedPamEnv = append(reservedPamEnv, kv)
			}
		}

		if len(reservedPamEnv) == 0 {
			err = os.Remove(pamEnvFile)
		} else {
			err = savePamEnv(pamEnvFile, reservedPamEnv)
		}
		if err != nil {
			return err
		}

		// save the current env to ~/.dde-env
		currentEnv := make(map[string]string)
		for _, envKey := range []string{"LANG", "LANGUAGE"} {
			envValue, ok := os.LookupEnv(envKey)
			if ok {
				currentEnv[envKey] = envValue
			}
		}
		err = userenv.Save(currentEnv)
		return err
	} else if err != nil {
		return err
	}
	return nil
}

type pamEnvKeyValue struct {
	key, value string
}

func loadPamEnv(filename string) ([]pamEnvKeyValue, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var result []pamEnvKeyValue
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}

		fields := strings.SplitN(line, "=", 2)
		if len(fields) == 2 {
			result = append(result, pamEnvKeyValue{
				key:   fields[0],
				value: fields[1],
			})
		}
	}
	return result, nil
}

func savePamEnv(filename string, pamEnv []pamEnvKeyValue) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	bw := bufio.NewWriterSize(f, 256)
	for _, kv := range pamEnv {
		_, err = fmt.Fprintf(bw, "%s=%s\n", kv.key, kv.value)
		if err != nil {
			return err
		}
	}

	err = bw.Flush()
	return err
}

func syncConfigToSoundThemePlayer() error {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	loginEnabled := soundutils.CanPlayEvent(soundutils.EventDesktopLogin)
	player := soundthemeplayer.NewSoundThemePlayer(sysBus)
	err = player.EnableSoundDesktopLogin(0, loginEnabled)
	return err
}
