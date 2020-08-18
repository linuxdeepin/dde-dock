package service_trigger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/godbus/dbus"
)

type Service struct {
	filename string
	basename string
	Monitor  struct {
		Type string
		DBus *struct { // dbus signal monitor
			BusType   string // System or Session
			Sender    string
			Interface string
			Signal    string
			Path      string // optional
		}
	}

	Name        string
	Description string
	Exec        []string
}

func (service *Service) getDBusMatchRule() string {
	dbusField := service.Monitor.DBus

	rule := "type='signal'"
	rule += fmt.Sprintf(",sender='%s'", dbusField.Sender)
	rule += fmt.Sprintf(",interface='%s'", dbusField.Interface)
	rule += fmt.Sprintf(",member='%s'", dbusField.Signal)

	if dbusField.Path != "" {
		rule += fmt.Sprintf(",path='%s'", dbusField.Path)
	}
	return rule
}

func (service *Service) check() error {
	if service.Monitor.Type != "DBus" {
		return fmt.Errorf("unknown Monitor.Type %q" + service.Monitor.Type)
	}

	if service.Monitor.Type == "DBus" {
		err := service.checkDBus()
		if err != nil {
			return err
		}
	}

	if service.Name == "" {
		return errors.New("field Name is empty")
	}

	if len(service.Exec) == 0 {
		return errors.New("field Exec is empty")
	}

	return nil
}

func (service *Service) checkDBus() error {
	dbusField := service.Monitor.DBus
	if dbusField == nil {
		return errors.New("field Monitor.DBus is nil")
	}

	if !(dbusField.BusType == "System" ||
		dbusField.BusType == "Session") {
		return errors.New("field Monitor.DBus.BusType is invalid")
	}

	if dbusField.Path != "" {
		if !dbus.ObjectPath(dbusField.Path).IsValid() {
			return errors.New("field Monitor.DBus.Path is invalid")
		}
	}
	if dbusField.Sender == "" {
		return errors.New("field Monitor.DBus.Sender is empty")
	}

	if dbusField.Interface == "" {
		return errors.New("field Monitor.DBus.Interface is empty")
	}

	if dbusField.Signal == "" {
		return errors.New("field Monitor.DBus.Signal is empty")
	}
	return nil
}

func loadService(filename string) (*Service, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var service Service
	err = json.Unmarshal(data, &service)
	if err != nil {
		return nil, err
	}

	err = service.check()
	if err != nil {
		return nil, err
	}

	service.filename = filename
	service.basename = strings.TrimSuffix(filepath.Base(filename), serviceFileExt)
	return &service, nil
}

func (service *Service) String() string {
	if service == nil {
		return "<Service nil>"
	}
	return fmt.Sprintf("<Service %s %q>", service.basename, service.Name)
}
