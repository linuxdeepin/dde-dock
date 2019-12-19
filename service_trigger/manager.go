package service_trigger

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/log"
)

type Manager struct {
	serviceMap map[string]*Service

	systemSigMonitor  *DBusSignalMonitor
	sessionSigMonitor *DBusSignalMonitor
}

func newManager() *Manager {
	m := &Manager{
		systemSigMonitor:  newDBusSignalMonitor(busTypeSystem),
		sessionSigMonitor: newDBusSignalMonitor(busTypeSession),
	}
	return m
}

func (m *Manager) start() {
	m.loadServices()

	m.sessionSigMonitor.init()
	go m.sessionSigMonitor.signalLoop(m)

	m.systemSigMonitor.init()
	go m.systemSigMonitor.signalLoop(m)
}

func (m *Manager) stop() error {
	err := m.sessionSigMonitor.stop()
	if err != nil {
		return err
	}

	return m.systemSigMonitor.stop()
}

func (m *Manager) loadServices() {
	m.serviceMap = make(map[string]*Service)
	m.loadServicesFromDir("/usr/lib/deepin-daemon/" + moduleName)
	m.loadServicesFromDir("/etc/deepin-daemon/" + moduleName)

	for _, service := range m.serviceMap {
		if service.Monitor.Type == "DBus" {
			dbusField := service.Monitor.DBus
			if dbusField.BusType == "System" {
				m.systemSigMonitor.appendService(service)
			} else if dbusField.BusType == "Session" {
				m.sessionSigMonitor.appendService(service)
			}
		}
	}
}

const serviceFileExt = ".service.json"

func (m *Manager) loadServicesFromDir(dirname string) {
	fileInfoList, _ := ioutil.ReadDir(dirname)
	for _, fileInfo := range fileInfoList {
		if fileInfo.IsDir() {
			continue
		}

		name := fileInfo.Name()
		if !strings.HasSuffix(name, serviceFileExt) {
			continue
		}

		filename := filepath.Join(dirname, name)
		service, err := loadService(filename)
		if err != nil {
			logger.Warningf("failed to load %q: %v", filename, err)
			continue
		} else {
			logger.Debugf("load %q ok", filename)
		}

		_, ok := m.serviceMap[name]
		if ok {
			logger.Debugf("file %q overwrites the old", filename)
		}
		m.serviceMap[name] = service
	}
}

func getNameOwner(conn *dbus.Conn, name string) (string, error) {
	var owner string
	err := conn.BusObject().Call("org.freedesktop.DBus.GetNameOwner",
		0, name).Store(&owner)
	if err != nil {
		return "", err
	}
	return owner, err
}

func newReplacer(signal *dbus.Signal) *strings.Replacer {
	var oldNewSlice []string
	for idx, item := range signal.Body {
		oldStr := fmt.Sprintf("%%{arg%d}", idx)
		newStr := fmt.Sprintf("%v", item)
		logger.Debugf("old %q => new %q", oldStr, newStr)
		oldNewSlice = append(oldNewSlice, oldStr, newStr)
	}
	return strings.NewReplacer(oldNewSlice...)
}

func (m *Manager) execService(service *Service, signal *dbus.Signal) {
	if len(service.Exec) == 0 {
		logger.Warning("service Exec empty")
		return
	}

	var args []string
	execArgs := service.Exec[1:]
	if logger.GetLogLevel() == log.LevelDebug {
		if service.Exec[0] == "sh" {
			// add -x option for debug shell
			execArgs = append([]string{"-x"}, execArgs...)
		}
	}

	replacer := newReplacer(signal)
	for _, arg := range execArgs {
		args = append(args, replacer.Replace(arg))
	}

	logger.Debugf("run cmd %q %#v", service.Exec[0], args)
	cmd := exec.Command(service.Exec[0], args...)
	out, err := cmd.CombinedOutput()
	logger.Debugf("cmd combined output: %s", out)
	if err != nil {
		logger.Warning(err)
	}
}
