package uadpagent

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	dbus "github.com/godbus/dbus"
	uadp "github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.uadp"
	secrets "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.secrets"
	"pkg.deepin.io/dde/daemon/session/common"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/procfs"
)

var logger = log.NewLogger("daemon/session/UadpAgent")

const (
	dbusServiceName = "com.deepin.daemon.UadpAgent"
	dbusPath        = "/com/deepin/daemon/UadpAgent"
	dbusInterface   = dbusServiceName
)

const (
	keyringTagExePath = "executablePath"
)

func (*UadpAgent) GetInterfaceName() string {
	return dbusInterface
}

type UadpAgent struct {
	service             *dbusutil.Service
	secretService       *secrets.Service
	secretSessionPath   dbus.ObjectPath
	defaultCollection   *secrets.Collection
	defaultCollectionMu sync.Mutex
	uadpDaemon          *uadp.Uadp // 提供加解密接口
	mu                  sync.Mutex

	secretData map[string]string // 密钥缓存

	methods *struct {
		SetDataKey func() `in:"keyName, dataKey"`
		GetDataKey func() `in:"keyName" out:"dataKey"`
	}
}

func newUadpAgent(service *dbusutil.Service) (*UadpAgent, error) {
	sessionBus := service.Conn()
	secretsObj := secrets.NewService(sessionBus)
	_, sessionPath, err := secretsObj.OpenSession(0, "plain", dbus.MakeVariant(""))
	if err != nil {
		logger.Warning("failed to get sessionPath:", err)
		return nil, err
	}
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	uadpDaemon := uadp.NewUadp(sysBus)
	u := &UadpAgent{
		service:           service,
		secretData:        make(map[string]string),
		secretService:     secretsObj,
		secretSessionPath: sessionPath,
		uadpDaemon:        uadpDaemon,
	}
	err = common.ActivateSysDaemonService(u.uadpDaemon.ServiceName_())
	if err != nil {
		logger.Warning(err)
	}
	return u, nil

}

// 提供给应用调用,用户通过此接口存储密钥
func (u *UadpAgent) SetDataKey(sender dbus.Sender, keyName string, dataKey string) *dbus.Error {
	executablePath, keyringKey, err := u.getExePathAndKeyringKey(sender, true)
	if err != nil {
		logger.Warning("failed to get exePath and keyringKey:", err)
		return dbusutil.ToError(err)
	}
	// 将用户希望存储的密钥加密存储
	err = u.uadpDaemon.SetDataKey(0, executablePath, keyName, dataKey, keyringKey)
	if err != nil {
		logger.Warning("failed to save data key:", err)
		return dbusutil.ToError(err)
	}
	return nil
}

// 提供给用户调用,用户通过此接口获取密钥
func (u *UadpAgent) GetDataKey(sender dbus.Sender, keyName string) (string, *dbus.Error) {
	executablePath, keyringKey, err := u.getExePathAndKeyringKey(sender, false)
	if err != nil {
		logger.Warning("failed to get exePath and keyringKey:", err)
		return "", dbusutil.ToError(err)
	}
	// 将密钥解密,并返回给用户
	key, err := u.uadpDaemon.GetDataKey(0, executablePath, keyName, keyringKey)
	if err != nil {
		logger.Warning("failed to get data key:", err)
		return "", dbusutil.ToError(err)
	}
	return key, nil
}

// 获取调用者二进制可执行文件路径和用于加解密数据的密钥
func (u *UadpAgent) getExePathAndKeyringKey(sender dbus.Sender, createIfNotExist bool) (string, string, error) {
	executablePath, err := u.getExePath(sender)
	if err != nil {
		logger.Warning("failed to get exePath:", err)
		return "", "", err
	}

	keyringKey, err := u.getKeyringKey(executablePath, createIfNotExist)
	if err != nil {
		logger.Warning("failed to get keyringKey:", err)
		return "", "", err
	}
	return executablePath, keyringKey, nil
}

func (u *UadpAgent) getKeyringKey(exePath string, createIfNotExist bool) (string, error) {
	var keyringKey string
	var err error
	keyringKey = u.secretData[exePath]

	if keyringKey == "" {
		keyringKey, err = u.getExeKeyringKey(exePath)
		if err != nil {
			logger.Warning(err)
		}
		if keyringKey == "" {
			if createIfNotExist {
				keyringKey, err = getRandKey(16)
				if err != nil {
					logger.Warning(err)
					return "", err
				}
				err = u.saveExeKeyringKey(exePath, keyringKey)
				if err != nil {
					logger.Warning("failed to save secret to keyring:", err)
				}
				u.mu.Lock()
				u.secretData[exePath] = keyringKey
				u.mu.Unlock()
			} else {
				return "", errors.New("unable to retrieve the password, please store the password first")
			}
		}
	}
	return keyringKey, nil
}

func (u *UadpAgent) getExePath(sender dbus.Sender) (string, error) {
	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		logger.Warning("failed to get PID:", err)
		return "", err
	}

	process := procfs.Process(pid)
	executablePath, err := process.Exe()
	if err != nil {
		logger.Warning("failed to get executablePath:", err)
		return "", err
	}
	return executablePath, nil
}

func (u *UadpAgent) saveExeKeyringKey(exePath, keyringKey string) error {
	label := fmt.Sprintf("UadpAgent code/decode secret for %s", exePath)
	logger.Debugf("set label: %q, exePath: %q, keyringKey: %q", label, exePath, keyringKey)
	itemSecret := secrets.Secret{
		Session:     u.secretSessionPath,
		Value:       []byte(keyringKey),
		ContentType: "text/plain",
	}

	properties := map[string]dbus.Variant{
		"org.freedesktop.Secret.Item.Label": dbus.MakeVariant(label),
		"org.freedesktop.Secret.Item.Type":  dbus.MakeVariant("org.freedesktop.Secret.Generic"),
		"org.freedesktop.Secret.Item.Attributes": dbus.MakeVariant(map[string]string{
			keyringTagExePath: exePath,
		}),
	}

	defaultCollection, err := u.getDefaultCollection()
	if err != nil {
		logger.Warning("failed to get defaultCollection:", err)
		return err
	}

	_, _, err = defaultCollection.CreateItem(0, properties, itemSecret, true)

	return err
}

func (u *UadpAgent) getExeKeyringKey(exePath string) (string, error) {
	attributes := map[string]string{
		keyringTagExePath: exePath,
	}

	defaultCollection, err := u.getDefaultCollection()
	if err != nil {
		logger.Warning("failed to get defaultCollection:", err)
		return "", err
	}
	items, err := defaultCollection.SearchItems(0, attributes)
	if err != nil {
		logger.Warning("failed to get items:", err)
		return "", err
	}

	secretData, err := u.secretService.GetSecrets(0, items, u.secretSessionPath)
	if err != nil {
		logger.Warning("failed to get secretData:", err)
		return "", err
	}
	var keyringKey string
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		logger.Warning(err)
		return "", err
	}
	for itemPath, secret := range secretData {
		itemObj, err := secrets.NewItem(sessionBus, itemPath)
		if err != nil {
			logger.Warning(err)
			return "", err
		}
		attributes, _ := itemObj.Attributes().Get(0)
		if attributes[keyringTagExePath] != "" {
			keyringKey = string(secret.Value)
		}
	}
	return keyringKey, nil
}

func (u *UadpAgent) getDefaultCollection() (*secrets.Collection, error) {
	u.defaultCollectionMu.Lock()
	defer u.defaultCollectionMu.Unlock()

	if u.defaultCollection != nil {
		return u.defaultCollection, nil
	}

	cPath, err := u.secretService.ReadAlias(0, "default")
	if err != nil {
		logger.Warning("failed to get collectionPath:", err)
		return nil, err
	}

	if cPath == "/" {
		return nil, errors.New("failed to get default collection path")
	}

	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	collectionObj, err := secrets.NewCollection(sessionBus, cPath)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}
	u.defaultCollection = collectionObj
	return collectionObj, nil
}

func getRandKey(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		logger.Warning(err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
