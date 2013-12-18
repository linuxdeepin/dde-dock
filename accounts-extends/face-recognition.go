package main

import (
	"dlib/glib-2.0"
	"fmt"
	"os"
	"os/user"
	"strconv"
)

type FaceRecogInfo struct {
	CanFaceRecog bool
	PersonName   string
}

const (
	_ENABLE      = "enable"
	_PERSON_NAME = "person_name"
	_CONFIG_DIR  = "/.config/deepin-system-settings/account/"
	_CONFIG_FILE = "/.config/deepin-system-settings/account/face_recognition.cfg"
)

func (info *AccountExtendsManager) CanFaceRecognition(id string) *FaceRecogInfo {
	var (
		err     error
		success bool
		homeDir string
		uuidStr string
	)

	homeDir, err = GetHomeDirById(id)
	if err != nil {
		fmt.Println(err)
		return &FaceRecogInfo{CanFaceRecog: false, PersonName: ""}
	}

	if !FileIsExist(homeDir + _CONFIG_FILE) {
		return &FaceRecogInfo{CanFaceRecog: false, PersonName: ""}
	}

	configFile := glib.NewKeyFile()
	_, err = configFile.LoadFromFile(homeDir+_CONFIG_FILE,
		glib.KeyFileFlagsNone)
	if err != nil {
		fmt.Println(err)
		return &FaceRecogInfo{CanFaceRecog: false, PersonName: ""}
	}

	userName, _ := GetUserNameById(id)
	success, err = configFile.GetBoolean(userName, _ENABLE)
	if err != nil {
		fmt.Println(err)
		return &FaceRecogInfo{CanFaceRecog: false, PersonName: ""}
	}

	uuidStr, err = configFile.GetString(userName, _PERSON_NAME)
	if err != nil {
		fmt.Println(err)
		return &FaceRecogInfo{CanFaceRecog: false, PersonName: ""}
	}

	return &FaceRecogInfo{CanFaceRecog: success, PersonName: uuidStr}
}

func (info *AccountExtendsManager) SetFaceRecognition(id string, enable bool) bool {
	var (
		userName string
		uuidStr  string
		data     string
		err      error
	)

	userName, err = GetUserNameById(id)
	if err != nil {
		fmt.Println(err)
		return false
	}

	homeDir, _ := GetHomeDirById(id)
	if !DirIsExist(homeDir + _CONFIG_DIR) {
		return false
	}

	if !FileIsExist(homeDir + _CONFIG_FILE) {
		f, err1 := os.Create(homeDir + _CONFIG_FILE)
		if err1 != nil {
			fmt.Println("Create face config file failed:", err1)
			return false
		}
		userInfo, _ := user.LookupId(id)
		uid, _ := strconv.ParseInt(userInfo.Uid, 10, 64)
		gid, _ := strconv.ParseInt(userInfo.Gid, 10, 64)
		f.Chown(int(uid), int(gid))

		f.Close()
	}

	configFile := glib.NewKeyFile()
	_, err = configFile.LoadFromFile(homeDir+_CONFIG_FILE,
		glib.KeyFileFlagsNone)
	if err != nil {
		fmt.Println(err)
		return false
	}

	configFile.SetBoolean(userName, _ENABLE, enable)
	uuidStr, err = configFile.GetString(userName, _PERSON_NAME)
	if err != nil || uuidStr == "" {
		uuid := CreateUUID()
		configFile.SetString(userName, _PERSON_NAME, uuid)
	}

	_, data, err = configFile.ToData()
	if err != nil {
		fmt.Println(err)
		return false
	}
	if !WriteKeyFile(homeDir+_CONFIG_FILE, data) {
		return false
	}

	return true
}

func (info *AccountExtendsManager) CreateConfigFile(id string) bool {
	var (
		f        *os.File
		err      error
		homeDir  string
		userName string
	)

	userInfo, err := user.LookupId(id)
	if err != nil {
		fmt.Println(err)
		return false
	}

	homeDir = userInfo.HomeDir
	if !DirIsExist(homeDir + _CONFIG_DIR) {
		return false
	}

	userName = userInfo.Username
	filename := homeDir + _CONFIG_FILE
	f, err = os.Create(filename)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer f.Close()

	uid, _ := strconv.ParseInt(userInfo.Uid, 10, 64)
	gid, _ := strconv.ParseInt(userInfo.Gid, 10, 64)
	f.Chown(int(uid), int(gid))

	configFile := glib.NewKeyFile()
	_, err = configFile.LoadFromFile(filename,
		glib.KeyFileFlagsNone)
	if err != nil {
		fmt.Println(err)
		return false
	}
	configFile.SetBoolean(userName, _ENABLE, false)
	uuid := CreateUUID()
	configFile.SetString(userName, _PERSON_NAME, uuid)
	_, data, _ := configFile.ToData()
	_, err = f.WriteString(data)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func WriteKeyFile(filename, data string) bool {
	var (
		f   *os.File
		err error
	)
	f, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0664)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer f.Close()

	_, err = f.WriteString(data)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func CreateUUID() string {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer f.Close()
	b := make([]byte, 16)
	f.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6],
		b[6:8], b[8:10], b[10:])

	return uuid
}

func GetUserNameById(id string) (string, error) {
	userInfo, err := user.LookupId(id)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return userInfo.Username, nil
}

func DirIsExist(path string) bool {
	err := os.MkdirAll(path, 0666)
	if err != nil {
		fmt.Printf("Dir '%s' failed: %s\n", err)
		return false
	}

	return true
}
