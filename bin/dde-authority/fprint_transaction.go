package main

import (
	"errors"
	"os/exec"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gosexy/gettext"
	fprint "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.fprintd"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
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
	methods        *struct {
		SetUser func() `in:"user"`
	}
	quit    chan struct{}
	release chan struct{}
}

func (tx *FPrintTransaction) setPropAuthenticating(value bool) {
	if tx.Authenticating != value {
		tx.Authenticating = value
		err := tx.parent.service.EmitPropertyChanged(tx, "Authenticating", value)
		if err != nil {
			logger.Warning(err)
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

func (tx *FPrintTransaction) Authenticate(sender dbus.Sender) *dbus.Error {
	tx.parent.service.DelayAutoQuit()
	if err := tx.checkSender(sender); err != nil {
		return err
	}

	tx.PropsMu.Lock()
	defer tx.PropsMu.Unlock()

	if tx.Authenticating {
		return dbusutil.ToError(errors.New("transaction busy"))
	}

	user := tx.getUser()
	if user == "" {
		return dbusutil.ToError(errors.New("user empty"))
	}

	deviceObj, err := tx.getDevice()
	if err != nil {
		return dbusutil.ToError(err)
	}
	go func() {
		tx.PropsMu.Lock()
		tx.setPropAuthenticating(true)
		tx.quit = make(chan struct{})
		tx.release = make(chan struct{})
		tx.PropsMu.Unlock()

		err := tx.authenticate(deviceObj, user)

		tx.PropsMu.Lock()
		tx.setPropAuthenticating(false)
		tx.quit = nil
		tx.release = nil
		tx.PropsMu.Unlock()

		tx.sendResult(err == nil)
		if err != nil {
			logger.Warning(err)
		}
	}()

	return nil
}

type verifyResult struct {
	status string
	done   bool
}

func claimFprintDevice(deviceObj *fprint.Device, user string) error {
	logger.Debug("claim device")
	err := deviceObj.Claim(0, user)
	if err != nil {
		if strings.Contains(err.Error(), "Could not attempt device open") {
			killFPrintDaemon()
		}
		return err
	}
	return nil

	//caps, err := deviceObj.GetCapabilities(0)
	//if err != nil {
	//	logger.Warning(err)
	//}
	//
	//if strv.Strv(caps).Contains("ClaimForce") {
	//	err = deviceObj.ClaimForce(0, user)
	//	if err != nil {
	//		return err
	//	}
	//
	//} else {
	//	err = deviceObj.Claim(0, user)
	//	if err != nil {
	//		if strings.Contains(err.Error(), "Could not attempt device open") {
	//			killFPrintDaemon()
	//		}
	//		return err
	//	}
	//}
}

func (tx *FPrintTransaction) authenticate(deviceObj *fprint.Device, user string) error {
	logger.Debugf("%v authenticate device: %v, user: %s", tx, deviceObj.Path_(), user)
	err := claimFprintDevice(deviceObj, user)
	if err != nil {
		return err
	}

	scanType, err := deviceObj.ScanType().Get(0)
	if err != nil {
		logger.Warning(err)
	}
	var isSwipe bool
	if scanType == "swipe" {
		isSwipe = true
	}

	locale := tx.getUserLocale()

	deviceObj.InitSignalExt(tx.parent.sigLoop, true)
	_, err = deviceObj.ConnectVerifyFingerSelected(func(finger string) {
		var msg string
		if isSwipe {
			msg = "Swipe your finger across the fingerprint reader"
		} else {
			msg = "Place your finger on the fingerprint reader"
		}
		msg = getFprintMsg(locale, msg)
		err := tx.displayTextInfo(msg)
		if err != nil {
			logger.Warning(err)
		}
	})
	if err != nil {
		logger.Warning(err)
	}

	verifyResultCh := make(chan verifyResult)

	_, err = deviceObj.ConnectVerifyStatus(func(result string, done bool) {
		logger.Debug(tx, "signal VerifyStatus", result, done)

		msg := verifyResultStrToMsg(result, isSwipe)
		msg = getFprintMsg(locale, msg)
		if msg != "" {
			err := tx.displayErrorMsg(msg)
			if err != nil {
				logger.Warning(err)
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
		maxTries  = 3
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
	close(verifyResultCh)

	logger.Debug(tx, "release device")
	err = deviceObj.Release(0)
	if err != nil {
		logger.Warning(err)
	}
	close(tx.release)

	if verifyOk {
		return nil
	}
	return errors.New("verify failed")
}

func (tx *FPrintTransaction) doVerify(deviceObj *fprint.Device, verifyResultCh chan verifyResult,
	locale string) (ok, continue0 bool) {

	logger.Debug(tx, "VerifyStart")
	err := deviceObj.VerifyStart(0, "any")
	if err != nil {
		logger.Warning(err)
		return false, false
	}

	var status string
loop:
	for {
		select {
		case <-time.After(10 * time.Second):
			logger.Warning(tx, "timed out")
			msg := getFprintMsg(locale, msgVerificationTimedOut)
			err := tx.displayErrorMsg(msg)
			if err != nil {
				logger.Warning(err)
			}
			break loop

		case <-tx.quit:
			logger.Debug(tx, "receive quit")
			break loop

		case result := <-verifyResultCh:
			logger.Debug(tx, "receive result", result)
			if result.done {
				status = result.status
				break loop
			}
		}
	}

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
	if err := tx.checkSender(sender); err != nil {
		return err
	}
	logger.Debug(tx, "end")
	tx.clearCookie()
	tx.PropsMu.Lock()
	if tx.Authenticating {
		logger.Debug(tx, "force quit")
		close(tx.quit)
		<-tx.release // wait tx.release chan closed
	}
	tx.PropsMu.Unlock()
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
