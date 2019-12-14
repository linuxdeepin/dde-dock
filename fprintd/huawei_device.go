package fprintd

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	accounts "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.accounts"
	huawei_fprint "github.com/linuxdeepin/go-dbus-factory/com.huawei.fingerprint"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/strv"
)

const (
	huaweiDevicePath = dbusPath + "/Device/huawei"
	huaweiFprintDir  = "/var/lib/dde-daemon/fingerprint/huawei"

	huaweiDeleteTypeOne = 0
	huaweiDeleteTypeAll = 1
)

type HuaweiDevice struct {
	service *dbusutil.Service
	core    *huawei_fprint.Fingerprint

	mu       sync.Mutex
	claimed  bool
	sender   string
	username string
	userUuid string

	ScanType string // const

	methods *deviceMethods
	signals *deviceSignals
}

func (d *HuaweiDevice) destroy() {
}

func (d *HuaweiDevice) getCorePath() dbus.ObjectPath {
	return huaweiDevicePath
}

func (d *HuaweiDevice) getPath() dbus.ObjectPath {
	return huaweiDevicePath
}

const (
	huaweiDeviceStatusBusy = 1
	huaweiDeviceStatusIdle = 0
)

func getUserUuid(username string) (string, error) {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return "", err
	}
	accountsObj := accounts.NewAccounts(sysBus)
	userPath, err := accountsObj.FindUserByName(0, username)
	if err != nil {
		return "", err
	}
	userObj, err := accounts.NewUser(sysBus, dbus.ObjectPath(userPath))
	if err != nil {
		return "", err
	}
	uuid, err := userObj.UUID().Get(0)
	if err != nil {
		return "", err
	}
	if uuid == "" {
		return "", errors.New("get empty uuid")
	}

	return uuid, nil
}

func (dev *HuaweiDevice) isFree() (bool, error) {
	dev.mu.Lock()
	defer dev.mu.Unlock()
	return !dev.claimed, nil
}

func (dev *HuaweiDevice) claim(sender, username string) error {
	dev.mu.Lock()
	defer dev.mu.Unlock()

	if dev.claimed {
		return errors.New("device is claimed by user " + dev.username)
	}

	userUuid, err := getUserUuid(username)
	if err != nil {
		return err
	}

	status, err := dev.core.GetStatus(0)
	if err != nil {
		return err
	}

	if status == huaweiDeviceStatusBusy {
		return errors.New("device is busy")
	}

	dev.claimed = true
	dev.sender = sender
	dev.username = username
	dev.userUuid = userUuid
	return nil
}

func (dev *HuaweiDevice) claimForce(sender, username string) error {
	dev.mu.Lock()
	if dev.claimed {
		err := dev.close()
		if err != nil {
			dev.mu.Unlock()
			return err
		}
		dev.releaseAux()
	}
	dev.mu.Unlock()

	return dev.claim(sender, username)
}

func (dev *HuaweiDevice) close() error {
	status, err := dev.core.GetStatus(0)
	if err != nil {
		return err
	}

	if status == huaweiDeviceStatusBusy {
		closeRet, err := dev.core.Close(0)
		if err != nil {
			return err
		}

		if closeRet == -1 {
			return errors.New("failed to close")
		}
	} // else status is idle, no need call close
	return nil
}

func (dev *HuaweiDevice) Claim(sender dbus.Sender, username string) *dbus.Error {
	err := dev.claim(string(sender), username)
	if err != nil {
		logger.Debugf("claim() sender: %q, username: %q, err %v", sender, username, err)
	} else {
		logger.Debugf("claim() sender: %q, username: %q, ok", sender, username)
	}
	return dbusutil.ToError(err)
}

func (dev *HuaweiDevice) ClaimForce(sender dbus.Sender, username string) *dbus.Error {
	err := dev.claimForce(string(sender), username)
	if err != nil {
		logger.Debugf("claimForce() sender: %q, username: %q, err %v", sender, username, err)
	} else {
		logger.Debugf("claimForce() sender: %q, username: %q, ok", sender, username)
	}
	return dbusutil.ToError(err)
}

func (dev *HuaweiDevice) GetCapabilities() ([]string, *dbus.Error) {
	return []string{"ClaimForce", "DeleteEnrolledFinger"}, nil
}

func (dev *HuaweiDevice) Release(sender dbus.Sender) *dbus.Error {
	err := dev.release(string(sender))
	if err != nil {
		logger.Debugf("release() sender: %q, err: %v", sender, err)
	} else {
		logger.Debugf("release() sender: %q, ok", sender)
	}
	return dbusutil.ToError(err)
}

func (dev *HuaweiDevice) releaseAux() {
	dev.claimed = false
	dev.sender = ""
	dev.username = ""
	dev.userUuid = ""
}

func (dev *HuaweiDevice) release(sender string) error {
	dev.mu.Lock()
	defer dev.mu.Unlock()

	if !dev.claimed {
		return errors.New("device is not claimed")
	}

	if dev.sender != sender {
		return errors.New("sender not match")
	}

	dev.releaseAux()
	return nil
}

func (dev *HuaweiDevice) checkClaimed(sender dbus.Sender) (userUuid string, err error) {
	dev.mu.Lock()
	defer dev.mu.Unlock()

	if !dev.claimed {
		return "", errors.New("device is not claimed")
	}

	if dev.sender != string(sender) {
		return "", errors.New("sender not match")
	}

	return dev.userUuid, nil
}

var fprintdFingerprintNames = strv.Strv{
	"left-thumb",
	"left-index-finger",
	"left-middle-finger",
	"left-ring-finger",
	"left-little-finger",

	"right-thumb",
	"right-index-finger",
	"right-middle-finger",
	"right-ring-finger",
	"right-little-finger",
}

func (dev *HuaweiDevice) enrollStart(sender dbus.Sender, finger string) error {
	err := checkAuth(actionIdEnroll, string(sender))
	if err != nil {
		return err
	}

	if !fprintdFingerprintNames.Contains(finger) {
		return errors.New("invalid fingerprint name")
	}

	userUuid, err := dev.checkClaimed(sender)
	if err != nil {
		return err
	}

	dir, err := ensureHuaweiFprintDir(userUuid)
	if err != nil {
		return err
	}

	filename := filepath.Join(dir, finger)
	_, err = os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		} // else file not exist, pass
	} else {
		// file exist
		err = os.Remove(filename)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		reloadRet, err := dev.core.Reload(0, huaweiDeleteTypeOne)
		if err != nil {
			return err
		}
		if reloadRet == -1 {
			return errors.New("failed to load")
		}
	}

	err = dev.core.Enroll(dbus.FlagNoReplyExpected, filename, userUuid)
	return err
}

func ensureHuaweiFprintDir(userUuid string) (dir string, err error) {
	err = os.MkdirAll(huaweiFprintDir, 0755)
	if err != nil {
		return
	}
	dir = filepath.Join(huaweiFprintDir, userUuid)
	err = os.Mkdir(dir, 0700)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	return dir, nil
}

func (dev *HuaweiDevice) EnrollStart(sender dbus.Sender, finger string) *dbus.Error {
	err := dev.enrollStart(sender, finger)
	if err != nil {
		logger.Debugf("enrollStart() sender: %q, finger: %q, err: %v", sender, finger, err)
	} else {
		logger.Debugf("enrollStart() sender: %q, finger: %q, ok", sender, finger)
	}
	return dbusutil.ToError(err)
}

func (dev *HuaweiDevice) stop(sender dbus.Sender) error {
	_, err := dev.checkClaimed(sender)
	if err != nil {
		return err
	}
	return dev.close()
}

func (dev *HuaweiDevice) EnrollStop(sender dbus.Sender) *dbus.Error {
	err := dev.stop(sender)
	if err != nil {
		logger.Debugf("enrollStop() sender: %q, err: %v", sender, err)
	} else {
		logger.Debugf("enrollStop() sender: %q, ok", sender)
	}
	return dbusutil.ToError(err)
}

func (dev *HuaweiDevice) verifyStart(sender dbus.Sender) error {
	userUuid, err := dev.checkClaimed(sender)
	if err != nil {
		return err
	}

	err = dev.core.Identify(dbus.FlagNoReplyExpected, userUuid)
	return err
}

func (dev *HuaweiDevice) VerifyStart(sender dbus.Sender, _finger string) *dbus.Error {
	err := dev.verifyStart(sender)
	return dbusutil.ToError(err)
}

func (dev *HuaweiDevice) VerifyStop(sender dbus.Sender) *dbus.Error {
	err := dev.stop(sender)
	return dbusutil.ToError(err)
}

func (dev *HuaweiDevice) deleteEnrolledFingers(sender dbus.Sender, username string) error {
	err := checkAuth(actionIdDelete, string(sender))
	if err != nil {
		return err
	}

	userUuid, err := getUserUuid(username)
	if err != nil {
		return err
	}
	dir, err := ensureHuaweiFprintDir(userUuid)
	if err != nil {
		return err
	}

	fileInfoList, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fileInfo := range fileInfoList {
		filename := filepath.Join(dir, fileInfo.Name())
		err = os.Remove(filename)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	reloadRet, err := dev.core.Reload(0, huaweiDeleteTypeAll)
	if err != nil {
		return err
	}
	if reloadRet == -1 {
		return errors.New("failed to reload")
	}
	return nil
}

func (dev *HuaweiDevice) DeleteEnrolledFingers(sender dbus.Sender, username string) *dbus.Error {
	err := dev.deleteEnrolledFingers(sender, username)
	return dbusutil.ToError(err)
}

func (dev *HuaweiDevice) DeleteEnrolledFinger(sender dbus.Sender, username, finger string) *dbus.Error {
	err := dev.deleteEnrolledFinger(sender, username, finger)
	return dbusutil.ToError(err)
}

func (dev *HuaweiDevice) deleteEnrolledFinger(sender dbus.Sender, username, finger string) error {
	err := checkAuth(actionIdDelete, string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}

	userUuid, err := getUserUuid(username)
	if err != nil {
		return err
	}

	dir, err := ensureHuaweiFprintDir(userUuid)
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join(dir, finger))
	if err != nil {
		if os.IsNotExist(err) {
			return errors.New("not found finger")
		}
		return err
	}

	reloadRet, err := dev.core.Reload(0, huaweiDeleteTypeOne)
	if err != nil {
		return err
	}
	if reloadRet == -1 {
		return errors.New("failed to reload")
	}

	return nil
}

func (dev *HuaweiDevice) ListEnrolledFingers(username string) ([]string, *dbus.Error) {
	result, err := dev.listEnrolledFingers(username)
	if err != nil {
		logger.Warningf("ListEnrolledFingers() username: %q, err: %v", username, err)
	} else {
		logger.Debugf("ListEnrolledFingers() username: %q, ret: %v", username, result)
	}
	return result, dbusutil.ToError(err)
}

func (dev *HuaweiDevice) listEnrolledFingers(username string) ([]string, error) {
	userUuid, err := getUserUuid(username)
	if err != nil {
		return nil, err
	}

	dir, err := ensureHuaweiFprintDir(userUuid)
	if err != nil {
		return nil, err
	}

	fileInfoList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, fileInfo := range fileInfoList {
		if fileInfo.IsDir() {
			continue
		}

		result = append(result, fileInfo.Name())
	}

	return result, nil
}

func (*HuaweiDevice) GetInterfaceName() string {
	return dbusDeviceInterface
}

const (
	fprintdEnrollStatusCompleted          = "enroll-completed"
	fprintdEnrollStatusFailed             = "enroll-failed"
	fprintdEnrollStatusStagePassed        = "enroll-stage-passed"
	fprintdEnrollStatusRetryScan          = "enroll-retry-scan"
	fprintdEnrollStatusSwipeTooShort      = "enroll-swipe-too-short"
	fprintdEnrollStatusFingerNotCentered  = "enroll-finger-not-centered"
	fprintdEnrollStatusRemoveAndRetry     = "enroll-remove-and-retry"
	fprintdEnrollStatusEnrollDisconnected = "enroll-disconnected"
	fprintdEnrollStatusEnrollUnknownError = "enroll-unknown-error"

	fprintdVerifyStatusNoMatch      = "verify-no-match"
	fprintdVerifyStatusMatch        = "verify-match"
	fprintdVerifyStatusUnknownError = "verify-unknown-error"
)

func (dev *HuaweiDevice) handleSignalEnrollStatus(progress int32, result int32) {
	logger.Debug("signal EnrollStatus", progress, result)
	var done bool
	var status string
	switch result {
	case -2:
		// 没进行设备初始化就进行录入操作
		done = true
		status = fprintdEnrollStatusFailed
		logger.Debug("failed, no device initialization")
	case -1:
		// 指纹录入错误（指纹录入错误，结束指纹录入，多为函数的参数问题引发的错误）
		done = true
		status = fprintdEnrollStatusFailed
		logger.Debug("failed")
	case 1:
		if progress == 100 {
			// 指纹录入完成（指纹录入以及指纹模板保存完成，结束指纹录入）
			done = true
			status = fprintdEnrollStatusCompleted
			logger.Debug("completed")
		} else {
			logger.Warningf("ignore invalid signal EnrollStatus(%d,%d)", progress, result)
		}
	case 2:
		// 指纹录入失败（指纹录入失败，结束指纹录入，多为指纹设备异常出现的错误）
		done = true
		status = fprintdEnrollStatusFailed
		logger.Debug("failed")
	case 3:
		// 单张指纹图像采图完成
		status = fprintdEnrollStatusStagePassed
		logger.Debug("Single fingerprint image acquisition completed")

	case 4:
		// 当前手指指纹模板已存在，需换其他手指录入指纹
		status = fprintdEnrollStatusRetryScan
		logger.Debug("The current finger fingerprint template already exists. You need to change the fingerprint of other fingers.")

	case 104:
		// TODO
		status = fprintdEnrollStatusRetryScan
		logger.Warning("unknown enroll result", result)

	case 100, 105:
		// 指纹图像质量太差，或其他设备扫描的原因需要重新录入指纹
		status = fprintdEnrollStatusRetryScan
		logger.Debug("The fingerprint image quality is too bad, or the reason for other device scanning needs to re-enter the fingerprint")

	case 106:
		// 生成的指纹模板已重复（指纹采图结束后自动生成的指纹模板与存在的指纹模板重复，结束指纹录入）
		done = true
		status = fprintdEnrollStatusFailed
		logger.Debug("The generated fingerprint template has been duplicated")

	case 107:
		// 指纹向左移动
		status = fprintdEnrollStatusStagePassed
		logger.Debug("move left")
	case 108:
		// 指纹向下移动
		status = fprintdEnrollStatusStagePassed
		logger.Debug("move down")
	case 109:
		// 指纹向右移动
		status = fprintdEnrollStatusStagePassed
		logger.Debug("move right")
	case 110:
		// 指纹向上移动
		status = fprintdEnrollStatusStagePassed
		logger.Debug("move up")

	default:
		logger.Warning("unknown EnrollStatus result", result)
		return
	}

	// TODO
	//status = fmt.Sprintf("%s;%d;%d", status, progress, result)
	dev.emitSignalEnrollStatus(status, done)
}

func (dev *HuaweiDevice) emitSignalEnrollStatus(status string, done bool) {
	err := dev.service.Emit(dev, "EnrollStatus", status, done)
	if err != nil {
		logger.Warning(err)
	}
}

func (dev *HuaweiDevice) handleSignalIdentifyStatus(result int32) {
	logger.Debug("signal IdentifyStatus", result)
	var done bool
	var status string

	switch result {
	case -2:
		// 设备初始化失败
		done = true
		status = fprintdVerifyStatusUnknownError

	case -1:
		// 认证超时，以后会被废弃
		done = true
		status = fprintdVerifyStatusNoMatch

	case 0:
		// 认证成功
		done = true
		status = fprintdVerifyStatusMatch

	case 1:
		// 认证失败
		done = true
		status = fprintdVerifyStatusNoMatch

	default:
		logger.Warning("unknown IdentifyStatus result", result)
		return
	}

	dev.emitSignalVerifyStatus(status, done)
}

func (dev *HuaweiDevice) emitSignalVerifyStatus(status string, done bool) {
	err := dev.service.Emit(dev, "VerifyStatus", status, done)
	if err != nil {
		logger.Warning(err)
	}
}

func (dev *HuaweiDevice) handleNameLost(name string) {
	dev.mu.Lock()
	defer dev.mu.Unlock()

	if !dev.claimed {
		return
	}

	if dev.sender == name {
		logger.Debugf("name %s lost, auto release", name)
		dev.releaseAux()
		err := dev.close()
		if err != nil {
			logger.Warning(err)
		}
	}
}
