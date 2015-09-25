package main

import (
	"fmt"
	"os"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/lib/log"
)

func runMainLoop() {
	ddeSessionRegister()
	dbus.DealWithUnhandledMessage()
	go glib.StartLoop()

	if err := dbus.Wait(); err != nil {
		logger.Errorf("Lost dbus: %v", err)
		os.Exit(-1)
	}

	logger.Info("dbus connection is closed by user")
	os.Exit(0)
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
	flags           *Flags
	logLevel        log.Priority
	log             *log.Logger
	settings        *gio.Settings
	enabledModules  map[string]loader.Module
	disabledModules map[string]loader.Module
}

func (*SessionDaemon) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Daemon",
		ObjectPath: "/com/deepin/daemon/Daemon",
		Interface:  "com.deepin.daemon.Daemon",
	}
}

func NewSessionDaemon(flags *Flags, settings *gio.Settings, logger *log.Logger) *SessionDaemon {
	session := &SessionDaemon{
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
	if !lib.UniqueOnSession(s.GetDBusInfo().Dest) {
		s.log.Warning("There already has a dde daemon running.")
		os.Exit(0)
	}

	if err := dbus.InstallOnSession(s); err != nil {
		logger.Fatal(err)
	}

}

func (s *SessionDaemon) defaultAction() {
	loader.SetLogLevel(s.logLevel)

	err := loader.EnableModules(s.getEnabledModules(), s.getDisabledModules(), getEnableFlag(s.flags))
	if err != nil {
		fmt.Println(err)
		// TODO: define exit code.
		os.Exit(3)
	}
}

func (s *SessionDaemon) execDefaultAction() {
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
