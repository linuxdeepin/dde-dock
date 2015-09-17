package main

//#cgo pkg-config:gtk+-3.0
//#include <gtk/gtk.h>
//void init(){gtk_init(0,0);}
import "C"
import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
	"pkg.deepin.io/dde/daemon/loader"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/proxy"
	"runtime/pprof"
)

// using go build -ldflags "-X main.__VERSION__ version" to set version.
var __VERSION__ = "unknown"

var logger = log.NewLogger("daemon/dde-session-daemon")

func startMemProfile(name string) {
	logger.Info("Start memory profile")
	f, err := os.Create(name)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			switch sig.String() {
			case "Interrupt":
				logger.Info("Memory profile done.")
				pprof.WriteHeapProfile(f)
				f.Close()
				close(c)
				os.Exit(0)
			}
		}
	}()
}

func main() {
	InitI18n()
	Textdomain("dde-daemon")

	cmd := kingpin.New("dde-session-daemon", "session daemon")
	cmd.Version("version " + __VERSION__)

	flags := new(Flags)
	flags.Verbose = cmd.Flag("verbose", "show much more message, the shorthand for --loglevel debug, if specificed, loglevel is ignored.").Short('v').Bool()
	flags.LogLevel = cmd.Flag("loglevel", "set log level, possible value is error/warn/info/debug/no.").Short('l').String()
	flags.IgnoreMissingModules = cmd.Flag("ignore", "ignore missing modules, --no-ignore to revert it.").Short('i').Default("true").Bool()
	flags.ForceStart = cmd.Flag("force", "force start disabled module.").Short('f').Bool()

	enablingModules := cmd.Command("enable", "enable modules and their dependencies, ignore settings.").Arg("module", "module names.").Required().Strings()
	disableModules := cmd.Command("disable", "disable modules, ignore settings.").Arg("module", "module names.").Required().Strings()
	listModule := cmd.Command("list", "list all the modules or the dependencies of one module.").Arg("module", "module name.").String()

	memprof := cmd.Flag("memprof", "memory profile").String()

	app := NewSessionDaemon(cmd, flags, daemonSettings, logger)

	subCmd, err := app.parse()
	if err != nil {
		fmt.Println(err)
		fmt.Println()
		app.printUsage()
		return
	}

	C.init()
	proxy.SetupProxy()

	if subCmd == "" {
		app.execDefaultAction()
		return
	}

	app.exitIfNotSingleton()

	if *flags.Verbose {
		app.logLevel = log.LevelDebug
	} else {
		logLevel, err := toLogLevel(*flags.LogLevel)
		if err != nil {
			fmt.Println(err)
			return
		}
		app.logLevel = logLevel
	}

	loader.SetLogLevel(app.logLevel)
	logger.Info("LogLevel is", app.logLevel)

	needRunMainLoop := true
	switch subCmd {
	case "enable":
		err = app.EnableModules(*enablingModules)
	case "disable":
		err = app.DisableModules(*disableModules)
	case "list":
		err = app.ListModule(*listModule)
		needRunMainLoop = false
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if needRunMainLoop {
		if *memprof != "" {
			startMemProfile(*memprof)
		}

		runMainLoop()
	}
}
