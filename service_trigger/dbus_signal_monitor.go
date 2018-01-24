package service_trigger

import (
	"strings"

	"pkg.deepin.io/lib/dbus1"
)

type DBusSignalMonitor struct {
	Type       uint
	conn       *dbus.Conn
	services   []*Service
	nameMap    map[string]string
	signalChan chan *dbus.Signal
}

func newDBusSignalMonitor(Type uint) *DBusSignalMonitor {
	return &DBusSignalMonitor{
		Type:    Type,
		nameMap: make(map[string]string),
	}
}

func (sigMonitor *DBusSignalMonitor) appendService(service *Service) {
	sigMonitor.services = append(sigMonitor.services, service)
}

func (sigMonitor *DBusSignalMonitor) getConn() (*dbus.Conn, error) {
	if sigMonitor.conn != nil {
		return sigMonitor.conn, nil
	}

	var conn *dbus.Conn
	var err error
	if sigMonitor.Type == busTypeSystem {
		conn, err = dbus.SystemBusPrivate()
	} else {
		conn, err = dbus.SessionBusPrivate()
	}
	if err != nil {
		return nil, err
	}

	err = conn.Auth(nil)
	if err != nil {
		return nil, err
	}

	err = conn.Hello()
	if err != nil {
		return nil, err
	}

	sigMonitor.conn = conn
	return conn, nil
}

const (
	busTypeSystem  = 0
	busTypeSession = 1
)

func (sigMonitor *DBusSignalMonitor) findMatchedServices(signal *dbus.Signal) []*Service {
	var sender string
	if strings.HasPrefix(signal.Sender, ":") {
		sender = sigMonitor.nameMap[signal.Sender]
	} else {
		sender = signal.Sender
	}

	var matched []*Service
	for _, service := range sigMonitor.services {
		dbusField := service.Monitor.DBus
		if dbusField.Sender == sender &&
			signal.Name == dbusField.Interface+"."+dbusField.Signal {

			if dbusField.Path != "" {
				if dbusField.Path == string(signal.Path) {
					matched = append(matched, service)
				}
			} else {
				matched = append(matched, service)
			}
		}
	}
	return matched
}

const ruleNameOwnerChanged = "type='signal'" +
	",sender='org.freedesktop.DBus',path='/org/freedesktop/DBus'" +
	",interface='org.freedesktop.DBus',member='NameOwnerChanged'"

func (sigMonitor *DBusSignalMonitor) init() {
	conn, err := sigMonitor.getConn()
	if err != nil {
		logger.Warning(err)
		return
	}
	rules := []string{ruleNameOwnerChanged}
	for _, service := range sigMonitor.services {
		rules = append(rules, service.getDBusMatchRule())

		dbusField := service.Monitor.DBus

		// set nameMap
		var nameKnown bool
		for _, name := range sigMonitor.nameMap {
			if name == dbusField.Sender {
				nameKnown = true
				break
			}
		}
		if !nameKnown {
			owner, err := getNameOwner(conn, dbusField.Sender)
			if err == nil {
				logger.Debugf("The name %s is owned by %s", dbusField.Sender, owner)
				sigMonitor.nameMap[owner] = dbusField.Sender
			}
		}
	}

	for _, rule := range rules {
		err = addMatch(conn, rule)
		if err != nil {
			logger.Warning(err)
		}
	}

	sigMonitor.signalChan = make(chan *dbus.Signal, 20)
}

func addMatch(conn *dbus.Conn, rule string) error {
	logger.Debug("add rule", rule)
	return conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, rule).Err
}

func (sigMonitor *DBusSignalMonitor) handleNameOwnerChanged(signalBody []interface{}) {
	if len(signalBody) != 3 {
		return
	}

	name, ok := signalBody[0].(string)
	if !ok {
		return
	}
	oldOwner, ok := signalBody[1].(string)
	if !ok {
		return
	}

	newOwner, ok := signalBody[2].(string)
	if !ok {
		return
	}

	if name == "" {
		return
	}

	var found bool
	for _, service := range sigMonitor.services {
		if service.Monitor.DBus.Sender == name {
			found = true
			break
		}
	}

	if !found {
		return
	}

	if oldOwner == "" && newOwner != "" && name != newOwner {
		// newOwner acquire name
		logger.Debugf("The name %s is owned by %s", name, newOwner)
		sigMonitor.nameMap[newOwner] = name
	} else if oldOwner != "" && newOwner == "" && name != oldOwner {
		// oldOwner lost name

		logger.Debugf("The name %s does not have an owner", name)
		delete(sigMonitor.nameMap, oldOwner)
	}
}

func (sigMonitor *DBusSignalMonitor) signalLoop(m *Manager) {
	conn, err := sigMonitor.getConn()
	if err != nil {
		logger.Warning("failed to get conn", sigMonitor.Type)
		return
	}

	conn.Signal(sigMonitor.signalChan)

	for signal := range sigMonitor.signalChan {
		if signal.Sender == "org.freedesktop.DBus" &&
			signal.Name == "org.freedesktop.DBus.NameOwnerChanged" {
			sigMonitor.handleNameOwnerChanged(signal.Body)
		}

		services := sigMonitor.findMatchedServices(signal)
		for _, service := range services {
			logger.Debug("exec service", service)
			go m.execService(service, signal)
		}
	}
	logger.Debug("signalLoop return", sigMonitor.Type)
}

func (sigMonitor *DBusSignalMonitor) stop() error {
	if sigMonitor.conn != nil {
		return sigMonitor.conn.Close()
	}

	return nil
}
