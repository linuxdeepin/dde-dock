package service_trigger

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"pkg.deepin.io/lib/dbus1"
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

var argReg = regexp.MustCompile(`^%arg(\d+)$`)

func replaceArg(arg string, signal *dbus.Signal) string {
	submatch := argReg.FindStringSubmatch(arg)
	if submatch == nil {
		return arg
	}

	argIdx, err := strconv.Atoi(submatch[1])
	if err != nil {
		logger.Warning(err)
		return arg
	}
	argIdx--
	if argIdx < 0 || argIdx >= len(signal.Body) {
		logger.Warningf("replaceArg: arg %q index out of range", arg)
		return ""
	}
	signalItem := signal.Body[argIdx]
	return fmt.Sprintf("%v", signalItem)
}

func (m *Manager) execService(service *Service, signal *dbus.Signal) {
	var args []string
	for _, arg := range service.Exec[1:] {
		args = append(args, replaceArg(arg, signal))
	}

	logger.Debugf("run cmd %q %#v", service.Exec[0], args)
	cmd := exec.Command(service.Exec[0], args...)
	err := cmd.Run()
	if err != nil {
		logger.Warning(err)
	}
}
