package main

//#cgo pkg-config:gtk+-3.0
//#include <gtk/gtk.h>
//void init(){gtk_init(0,0);}
import "C"
import (
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"pkg.deepin.io/dde/daemon/loader"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/proxy"
)

// using go build -ldflags "-X main.__VERSION__ version" to set version.
var __VERSION__ = "unknown"

var logger = log.NewLogger("daemon/dde-session-daemon")

func main() {
	InitI18n()
	Textdomain("dde-daemon")

	cmd := kingpin.New("dde-session-daemon", "session daemon")
	cmd.Version("version " + __VERSION__)

	flags := new(Flags)
	flags.Verbose = cmd.Flag("verbose", "Show much more message, the shorthand for --loglevel debug, if specificed, loglevel is ignored.").Short('v').Bool()
	flags.LogLevel = cmd.Flag("loglevel", "Set log level, possible value is error/warn/info/debug/no.").Short('l').String()
	flags.IgnoreMissingModules = cmd.Flag("Ignore", "ignore missing modules, --no-ignore to revert it.").Short('i').Default("true").Bool()
	flags.ForceStart = cmd.Flag("force", "Force start disabled module.").Short('f').Bool()

	cmd.Command("auto", "Automatically get enabled and disabled modules from settings.").Default()
	enablingModules := cmd.Command("enable", "Enable modules and their dependencies, ignore settings.").Arg("module", "module names.").Required().Strings()
	disableModules := cmd.Command("disable", "Disable modules, ignore settings.").Arg("module", "module names.").Required().Strings()
	listModule := cmd.Command("list", "List all the modules or the dependencies of one module.").Arg("module", "module name.").String()

	memprof := cmd.Flag("memprof", "Write memory profile to specific file").String()
	cpuprof := cmd.Flag("cpuprof", "Write cpu profile to specific file").String()

	subCmd := kingpin.MustParse(cmd.Parse(os.Args[1:]))

	if *memprof != "" && *cpuprof != "" {
		logger.Fatal("can not use memprof and cpuprof together")
	}

	C.init()
	proxy.SetupProxy()

	app := NewSessionDaemon(flags, daemonSettings, logger)
	if err := app.register(); err != nil {
		logger.Info(err)
		os.Exit(0)
	}

	if *flags.Verbose {
		app.logLevel = log.LevelDebug
	} else {
		logLevel, err := toLogLevel(*flags.LogLevel)
		if err != nil {
			logger.Error(err)
			return
		}
		app.logLevel = logLevel
	}

	loader.SetLogLevel(app.logLevel)
	logger.Info("LogLevel is", app.logLevel)

	if *memprof != "" {
		startMemProfile(*memprof)
	}

	if *cpuprof != "" {
		startCPUProfile(*cpuprof)
	}

	var err error
	needRunMainLoop := true
	switch subCmd {
	case "auto":
		app.execDefaultAction()
	case "enable":
		err = app.EnableModules(*enablingModules)
	case "disable":
		err = app.DisableModules(*disableModules)
	case "list":
		err = app.ListModule(*listModule)
		needRunMainLoop = false
	}

	if err != nil {
		logger.Warning(err)
		os.Exit(1)
	}

	if needRunMainLoop {
		runMainLoop()
	}
}
