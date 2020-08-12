package clipboard

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/xfixes"
	"pkg.deepin.io/lib/log"
)

var (
	atomClipboardManager     x.Atom
	atomClipboard            x.Atom
	atomSaveTargets          x.Atom
	atomTargets              x.Atom
	atomMultiple             x.Atom
	atomDelete               x.Atom
	atomInsertProperty       x.Atom
	atomInsertSelection      x.Atom
	atomAtomPair             x.Atom //nolint
	atomIncr                 x.Atom
	atomTimestamp            x.Atom
	atomNull                 x.Atom //nolint
	atomTimestampProp        x.Atom
	atomFromClipboardManager x.Atom

	selectionMaxSize int
)

func initAtoms(xConn *x.Conn) {
	atomClipboardManager, _ = xConn.GetAtom("CLIPBOARD_MANAGER")
	atomClipboard, _ = xConn.GetAtom("CLIPBOARD")
	atomSaveTargets, _ = xConn.GetAtom("SAVE_TARGETS")
	atomTargets, _ = xConn.GetAtom("TARGETS")
	atomMultiple, _ = xConn.GetAtom("MULTIPLE")
	atomDelete, _ = xConn.GetAtom("DELETE")
	atomInsertProperty, _ = xConn.GetAtom("INSERT_PROPERTY")
	atomInsertSelection, _ = xConn.GetAtom("INSERT_SELECTION")
	atomAtomPair, _ = xConn.GetAtom("ATOM_PAIR")
	atomIncr, _ = xConn.GetAtom("INCR")
	atomTimestamp, _ = xConn.GetAtom("TIMESTAMP")
	atomTimestampProp, _ = xConn.GetAtom("_TIMESTAMP_PROP")
	atomNull, _ = xConn.GetAtom("NULL")
	atomFromClipboardManager, _ = xConn.GetAtom("FROM_DEEPIN_CLIPBOARD_MANAGER")
	selectionMaxSize = 65432
	logger.Debug("selectionMaxSize:", selectionMaxSize)
}

type TargetData struct {
	Target x.Atom
	Type   x.Atom
	Format uint8
	Data   []byte
}

func (td *TargetData) needINCR() bool {
	return len(td.Data) > selectionMaxSize
}

type Manager struct {
	xc        XClient
	window    x.Window
	ec        *eventCaptor
	timestamp x.Timestamp

	contentMu sync.Mutex
	content   []*TargetData
	//nolint
	methods *struct {
		RemoveTarget func() `in:"target"`
	}
}

func (m *Manager) getTargetData(target x.Atom) *TargetData {
	m.contentMu.Lock()
	defer m.contentMu.Unlock()

	for _, td := range m.content {
		if td.Target == target {
			return td
		}
	}
	return nil
}

func (m *Manager) addTargetData(targetData *TargetData) {
	m.contentMu.Lock()
	defer m.contentMu.Unlock()

	for idx, td := range m.content {
		if td.Target == targetData.Target {
			m.content[idx] = targetData
			return
		}
	}
	m.content = append(m.content, targetData)
}

func (m *Manager) start() error {
	owner, err := m.xc.GetSelectionOwner(atomClipboardManager)
	if err != nil {
		return err
	}
	if owner != 0 {
		return fmt.Errorf("another clipboard manager is already running, owner: %d", owner)
	}

	m.window, err = m.xc.CreateWindow()
	if err != nil {
		return err
	}
	logger.Debug("m.window:", m.window)

	err = m.xc.SelectSelectionInputE(m.window, atomClipboard,
		xfixes.SelectionEventMaskSetSelectionOwner|
			xfixes.SelectionEventMaskSelectionClientClose|
			xfixes.SelectionEventMaskSelectionWindowDestroy)
	if err != nil {
		logger.Warning(err)
	}

	m.ec = newEventCaptor()
	eventChan := make(chan x.GenericEvent, 50)
	m.xc.Conn().AddEventChan(eventChan)
	go func() {
		for ev := range eventChan {
			m.handleEvent(ev)
		}
	}()

	ts, err := m.getTimestamp()
	if err != nil {
		return err
	}
	m.timestamp = ts

	logger.Debug("ts:", ts)
	err = setSelectionOwner(m.xc, m.window, atomClipboardManager, ts)
	if err != nil {
		return err
	}

	err = announceManageSelection(m.xc.Conn(), m.window, atomClipboardManager, ts)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) handleEvent(ev x.GenericEvent) {
	xfixesExtData := m.xc.Conn().GetExtensionData(xfixes.Ext())
	code := ev.GetEventCode()
	switch code {
	case x.SelectionRequestEventCode:
		event, _ := x.NewSelectionRequestEvent(ev)
		logger.Debug(selReqEventToString(event))

		if event.Selection == atomClipboardManager {
			go m.convertClipboardManager(event)
		} else if event.Selection == atomClipboard {
			go m.convertClipboard(event)
		}

	case x.PropertyNotifyEventCode:
		event, _ := x.NewPropertyNotifyEvent(ev)

		if m.ec.handleEvent(event) {
			logger.Debug("->", propNotifyEventToString(event))
			return
		}
		logger.Debug(">>", propNotifyEventToString(event))

	case x.SelectionNotifyEventCode:
		event, _ := x.NewSelectionNotifyEvent(ev)

		if m.ec.handleEvent(event) {
			logger.Debug("->", selNotifyEventToString(event))
			return
		}
		logger.Debug(">>", selNotifyEventToString(event))

	case x.DestroyNotifyEventCode:
		event, _ := x.NewDestroyNotifyEvent(ev)

		logger.Debug(destroyNotifyEventToString(event))

	case x.SelectionClearEventCode:
		event, _ := x.NewSelectionClearEvent(ev)
		logger.Debug(selClearEventToString(event))

	case xfixes.SelectionNotifyEventCode + xfixesExtData.FirstEvent:
		event, _ := xfixes.NewSelectionNotifyEvent(ev)
		logger.Debug(xfixesSelNotifyEventToString(event))
		switch event.Subtype {
		case xfixes.SelectionEventSetSelectionOwner:
			if event.Selection == atomClipboard {
				if event.Owner == m.window {
					logger.Debug("i have become the owner of CLIPBOARD")
				} else {
					logger.Debug("other app have become the owner of CLIPBOARD")
				}
			}

		case xfixes.SelectionEventSelectionWindowDestroy, xfixes.SelectionEventSelectionClientClose:
			err := m.becomeClipboardOwner(event.Timestamp)
			if err != nil {
				logger.Warning(err)
			}
		}
	}
}

func setSelectionOwner(xc XClient, win x.Window, selection x.Atom, ts x.Timestamp) error {
	xc.SetSelectionOwner(win, selection, ts)
	owner, err := xc.GetSelectionOwner(selection)
	if err != nil {
		return err
	}
	if owner != win {
		return errors.New("failed to set selection owner")
	}
	return nil
}

func (m *Manager) becomeClipboardOwner(ts x.Timestamp) error {
	err := setSelectionOwner(m.xc, m.window, atomClipboard, ts)
	if err != nil {
		return err
	}
	logger.Debug("set clipboard selection owner to me")
	return nil
}

func (m *Manager) getClipboardTargets(ts x.Timestamp) ([]x.Atom, error) {
	selNotifyEvent, err := m.ec.captureSelectionNotifyEvent(func() error {
		m.xc.ConvertSelection(m.window, atomClipboard,
			atomTargets, atomTargets, ts)
		return m.xc.Flush()
	}, func(event *x.SelectionNotifyEvent) bool {
		return event.Target == atomTargets &&
			event.Selection == atomClipboard &&
			event.Requestor == m.window
	})
	if err != nil {
		return nil, err
	}

	if selNotifyEvent.Property == x.None {
		return nil, errors.New("failed to convert clipboard targets")
	}

	propReply, err := m.getProperty(m.window, selNotifyEvent.Property, true)
	if err != nil {
		return nil, err
	}

	targets, err := getAtomListFormReply(propReply)
	if err != nil {
		return nil, err
	}

	return targets, nil
}

// convert CLIPBOARD_MANAGER selection
func (m *Manager) convertClipboardManager(ev *x.SelectionRequestEvent) {
	logger.Debug("convert CLIPBOARD_MANAGER selection")
	switch ev.Target {
	case atomSaveTargets:
		logger.Debug("SAVE_TARGETS")
		err := m.xc.ChangeWindowEventMask(ev.Requestor, x.EventMaskStructureNotify)
		if err != nil {
			logger.Warning(err)
			m.finishSelectionRequest(ev, false)
			return
		}

		var targets []x.Atom
		var replyType x.Atom
		if ev.Property != x.None {
			reply, err := m.xc.GetProperty(true, ev.Requestor, ev.Property,
				x.AtomAtom, 0, 0x1FFFFFFF)
			if err != nil {
				logger.Warning(err)
				m.finishSelectionRequest(ev, false)
				return
			}

			replyType = reply.Type
			if reply.Type != x.None {
				targets, err = getAtomListFormReply(reply)
				if err != nil {
					logger.Warning(err)
					m.finishSelectionRequest(ev, false)
					return
				}
			}
		}

		if replyType == x.None {
			logger.Debugf("need convert clipboard targets")
			targets, err = m.getClipboardTargets(ev.Time)
			if err != nil {
				logger.Warning(err)
				m.finishSelectionRequest(ev, false)
				return
			}
		}

		m.saveTargets(targets, ev.Time)
		// add special target
		m.addTargetData(&TargetData{
			Target: atomFromClipboardManager,
			Type:   x.AtomString,
			Format: 8,
			Data:   []byte("1"),
		})

		m.finishSelectionRequest(ev, true)

	case atomTargets:
		w := x.NewWriter()
		w.Write4b(uint32(atomTargets))
		w.Write4b(uint32(atomSaveTargets))
		w.Write4b(uint32(atomTimestamp))
		err := m.xc.ChangePropertyE(x.PropModeReplace, ev.Requestor,
			ev.Property, x.AtomAtom, 32, w.Bytes())
		if err != nil {
			logger.Warning(err)
		}
		m.finishSelectionRequest(ev, err == nil)

	case atomTimestamp:
		w := x.NewWriter()
		w.Write4b(uint32(m.timestamp))
		err := m.xc.ChangePropertyE(x.PropModeReplace, ev.Requestor,
			ev.Property, x.AtomInteger, 32, w.Bytes())
		if err != nil {
			logger.Warning(err)
		}
		m.finishSelectionRequest(ev, err == nil)

	default:
		m.finishSelectionRequest(ev, false)
	}
}

// convert CLIPBOARD selection
func (m *Manager) convertClipboard(ev *x.SelectionRequestEvent) {
	targetName, _ := m.xc.GetAtomName(ev.Target)
	logger.Debugf("convert clipboard target %s %d", targetName, ev.Target)

	if ev.Target == atomTargets {
		w := x.NewWriter()
		w.Write4b(uint32(atomTargets))
		m.contentMu.Lock()
		for _, targetData := range m.content {
			w.Write4b(uint32(targetData.Target))
		}
		m.contentMu.Unlock()

		err := m.xc.ChangePropertyE(x.PropModeReplace, ev.Requestor,
			ev.Property, x.AtomAtom, 32, w.Bytes())
		if err != nil {
			logger.Warning(err)
		}
		m.finishSelectionRequest(ev, err == nil)

	} else {
		targetData := m.getTargetData(ev.Target)
		if targetData == nil {
			m.finishSelectionRequest(ev, false)
			return
		}

		if targetData.needINCR() {
			err := m.sendTargetIncr(targetData, ev)
			if err != nil {
				logger.Warning(err)
			}
		} else {
			err := m.xc.ChangePropertyE(x.PropModeReplace, ev.Requestor,
				ev.Property, targetData.Type, targetData.Format, targetData.Data)
			if err != nil {
				logger.Warning(err)
			}
			m.finishSelectionRequest(ev, err == nil)
		}
	}
}

func (m *Manager) sendTargetIncr(targetData *TargetData, ev *x.SelectionRequestEvent) error {
	err := m.xc.ChangeWindowEventMask(ev.Requestor, x.EventMaskPropertyChange)
	if err != nil {
		return err
	}

	_, err = m.ec.capturePropertyNotifyEvent(func() error {
		w := x.NewWriter()
		w.Write4b(uint32(len(targetData.Data)))
		err = m.xc.ChangePropertyE(x.PropModeReplace, ev.Requestor, ev.Property,
			atomIncr, 32, w.Bytes())
		if err != nil {
			logger.Warning(err)
		}
		m.finishSelectionRequest(ev, err == nil)
		return err
	}, func(pev *x.PropertyNotifyEvent) bool {
		return pev.Window == ev.Requestor &&
			pev.State == x.PropertyDelete &&
			pev.Atom == ev.Property
	})

	if err != nil {
		return err
	}

	var offset int
	for {
		data := targetData.Data[offset:]
		length := len(data)
		if length > selectionMaxSize {
			length = selectionMaxSize
		}
		offset += length

		_, err = m.ec.capturePropertyNotifyEvent(func() error {
			logger.Debug("send incr data", length)
			err = m.xc.ChangePropertyE(x.PropModeReplace, ev.Requestor, ev.Property,
				targetData.Type, targetData.Format, data[:length])
			if err != nil {
				logger.Warning(err)
			}
			return err
		}, func(pev *x.PropertyNotifyEvent) bool {
			return pev.Window == ev.Requestor &&
				pev.State == x.PropertyDelete &&
				pev.Atom == ev.Property
		})

		if length == 0 {
			break
		}
	}

	return nil
}

func (m *Manager) finishSelectionRequest(ev *x.SelectionRequestEvent, success bool) {
	var property x.Atom
	if success {
		property = ev.Property
	}

	event := &x.SelectionNotifyEvent{
		Time:      ev.Time,
		Requestor: ev.Requestor,
		Selection: ev.Selection,
		Target:    ev.Target,
		Property:  property,
	}

	err := m.xc.SendEventE(false, ev.Requestor, x.EventMaskNoEvent,
		event)
	if err != nil {
		logger.Warning(err)
	}

	logger.Debugf("finish selection request %v {Requestor: %d, Selection: %d,"+
		" Target: %d, Property: %d}",
		success, ev.Requestor, ev.Selection, ev.Target, ev.Property)
}

func (m *Manager) saveTargets(targets []x.Atom, ts x.Timestamp) {
	m.contentMu.Lock()
	m.content = nil
	m.contentMu.Unlock()

	for _, target := range targets {
		targetName, err := m.xc.GetAtomName(target)
		if err != nil {
			logger.Warning(err)
			continue
		}
		if shouldIgnoreSaveTarget(target, targetName) {
			logger.Debug("ignore target", target, targetName)
			continue
		}

		logger.Debug("save target", target, targetName)
		err = m.saveTarget(target, ts)
		if err != nil {
			logger.Warning(err)
		}
	}
}

func shouldIgnoreSaveTarget(target x.Atom, targetName string) bool {
	switch target {
	case atomTargets, atomSaveTargets,
		atomTimestamp, atomMultiple, atomDelete,
		atomInsertProperty, atomInsertSelection,
		x.AtomPixmap:
		return true
	}
	if strings.HasPrefix(targetName, "image/") {
		switch targetName {
		case "image/jpeg", "image/png", "image/bmp":
			return false
		default:
			return true
		}
	}
	return false
}

func (m *Manager) saveTarget(target x.Atom, ts x.Timestamp) error {
	selNotifyEvent, err := m.ec.captureSelectionNotifyEvent(func() error {
		m.xc.ConvertSelection(m.window, atomClipboard, target, target, ts)
		return m.xc.Flush()
	}, func(event *x.SelectionNotifyEvent) bool {
		return event.Selection == atomClipboard &&
			event.Requestor == m.window &&
			event.Target == target
	})
	if err != nil {
		return err
	}
	if selNotifyEvent.Property == x.None {
		return errors.New("failed to convert target")
	}

	propReply, err := m.getProperty(m.window, selNotifyEvent.Property, false)
	if err != nil {
		return err
	}

	if propReply.Type == atomIncr {
		err := m.recvTargetIncr(target, selNotifyEvent.Property)
		if err != nil {
			return err
		}
	} else {
		err = m.xc.DeletePropertyE(m.window, selNotifyEvent.Property)
		if err != nil {
			return err
		}
		logger.Debug("data len:", len(propReply.Value))
		m.addTargetData(&TargetData{
			Target: target,
			Type:   propReply.Type,
			Format: propReply.Format,
			Data:   propReply.Value,
		})
	}
	return nil
}

func (m *Manager) getProperty(win x.Window, propertyAtom x.Atom, delete bool) (*x.GetPropertyReply, error) {
	propReply, err := m.xc.GetProperty(false, win, propertyAtom,
		x.GetPropertyTypeAny, 0, 0)
	if err != nil {
		return nil, err
	}

	propReply, err = m.xc.GetProperty(delete, win, propertyAtom,
		x.GetPropertyTypeAny,
		0,
		(propReply.BytesAfter+uint32(x.Pad(int(propReply.BytesAfter))))/4,
	)
	if err != nil {
		return nil, err
	}
	return propReply, nil
}

func (m *Manager) recvTargetIncr(target, prop x.Atom) error {
	logger.Debug("start recvTargetIncr", target)
	var data [][]byte
	t0 := time.Now()
	total := 0
	for {
		propNotifyEvent, err := m.ec.capturePropertyNotifyEvent(func() error {
			err := m.xc.DeletePropertyE(m.window, prop)
			if err != nil {
				logger.Warning(err)
			}
			return err

		}, func(event *x.PropertyNotifyEvent) bool {
			return event.State == x.PropertyNewValue && event.Window == m.window &&
				event.Atom == prop
		})
		if err != nil {
			logger.Warning(err)
			return err
		}

		propReply, err := m.xc.GetProperty(false, propNotifyEvent.Window, propNotifyEvent.Atom,
			x.GetPropertyTypeAny,
			0, 0)
		if err != nil {
			logger.Warning(err)
			return err
		}
		propReply, err = m.xc.GetProperty(false, propNotifyEvent.Window, propNotifyEvent.Atom,
			x.GetPropertyTypeAny, 0,
			(propReply.BytesAfter+uint32(x.Pad(int(propReply.BytesAfter))))/4,
		)
		if err != nil {
			logger.Warning(err)
			return err
		}

		if propReply.ValueLen == 0 {
			logger.Debugf("end recvTargetIncr %d, took %v, total size: %d",
				target, time.Since(t0), total)

			err = m.xc.DeletePropertyE(propNotifyEvent.Window, propNotifyEvent.Atom)
			if err != nil {
				logger.Warning(err)
				return err
			}

			m.addTargetData(&TargetData{
				Target: target,
				Type:   propReply.Type,
				Format: propReply.Format,
				Data:   bytes.Join(data, nil),
			})
			return nil
		}
		if logger.GetLogLevel() == log.LevelDebug {
			logger.Debugf("recv data size: %d, md5sum: %s", len(propReply.Value), getBytesMd5sum(propReply.Value))
		}
		total += len(propReply.Value)
		data = append(data, propReply.Value)
	}
}

func (m *Manager) getTimestamp() (x.Timestamp, error) {
	propNotifyEvent, err := m.ec.capturePropertyNotifyEvent(func() error {
		return m.xc.ChangePropertyE(x.PropModeReplace, m.window, atomTimestampProp,
			x.AtomInteger, 32, nil)
	}, func(event *x.PropertyNotifyEvent) bool {
		return event.Window == m.window &&
			event.Atom == atomTimestampProp &&
			event.State == x.PropertyNewValue
	})

	if err != nil {
		return 0, err
	}

	return propNotifyEvent.Time, nil
}

func (m *Manager) GetInterfaceName() string {
	return dbusServiceName
}
