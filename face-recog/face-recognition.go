package main

import (
	"dlib/dbus"
	"dlib/glib-2.0"
	"fmt"
	"os"
)

type FaceRecogManager struct{}

type FaceRecogInfo struct {
	CanFaceRecog bool
	PersonName   string
}

const (
	_ENABLE      = "enable"
	_PERSON_NAME = "person_name"
	_CONFIG_FILE = "/.config/face_recognition.cfg"
)

func (info *FaceRecogManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.FaceRecogManager",
		"/com/deepin/daemon/FaceRecogManager",
		"com.deepin.daemon.FaceRecogManager",
	}
}

func (info *FaceRecogManager) GetCanFaceRecognition(userName string) FaceRecogInfo {
	var err error
	faceInfo := FaceRecogInfo{}

	if !FileIsExist(userName) {
		faceInfo.CanFaceRecog, faceInfo.PersonName = CreateFaceRecognition(userName)
		return faceInfo
	}

	configFile := glib.NewKeyFile()
	_, err = configFile.LoadFromFile(GetFaceRecogPath(userName),
		glib.KeyFileFlagsNone)
	if err != nil {
		fmt.Println(err)
		return faceInfo
	}

	faceInfo.CanFaceRecog, err = configFile.GetBoolean(userName, _ENABLE)
	if err != nil {
		fmt.Println(err)
		return faceInfo
	}

	faceInfo.PersonName, err = configFile.GetString(userName, _PERSON_NAME)
	if err != nil {
		fmt.Println(err)
		return faceInfo
	}

	return faceInfo
}

func (info *FaceRecogManager) SetCanFaceRecognition(userName string, enable bool) {
	var (
		data string
		err  error
	)

	configFile := glib.NewKeyFile()
	_, err = configFile.LoadFromFile(GetFaceRecogPath(userName),
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

func FileIsExist(userName string) bool {
	filename := GetFaceRecogPath(userName)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
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

func NewFaceRecogManager() *FaceRecogManager {
	return &FaceRecogManager{}
}

func main() {
	info := NewFaceRecogManager()
	dbus.InstallOnSession(info)
	select {}
}
