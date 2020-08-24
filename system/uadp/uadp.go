package uadp

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	dbus "github.com/godbus/dbus"
	polkit "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.policykit1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/procfs"
)

var logger = log.NewLogger("daemon/Uadp")

const (
	dbusServiceName = "com.deepin.daemon.Uadp"
	dbusPath        = "/com/deepin/daemon/Uadp"
	dbusInterface   = dbusServiceName
)

const (
	allowedProcess = "/usr/lib/deepin-daemon/dde-session-daemon"

	polkitActionUadp = "com.deepin.daemon.uadp.call"

	UadpDataDir = "/var/lib/dde-daemon/uadp"
)

func (*Uadp) GetInterfaceName() string {
	return dbusInterface
}

type Uadp struct {
	service    *dbusutil.Service
	appDataMap map[string]map[string][]byte // 应用加密数据缓存
	fileNames  map[string]string            // 文件索引缓存

	secretMu sync.Mutex
	mu       sync.Mutex

	methods *struct {
		SetDataKey func() `in:"exePath,keyName,dataKey,keyringKey"`
		GetDataKey func() `in:"exePath,keyName,keyringKey" out:"dataKey"`
	}
}

func newUadp(service *dbusutil.Service) (*Uadp, error) {
	u := &Uadp{
		service:    service,
		appDataMap: make(map[string]map[string][]byte),
		fileNames:  make(map[string]string),
	}
	return u, nil
}

// 加密用户存储的密钥并存储在文件中
func (u *Uadp) SetDataKey(sender dbus.Sender, exePath, keyName, dataKey, keyringKey string) *dbus.Error {
	_, err := u.verifyIdentity(sender)
	if err != nil {
		logger.Warning("failed to verify:", err)
		return dbusutil.ToError(err)
	}

	// 通过polkit，防止远程访问
	pass, err := u.checkAuth(string(sender), polkitActionUadp)
	if err != nil {
		logger.Warning("failed to pass authentication:", err)
		return dbusutil.ToError(err)
	}

	if !pass {
		return dbusutil.ToError(errors.New("not be authorized"))
	}

	err = u.setDataKey(exePath, keyName, dataKey, keyringKey)
	if err != nil {
		logger.Warning("failed to encrypt key:", err)
		return dbusutil.ToError(err)
	}
	return nil
}

func (u *Uadp) setDataKey(exePath, keyName, dataKey, keyringKey string) error {
	encryptedKey, err := aesEncryptKey(dataKey, keyringKey)
	if err != nil {
		logger.Warning("failed to encryptKey by aes:", err)
		return err
	}
	u.secretMu.Lock()
	if u.appDataMap[exePath] == nil {
		u.appDataMap[exePath] = make(map[string][]byte)
	}
	u.appDataMap[exePath][keyName] = encryptedKey
	u.secretMu.Unlock()

	err = u.updateDataFile(exePath)
	if err != nil {
		logger.Warning("failed to updateDataFile:", err)
		return err
	}
	return nil
}

func aesEncryptKey(origin, key string) ([]byte, error) {
	origData := []byte(origin)
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = pkcs7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	encrypted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(encrypted, origData)
	return encrypted, nil
}

func pkcs7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func (u *Uadp) GetDataKey(sender dbus.Sender, exePath, keyName, keyringKey string) (string, *dbus.Error) {
	_, err := u.verifyIdentity(sender)
	if err != nil {
		logger.Warning("failed to verify:", err)
		return "", dbusutil.ToError(err)
	}

	// 通过polkit授权防止被远程访问
	pass, err := u.checkAuth(string(sender), polkitActionUadp)
	if err != nil {
		logger.Warning("failed to pass authentication:", err)
		return "", dbusutil.ToError(err)
	}

	if !pass {
		return "", dbusutil.ToError(errors.New("not authorized"))
	}

	key, err := u.getDataKey(exePath, keyName, keyringKey)
	if err != nil {
		logger.Warning("failed to decrypt Key:", err)
		return "", dbusutil.ToError(err)
	}
	return key, nil
}

func (u *Uadp) getDataKey(exePath, keyName, keyringKey string) (string, error) {
	encryptedKey := u.findKeyFromCacheOrFile(exePath, keyName)
	if encryptedKey == nil {
		return "", errors.New("failed to find data used to be decrypted")
	}
	key, err := u.aesDecrypt(encryptedKey, keyringKey)
	if err != nil {
		logger.Warning("failed to aesDecrypt key:", err)
		return "", err
	}
	return key, nil
}

func (u *Uadp) aesDecrypt(encryptedKey []byte, key string) (string, error) {
	encryptedKeyByte := []byte(encryptedKey)
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		logger.Warning("failed to newCipher key:", err)
		return "", err
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	decrypted := make([]byte, len(encryptedKeyByte))
	// 解密
	blockMode.CryptBlocks(decrypted, encryptedKeyByte)
	// 去补全码
	decrypted = pkcs7UnPadding(decrypted)
	return string(decrypted), nil
}

func pkcs7UnPadding(originData []byte) []byte {
	length := len(originData)
	unpadding := int(originData[length-1])
	return originData[:(length - unpadding)]
}

func (u *Uadp) findKeyFromCacheOrFile(exePath, keyName string) []byte {
	if secretData, ok := u.appDataMap[exePath]; ok {
		if value, ok := secretData[keyName]; ok {
			return value
		}
	}

	secretData, err := u.loadDataFromFile(exePath)
	if err != nil {
		logger.Warning("failed to loadDataFromFile:", err)
		return nil
	}

	u.secretMu.Lock()
	defer u.secretMu.Unlock()

	u.appDataMap[exePath] = secretData

	return u.appDataMap[exePath][keyName]
}

func (u *Uadp) loadDataFromFile(exePath string) (map[string][]byte, error) {
	fileName, err := u.getFileName(exePath)
	if err != nil {
		logger.Warning("failed to get filename:", err)
		return nil, err
	}
	var secretData map[string][]byte
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Warning("cannot read data from file:", err)
		return nil, err
	}
	err = unmarshalGob(content, &secretData)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	return secretData, nil
}

func (u *Uadp) updateDataFile(exePath string) error {
	secretData := u.appDataMap[exePath]
	fileName, err := u.getFileName(exePath)
	if err != nil {
		logger.Warning("failed to get filename:", err)
		return err
	}

	newFileName := fileName + "-1"
	content, err := marshalGob(secretData)
	if err != nil {
		logger.Warning(err)
		return err
	}
	err = writeFile(newFileName, content, 0600)
	if err != nil {
		logger.Warning(err)
		return err
	}

	err = os.Rename(newFileName, fileName)
	if err != nil {
		logger.Warning(err)
		return err
	}
	return nil
}

func (u *Uadp) getFileName(exePath string) (string, error) {
	err := os.MkdirAll(UadpDataDir, 0755)
	if err != nil {
		logger.Warning(err)
		return "", err
	}
	var fileName string
	fileName = u.fileNames[exePath]

	var fileNames map[string]string

	if fileName == "" {
		fileMap := filepath.Join(UadpDataDir, "filemap")
		content, err := ioutil.ReadFile(fileMap)
		if err == nil {
			err = json.Unmarshal(content, &fileNames)
			if err != nil {
				logger.Warning(err)
				return "", err
			}
			u.mu.Lock()
			u.fileNames = fileNames
			fileName = u.fileNames[exePath]
			u.mu.Unlock()
		}

		if fileName == "" {
			// 新增文件索引
			fileName = fmt.Sprintf("%d", len(fileNames))
			u.mu.Lock()
			u.fileNames[exePath] = fileName
			u.mu.Unlock()
			content, err := json.Marshal(u.fileNames)
			if err != nil {
				logger.Warning(err)
				return "", err
			}
			err = writeFile(fileMap, content, 0600)
			if err != nil {
				logger.Warning(err)
				return "", err
			}
		}
	}

	file := filepath.Join(UadpDataDir, fileName)
	return file, nil
}

func (u *Uadp) verifyIdentity(sender dbus.Sender) (bool, error) {
	pid, err := u.service.GetConnPID(string(sender))
	if err != nil {
		logger.Warning("failed to get PID:", err)
		return false, err
	}

	process := procfs.Process(pid)
	executablePath, err := process.Exe()
	if err != nil {
		logger.Warning("failed to get executablePath:", err)
		return false, err
	}

	if executablePath == allowedProcess {
		return true, nil
	}

	return false, errors.New("process is not allowed to access")
}

func (u *Uadp) checkAuth(sysBusName, actionId string) (bool, error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return false, err
	}
	authority := polkit.NewAuthority(systemBus)
	subject := polkit.MakeSubject(polkit.SubjectKindSystemBusName)
	subject.SetDetail("name", sysBusName)

	ret, err := authority.CheckAuthorization(0, subject,
		actionId, nil,
		polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}
	return ret.IsAuthorized, nil
}

func writeFile(filename string, data []byte, perm os.FileMode) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		logger.Warning(err)
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	_, err = file.Write(data)
	if err != nil {
		logger.Warning(err)
		return err
	}

	err = file.Sync()
	if err != nil {
		logger.Warning(err)
		return err
	}
	return err
}

func unmarshalGob(content []byte, secretData interface{}) error {
	r := bytes.NewReader(content)
	dec := gob.NewDecoder(r)
	if err := dec.Decode(secretData); err != nil {
		logger.Warning("decode datas failed:", err)
		return err
	}
	return nil
}

func marshalGob(secretData map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(secretData)
	if err != nil {
		logger.Warning(err)
		return nil, err
	}

	return buf.Bytes(), nil
}
