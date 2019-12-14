package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
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

	core *pam.Transaction
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
		err := tx.displayErrorMsg("", msg)
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
	_, err := rand.Read(buf)
	// NOTE: err == nil only if we read len(buf) bytes.
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
	tx.setPropAuthenticating(false)
	tx.PropsMu.Unlock()

	if tx.hasEnded() {
		tx.terminate()
		return errTxEnd
	}
	return err
}

func (tx *PAMTransaction) Authenticate(sender dbus.Sender) *dbus.Error {
	tx.parent.service.DelayAutoQuit()
	if err := tx.checkSender(sender); err != nil {
		return err
	}

	logger.Debugf("%s Authenticate sender: %q", tx, sender)
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
	err := tx.checkSender(sender)
	if err != nil {
		return err
	}

	logger.Debugf("%s End sender: %s", tx, sender)
	if tx.hasEnded() {
		logger.Warningf("%s End sender: %s, tx has ended", tx, sender)
		return dbusutil.ToError(errTxEnd)
	}

	tx.clearSecret()
	tx.markEnd()

	tx.PropsMu.Lock()
	inAuth := tx.Authenticating
	tx.PropsMu.Unlock()

	if !inAuth {
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
