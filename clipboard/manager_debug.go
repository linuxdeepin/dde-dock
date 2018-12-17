package clipboard

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/linuxdeepin/go-x11-client"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

func (m *Manager) saveClipboard() error {
	owner, err := getSelectionOwner(m.xConn, atomClipboard)
	if err != nil {
		return err
	}

	logger.Debug("clipboard selection owner:", owner)

	ts, err := m.getTimestamp()
	if err != nil {
		return err
	}

	targets, err := m.getClipboardTargets(ts)
	if err != nil {
		return err
	}
	logger.Debug("targets:", targets)

	m.saveTargets(targets, ts)
	m.contentMu.Lock()
	for _, targetData := range m.content {
		logger.Debugf("target %d type: %v", targetData.Target, targetData.Type)
	}
	m.contentMu.Unlock()

	return nil
}

func (m *Manager) SaveClipboard() *dbus.Error {
	err := m.saveClipboard()
	return dbusutil.ToError(err)
}

func (m *Manager) writeContent() error {
	dir := "/tmp/dde-session-daemon-clipboard"

	err := os.Mkdir(dir, 0700)
	if err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	err = emptyDir(dir)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	m.contentMu.Lock()
	for _, targetData := range m.content {
		target := targetData.Target
		targetName, _ := m.xConn.GetAtomName(target)
		_, err = fmt.Fprintf(&buf, "%d,%s\n", target, targetName)
		if err != nil {
			m.contentMu.Unlock()
			return err
		}

		err = ioutil.WriteFile(filepath.Join(dir, strconv.Itoa(int(target))), targetData.Data, 0644)
		if err != nil {
			m.contentMu.Unlock()
			return err
		}
	}
	m.contentMu.Unlock()
	err = ioutil.WriteFile(filepath.Join(dir, "index.txt"), buf.Bytes(), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) WriteContent() *dbus.Error {
	err := m.writeContent()
	return dbusutil.ToError(err)
}

func (m *Manager) BecomeClipboardOwner() *dbus.Error {
	ts, err := m.getTimestamp()
	if err != nil {
		return dbusutil.ToError(err)
	}
	err = m.becomeClipboardOwner(ts)
	return dbusutil.ToError(err)
}

func (m *Manager) removeTarget(target x.Atom) {
	m.contentMu.Lock()
	newContent := make([]*TargetData, 0, len(m.content))
	for _, td := range m.content {
		if td.Target != target {
			newContent = append(newContent, td)
		}
	}
	m.content = newContent
	m.contentMu.Unlock()
}

func (m *Manager) RemoveTarget(target uint32) *dbus.Error {
	m.removeTarget(x.Atom(target))
	return nil
}
