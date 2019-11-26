package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"sync"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/pam"
)

type PAMTransaction struct {
	baseTransaction
	PropsMu        sync.RWMutex
	Authenticating bool
	Sender         string
	methods        *struct {
		SetUser func() `in:"user"`
	}

	core    *pam.Transaction
	markEnd bool
}

func (tx *PAMTransaction) setPropAuthenticating(value bool) {
	if tx.Authenticating != value {
		tx.Authenticating = value
		err := tx.parent.service.EmitPropertyChanged(tx, "Authenticating", value)
		if err != nil {
			logger.Warning(err)
		}
	}
}

func (tx *PAMTransaction) RespondPAM(style pam.Style, msg string) (string, error) {
	switch style {
	case pam.PromptEchoOn:
		result, err := tx.requestEchoOn(msg)
		if err != nil {
			logger.Warning(err)
		} else {
			logger.Debug("RequestEchoOn result:", result)
			tx.setUser(result)
		}
		return result, err

	case pam.PromptEchoOff:
		result, err := tx.requestEchoOff(msg)
		if err != nil {
			logger.Warning(err)
		}
		tx.setAuthToken(result)
		return result, err

	case pam.ErrorMsg:
		err := tx.displayErrorMsg(msg)
		if err != nil {
			logger.Warning(err)
		}
		return "", nil
	case pam.TextInfo:
		err := tx.displayTextInfo(msg)
		if err != nil {
			logger.Warning(err)
		}
		return "", nil
	default:
		return "", errors.New("unrecognized message style")
	}
}

func genCookie() (string, error) {
	var buf = make([]byte, 256)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write(buf)
	encoded := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	return encoded, nil
}

func (tx *PAMTransaction) authenticate() error {
	logger.Debug(tx, "authenticate")
	tx.PropsMu.Lock()
	tx.setPropAuthenticating(true)
	tx.PropsMu.Unlock()

	// meet the requirement of pam_unix.so nullok_secure option,
	// allows any user with a blank password to unlock.
	err := tx.core.SetItemStr(pam.Tty, "tty1")
	if err != nil {
		logger.Warning("failed to set item tty:", err)
	}

	err = tx.core.Authenticate(0)

	tx.PropsMu.Lock()
	defer tx.PropsMu.Unlock()

	tx.setPropAuthenticating(false)
	if tx.markEnd {
		tx.terminate()
		return errors.New("mark end")
	}
	return err
}

func (tx *PAMTransaction) Authenticate(sender dbus.Sender) *dbus.Error {
	tx.parent.service.DelayAutoQuit()
	if err := tx.checkSender(sender); err != nil {
		return err
	}

	tx.PropsMu.Lock()
	defer tx.PropsMu.Unlock()

	if tx.Authenticating {
		return dbusutil.ToError(errors.New("transaction busy"))
	} else {
		go func() {
			err := tx.authenticate()
			tx.sendResult(err == nil)
			if err != nil {
				logger.Warning(err)
			}
		}()
	}
	return nil
}

func (tx *PAMTransaction) terminate() {
	logger.Debug(tx, "terminate")
	err := tx.core.End(tx.core.LastStatus())
	if err != nil {
		logger.Warning(err)
	}
	tx.parent.deleteTx(tx.id)
}

func (tx *PAMTransaction) End(sender dbus.Sender) *dbus.Error {
	tx.parent.service.DelayAutoQuit()
	if err := tx.checkSender(sender); err != nil {
		return err
	}

	logger.Debugf("%s End sender: %s", tx, sender)
	tx.clearSecret()
	tx.PropsMu.Lock()
	defer tx.PropsMu.Unlock()

	if tx.Authenticating {
		tx.markEnd = true
	} else {
		tx.terminate()
	}
	return nil
}

func (tx *PAMTransaction) SetUser(sender dbus.Sender, user string) *dbus.Error {
	tx.parent.service.DelayAutoQuit()
	if err := tx.checkSender(sender); err != nil {
		return err
	}
	logger.Debug(tx, "SetUser", sender, user)
	tx.PropsMu.Lock()
	defer tx.PropsMu.Unlock()

	if tx.Authenticating {
		return dbusutil.ToError(errors.New("transaction busy"))
	}

	err := tx.core.SetItemStr(pam.User, user)

	if err != nil {
		logger.Warning(err)
		return dbusutil.ToError(err)
	}

	tx.setUser(user)
	return nil
}
