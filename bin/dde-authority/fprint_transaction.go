package main

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	dbus "github.com/godbus/dbus"
	"github.com/gosexy/gettext"
	fprint "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.fprintd"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"pkg.deepin.io/lib/strv"
)

var (
	msgVerificationFailed   = Tr("Fingerprint verification failed")
	msgVerificationTimedOut = Tr("Fingerprint verification timed out")
)

func Tr(str string) string {
	return str
}

type FPrintTransaction struct {
	baseTransaction
	PropsMu        sync.RWMutex
	Authenticating bool
	Sender         string
	methods        *struct { //nolint
		SetUser func() `in:"user"`
	}
	quit       chan struct{}
	release    chan struct{}
	devicePath dbus.ObjectPath
}

func (tx *FPrintTransaction) setPropAuthenticating(value bool) {
	if tx.Authenticating != value {
		tx.Authenticating = value
		err := tx.parent.service.EmitPropertyChanged(tx, "Authenticating", value)
		if err != nil {
			logger.Warning(tx, err)
		}
	}
}

func (tx *FPrintTransaction) getUser() string {
	user := tx.baseTransaction.getUser()
	if user != "" {
		return user
	}

	result, err := tx.requestEchoOn("login:")
	if err != nil {
		logger.Warning(tx, "requestEchoOn return err:", err)
		return ""
	} else {
		logger.Debug(tx, "requestEchoOn(\"login:\") result:", result)
		tx.setUser(result)
		return result
	}
}

func (tx *FPrintTransaction) getDevice() (*fprint.Device, error) {
	devicePath, err := tx.parent.fprintManager.GetDefaultDevice(0)
	if err != nil {
		return nil, err
	}

	deviceObj, err := fprint.NewDevice(tx.parent.service.Conn(), devicePath)
	if err != nil {
		return nil, err
	}

	return deviceObj, nil
}

func (tx *FPrintTransaction) setDevicePath(devPath dbus.ObjectPath) {
	tx.mu.Lock()
	tx.devicePath = devPath
	tx.mu.Unlock()
}

func (tx *FPrintTransaction) isClaimOk() bool {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	return tx.devicePath != ""
}

func (tx *FPrintTransaction) authenticate() error {
	tx.PropsMu.Lock()
	inAuth := tx.Authenticating
	tx.PropsMu.Unlock()

	if inAuth {
		return errors.New("transaction busy")
	}

	user := tx.getUser()
	if user == "" {
		return errors.New("user empty")
	}

	deviceObj, err := tx.getDevice()
	if err != nil {
		return err
	}

	fingers, err := deviceObj.ListEnrolledFingers(0, user)
	if err != nil {
		return err
	}

	if len(fingers) == 0 {
		return errors.New("user do not have fingerprint")
	}

	scanType, err := deviceObj.ScanType().Get(0)
	if err != nil {
		logger.Warning(err)
	}

	err = tx.claimDevice(deviceObj, user)
	if err != nil {
		return err
	}

	go func() {
		var err error
		tx.PropsMu.Lock()
		tx.setPropAuthenticating(true)
		tx.PropsMu.Unlock()
		tx.quit = make(chan struct{})
		tx.release = make(chan struct{})

		verifyErr := tx.verify(deviceObj, user, scanType)
		if verifyErr != nil {
			logger.Warning(tx, verifyErr)
		}

		if tx.isClaimOk() {
			logger.Debug(tx, "release device")
			tx.setDevicePath("")
			err = deviceObj.Release(0)
			if err != nil {
				logger.Warning(tx, err)
			}
		} else {
			logger.Debug(tx, "claim lost, do not call Release")
		}

		close(tx.release)

		tx.PropsMu.Lock()
		tx.setPropAuthenticating(false)
		tx.PropsMu.Unlock()

		tx.sendResult(verifyErr == nil)
	}()

	return nil
}

func (tx *FPrintTransaction) Authenticate(sender dbus.Sender) *dbus.Error {
	tx.parent.service.DelayAutoQuit()
	if err := tx.checkSender(sender); err != nil {
		return err
	}
	logger.Debugf("%s Authenticate sender: %q", tx, sender)
	err := tx.authenticate()
	if err != nil {
		logger.Infof("%s failed to authenticate: %v", tx, err)
	}
	return dbusutil.ToError(err)
}

type verifyResult struct {
	status string
	done   bool
}

func (tx *FPrintTransaction) releaseOtherTx(devPath dbus.ObjectPath) {
	tx.parent.releaseFprintTransaction(tx.id, devPath)
}

func (tx *FPrintTransaction) claimDevice(deviceObj *fprint.Device, user string) error {
	logger.Debug(tx, "claim device")

	caps, err := deviceObj.GetCapabilities(0)
	if err != nil {
		logger.Warning(err)
	}

	if strv.Strv(caps).Contains("ClaimForce") {
		tx.releaseOtherTx(deviceObj.Path_())
		// huawei device
		err = deviceObj.ClaimForce(0, user)
		if err != nil {
			return err
		}
	} else {
		// fprintd device
		err = deviceObj.Claim(0, user)
		if err != nil {
			if strings.Contains(err.Error(), "Could not attempt device open") {
				killFPrintDaemon()
			}
			return err
		}
	}

	tx.setDevicePath(deviceObj.Path_())
	return nil
}

var errClaimLost = errors.New("claim lost")

func (tx *FPrintTransaction) verify(deviceObj *fprint.Device, user, scanType string) error {
	if !tx.isClaimOk() {
		return errClaimLost
	}

	logger.Debugf("%v verify device: %v, user: %s", tx, deviceObj.Path_(), user)
	var isSwipe bool
	if scanType == "swipe" {
		isSwipe = true
	}

	locale := tx.getUserLocale()

	deviceObj.InitSignalExt(tx.parent.sigLoop, true)

	var msg string
	if isSwipe {
		msg = "Swipe your finger across the fingerprint reader"
	} else {
		msg = "Place your finger on the fingerprint reader"
	}
	msg = getFprintMsg(locale, msg)
	err := tx.displayTextInfo(msg)
	if err != nil {
		logger.Warning(tx, err)
	}

	verifyResultCh := make(chan verifyResult)

	_, err = deviceObj.ConnectVerifyStatus(func(result string, done bool) {
		logger.Debug(tx, "signal VerifyStatus", result, done)
		if tx.hasEnded() {
			return
		}

		msg := verifyResultStrToMsg(result, isSwipe)
		msg = getFprintMsg(locale, msg)
		if msg != "" {
			err := tx.displayErrorMsg(result, msg)
			if err != nil {
				logger.Warning(tx, err)
			}
		}

		verifyResultCh <- verifyResult{result, done}
	})
	if err != nil {
		logger.Warning(err)
	}

	var (
		verifyOk  bool
		continue0 bool
		maxTries  = 5
	)

	for maxTries > 0 {
		verifyOk, continue0 = tx.doVerify(deviceObj, verifyResultCh, locale)
		logger.Debugf("%v doVerify verifyOk: %v, continue: %v", tx, verifyOk, continue0)
		if !continue0 {
			break
		}
		maxTries--
	}

	deviceObj.RemoveHandler(proxy.RemoveAllHandlers)

	if verifyOk {
		return nil
	}
	return errors.New("verify failed")
}

func shouldLimitVerifyTime(deviceObj *fprint.Device) bool {
	path := string(deviceObj.Path_())
	return filepath.Base(path) != "huawei"
}

func (tx *FPrintTransaction) doVerify(deviceObj *fprint.Device, verifyResultCh chan verifyResult,
	locale string) (ok, continue0 bool) {

	if !tx.isClaimOk() {
		return
	}

	logger.Debug(tx, "VerifyStart")
	err := deviceObj.VerifyStart(0, "any")
	if err != nil {
		logger.Warning(err)
		return false, false
	}

	var timeCh <-chan time.Time
	if shouldLimitVerifyTime(deviceObj) {
		timeCh = time.After(10 * time.Second)
	}

	var status string
loop:
	for {
		select {
		case <-timeCh:
			logger.Warning(tx, "timed out")
			msg := getFprintMsg(locale, msgVerificationTimedOut)
			err := tx.displayErrorMsg("verify-timed-out", msg)
			if err != nil {
				logger.Warning(err)
			}
			break loop

		case <-tx.quit:
			logger.Debug(tx, "receive quit")
			break loop

		case result, ok := <-verifyResultCh:
			if ok {
				logger.Debug(tx, "receive result", result)
				if result.done {
					status = result.status
					break loop
				}
			} else {
				logger.Debug("verifyResultCh closed")
				break loop
			}
		}
	}

	if tx.isClaimOk() {
		logger.Debug(tx, "VerifyStop")
		stopCallCh := make(chan *dbus.Call, 1)
		deviceObj.GoVerifyStop(0, stopCallCh)
		select {
		case call := <-stopCallCh:
			logger.Debug(tx, "VerifyStop done")
			if call.Err != nil {
				logger.Warning(err)
			}
		case <-time.After(10 * time.Second):
			logger.Warning(tx, "VerifyStop timed out")
		}
		switch status {
		case "verify-no-match":
			continue0 = true
		case "verify-match":
			ok = true
		}
	} else {
		logger.Debug(tx, "claim lost, do not call VerifyStop")
		// verify failed, no continue
	}

	return
}

func verifyResultStrToMsg(result string, isSwipe bool) string {
	switch result {
	case "verify-no-match":
		return msgVerificationFailed
	case "verify-retry-scan":
		if isSwipe {
			return "Swipe your finger again"
		} else {
			return "Place your finger on the reader again"
		}
	case "verify-swipe-too-short":
		return "Swipe was too short, try again"
	case "verify-finger-not-centered":
		return "Your finger was not centered, try swiping your finger again"
	case "verify-remove-and-retry":
		return "Remove your finger, and try swiping your finger again"
	default:
		return ""
	}
}

func getFprintMsg(locale, msgId string) string {
	gettext.SetLocale(gettext.LC_ALL, locale)
	domain := "fprintd"
	switch msgId {
	case "":
		return ""
	case msgVerificationFailed, msgVerificationTimedOut:
		domain = "dde-daemon"
	}
	text := gettext.DGettext(domain, msgId)
	if utf8.ValidString(text) {
		return text
	}
	return msgId
}

func (tx *FPrintTransaction) SetUser(sender dbus.Sender, user string) *dbus.Error {
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

	tx.setUser(user)
	return nil
}

func (tx *FPrintTransaction) End(sender dbus.Sender) *dbus.Error {
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

	if inAuth {
		logger.Debug(tx, "force quit")
		close(tx.quit)
		<-tx.release // wait tx.release chan closed
	}

	tx.parent.deleteTx(tx.id)
	return nil
}

func killFPrintDaemon() {
	logger.Debug("kill fprintd")
	err := exec.Command("pkill", "-f", "/usr/lib/fprintd/fprintd").Run()
	if err != nil {
		logger.Warning("failed to kill fprintd:", err)
	}
}
