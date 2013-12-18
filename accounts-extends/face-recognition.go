package main

import (
	"dlib/dbus"
	"dlib/glib-2.0"
	"fmt"
	"os"
	"os/user"
)

type FaceRecogInfo struct {
	CanFaceRecog bool
	PersonName   string
}

const (
	_ENABLE      = "enable"
	_PERSON_NAME = "person_name"
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
		data     string
		err      error
	)

	userName, err = GetUserNameById(id)
	if err != nil {
		fmt.Println(err)
		return false
	}

	homeDir, _ := GetHomeDirById(id)
	configFile := glib.NewKeyFile()
	_, err = configFile.LoadFromFile(homeDir+_CONFIG_FILE,
		glib.KeyFileFlagsNone)
	if err != nil {
		fmt.Println(err)
		return
	}

	configFile.SetBoolean(userName, _ENABLE, enable)

	_, data, err = configFile.ToData()
	if err != nil {
		fmt.Println(err)
		return
	}
	WriteKeyFile(userName, data)
}

func CreateFaceRecognition(userName string) (bool, string) {
	var (
		f   *os.File
		err error
	)

	filename := GetFaceRecogPath(userName)
	f, err = os.Create(filename)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	configFile := glib.NewKeyFile()
	_, err = configFile.LoadFromFile(GetFaceRecogPath(userName),
		glib.KeyFileFlagsNone)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}
	configFile.SetBoolean(userName, _ENABLE, false)
	uuid := CreateUUID()
	configFile.SetString(userName, _PERSON_NAME, uuid)
	_, data, _ := configFile.ToData()
	_, err = f.WriteString(data)
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	return false, uuid
}

func WriteKeyFile(userName, data string) {
	var (
		f   *os.File
		err error
	)
	f, err = os.Create(GetFaceRecogPath(userName))
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func GetFaceRecogPath(userName string) string {
	path := ""
	path += "/home/" + userName + _CONFIG_FILE

	return path
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
