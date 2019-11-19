package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusTxInterface  = dbusInterface + ".Transaction"
	dbusTxPathPrefix = dbusPath + "/Transaction"
)

func getTxObjPath(id uint64) dbus.ObjectPath {
	return dbus.ObjectPath(dbusTxPathPrefix + strconv.FormatUint(id, 10))
}

type Transaction interface {
	getUserCookie() (string, string)
	clearCookie()
	matchSender(name string) bool
	getId() uint64

	GetInterfaceName() string
	Authenticate(sender dbus.Sender) *dbus.Error
	SetUser(sender dbus.Sender, user string) *dbus.Error
	End(sender dbus.Sender) *dbus.Error
}

var _ Transaction = &PAMTransaction{}
var _ Transaction = &FPrintTransaction{}

type baseTransaction struct {
	authType string
	parent   *Authority
	id       uint64
	agent    dbus.BusObject
	user     string
	cookie   string
	mu       sync.Mutex
}

func (tx *baseTransaction) String() string {
	return fmt.Sprintf("tx(%s,%d)", tx.authType, tx.id)
}

func (*baseTransaction) GetInterfaceName() string {
	return dbusTxInterface
}

func (tx *baseTransaction) getId() uint64 {
	return tx.id
}

func (tx *baseTransaction) checkSender(sender dbus.Sender) *dbus.Error {
	if tx.agent.Destination() != string(sender) {
		return dbusutil.ToError(errors.New("sender not match"))
	}
	return nil
}

func (tx *baseTransaction) matchSender(name string) bool {
	return tx.agent.Destination() == name
}

func (tx *baseTransaction) requestEchoOn(msg string) (ret string, err error) {
	logger.Debug(tx, "RequestEchoOn:", msg)
	err = tx.agent.Call(dbusAgentInterface+".RequestEchoOn", 0, msg).Store(&ret)
	return
}

func (tx *baseTransaction) requestEchoOff(msg string) (ret string, err error) {
	logger.Debug(tx, "RequestEchoOff:", msg)
	err = tx.agent.Call(dbusAgentInterface+".RequestEchoOff", 0, msg).Store(&ret)
	return
}

func (tx *baseTransaction) displayErrorMsg(msg string) error {
	logger.Debug(tx, "DisplayErrorMsg:", msg)
	return tx.agent.Call(dbusAgentInterface+".DisplayErrorMsg", 0, msg).Err
}

func (tx *baseTransaction) displayTextInfo(msg string) error {
	logger.Debug(tx, "DisplayTextInfo:", msg)
	return tx.agent.Call(dbusAgentInterface+".DisplayTextInfo", 0, msg).Err
}

func (tx *baseTransaction) sendResult(success bool) {
	var cookie string
	var err error
	if success {
		cookie, err = genCookie()
		if err != nil {
			logger.Warning(tx, "failed to gen cookie:", err)
		} else {
			tx.setCookie(cookie)
		}
	}

	err = tx.agent.Call(dbusAgentInterface+".RespondResult", 0,
		cookie).Err
	if err != nil {
		logger.Warning(tx, err)
	}
}

func (tx *baseTransaction) getUserCookie() (string, string) {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	return tx.user, tx.cookie
}

func (tx *baseTransaction) setCookie(cookie string) {
	tx.mu.Lock()
	tx.cookie = cookie
	tx.mu.Unlock()
}

func (tx *baseTransaction) clearCookie() {
	tx.mu.Lock()
	tx.cookie = ""
	tx.mu.Unlock()
}

func (tx *baseTransaction) setUser(user string) {
	tx.mu.Lock()
	if tx.user != user {
		tx.user = user
		tx.cookie = ""
	}
	tx.mu.Unlock()
}

func (tx *baseTransaction) getUser() string {
	tx.mu.Lock()
	user := tx.user
	tx.mu.Unlock()
	return user
}

func (tx *baseTransaction) getUserLocale() string {
	locale, err := tx.parent.getUserLocale(tx.getUser())
	if err != nil {
		logger.Warning(tx, "failed to get user locale:", err)
		return "en_US.UTF-8"
	}
	return locale
}
