package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"pkg.deepin.io/dde-daemon/loader"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/lib/log"
	"strings"
)

func runMainLoop() {
	ddeSessionRegister()
	dbus.DealWithUnhandledMessage()
	glib.StartLoop()
}

func toBool(v string) bool {
	return v == "true"
}

func toLogLevel(name string) (log.Priority, error) {
	name = strings.ToLower(name)
	logLevel := log.LevelInfo
	var err error
	switch name {
	case "":
	case "error":
		logLevel = log.LevelError
	case "warn":
		logLevel = log.LevelWarning
	case "info":
		logLevel = log.LevelInfo
	case "debug":
		logLevel = log.LevelDebug
	case "no":
		logLevel = log.LevelDisable
	default:
		err = fmt.Errorf("%s is not support", name)
	}

	return logLevel, err
}

type Flags struct {
	Verbose              *bool
	LogLevel             *string
	IgnoreMissingModules *bool
	ForceStart           *bool
}

func getEnableFlag(flag *Flags) loader.EnableFlag {
	enableFlag := loader.EnableFlagIgnoreMissingModule

	if *flag.IgnoreMissingModules {
		enableFlag = loader.EnableFlagNone
	}

	if *flag.ForceStart {
		enableFlag |= loader.EnableFlagForceStart
	}

	return enableFlag
}

type SessionDaemon struct {
	cmd             *kingpin.Application
	ctx             *kingpin.ParseContext
	flags           *Flags
	logLevel        log.Priority
	log             *log.Logger
	settings        *gio.Settings
	enabledModules  map[string]loader.Module
	disabledModules map[string]loader.Module
}

func NewSessionDaemon(cmd *kingpin.Application, flags *Flags, settings *gio.Settings, logger *log.Logger) *SessionDaemon {
	session := &SessionDaemon{
		cmd:             cmd,
		flags:           flags,
		settings:        settings,
		logLevel:        log.LevelInfo,
		log:             logger,
		enabledModules:  map[string]loader.Module{},
		disabledModules: map[string]loader.Module{},
	}

	session.initModules()

	return session
}

func (s *SessionDaemon) exitIfNotSingleton() {
	if !lib.UniqueOnSession("com.deepin.daemon") {
		s.log.Warning("There already has a dde-daemon running.")
		os.Exit(0)
	}

}

func (s *SessionDaemon) parse() (string, error) {
	// kingpin doesn't support customized default action for now.
	// If no subcommand is selected, usage will be printed.

	// check the validity of command line arguments first.
	ctx, err := s.cmd.ParseContext(os.Args[1:])
	s.ctx = ctx
	if err != nil {
		return "", err
	}

	// using "" to indicate no subcommand is selected.
	if s.ctx.SelectedCommand == nil {
		return "", nil
	}

	return kingpin.MustParse(s.cmd.Parse(os.Args[1:])), nil
}

func (s *SessionDaemon) printUsage() {
	s.cmd.UsageForContext(s.ctx)
}

func (s *SessionDaemon) version() string {
	model := s.cmd.Model()
	return fmt.Sprintf("%s %s", model.Name, model.Version)
}

func (s *SessionDaemon) defaultAction() {
	loader.SetLogLevel(s.logLevel)

	err := loader.EnableModules(s.getEnabledModules(), s.getDisabledModules(), getEnableFlag(s.flags))
	if err != nil {
		fmt.Println(err)
		// TODO: define exit code.
		os.Exit(3)
	}

	go func() {
		if err := dbus.Wait(); err != nil {
			logger.Errorf("Lost dbus: %v", err)
			os.Exit(-1)
		} else {
			logger.Info("dbus connection is closed by user")
			os.Exit(0)
		}
	}()

	runMainLoop()
}

func (s *SessionDaemon) execDefaultAction() {
	verboseExist := false

	// If no subcommand is selected, flags won't be parsed, do it ourself.
	// Just handle the first flag, otherwise print usage.
	for _, el := range s.ctx.Elements {
		switch c := el.Clause.(type) {
		case *kingpin.FlagClause:
			switch c.Model().Name {
			case "help":
				s.printUsage()
				return
			case "version":
				fmt.Println(s.version())
				return
			case "ignore":
				*s.flags.IgnoreMissingModules = toBool(*el.Value)
			case "verbose":
				verboseExist = true
				*s.flags.Verbose = toBool(*el.Value)
				if *s.flags.Verbose {
					s.logLevel = log.LevelDebug
				}
			case "loglevel":
				if verboseExist {
					continue
				}
				logLevel, err := toLogLevel(*el.Value)
				if err != nil {
					fmt.Println(err)
					return
				}
				s.logLevel = logLevel
			case "force":
				*s.flags.ForceStart = toBool(*el.Value)
			}
		case *kingpin.ArgClause:
			fmt.Println("not support ArgClause now:", c.Model().Name)
			os.Exit(2)
		case *kingpin.CmdClause:
			fmt.Println("not support CmdClause now:", c.Model().Name)
			os.Exit(2)
		}
	}

	s.exitIfNotSingleton()
	s.defaultAction()
}

func (s *SessionDaemon) initModules() {
	allModules := loader.List()
	for _, module := range allModules {
		name := module.Name()
		if s.settings.GetBoolean(name) {
			s.enabledModules[name] = module
		} else {
			s.disabledModules[name] = module
		}
	}
}

func keys(m map[string]loader.Module) []string {
	keys := []string{}
	for key, _ := range m {
		keys = append(keys, key)
	}
	return keys
}

func (s *SessionDaemon) getDisabledModules() []string {
	return keys(s.disabledModules)
}

func (s *SessionDaemon) getEnabledModules() []string {
	return keys(s.enabledModules)
}

func (s *SessionDaemon) EnableModules(enablingModules []string) error {
	return loader.EnableModules(enablingModules, s.getDisabledModules(), getEnableFlag(s.flags))
}

func (s *SessionDaemon) DisableModules(disableModules []string) error {
	return loader.EnableModules(s.getEnabledModules(), disableModules, getEnableFlag(s.flags))
}

func (s *SessionDaemon) ListModule(name string) error {
	if name == "" {
		for _, module := range loader.List() {
			fmt.Println(module.Name())
		}
		return nil
	}

	module := loader.GetModule(name)
	if module == nil {
		return fmt.Errorf("no such a module named %s", name)
	}

	for _, m := range module.GetDependencies() {
		fmt.Println(m)
	}

	return nil
}
