package util

import (
	"bytes"
	"encoding/json"

	wm "github.com/linuxdeepin/go-dbus-factory/com.deepin.wm"
)

func MarshalJSON(v interface{}) (string, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

type KWinAccel struct {
	Id                string
	Keystrokes        []string `json:"Accels"`
	DefaultKeystrokes []string `json:"Default,omitempty"`
}

func GetAllKWinAccels(wm *wm.Wm) ([]KWinAccel, error) {
	allJson, err := wm.GetAllAccels(0)
	if err != nil {
		return nil, err
	}

	var result []KWinAccel
	err = json.Unmarshal([]byte(allJson), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
