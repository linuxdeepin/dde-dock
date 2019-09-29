package util

import (
	"bytes"
	"encoding/json"
	"strings"

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

func (kwa *KWinAccel) fix() {
	var keystrokes []string
	for _, ks := range kwa.Keystrokes {
		if ks == "" {
			continue
		}
		keystrokes = append(keystrokes, ks)
	}
	kwa.Keystrokes = keystrokes

	var defaultKeystrokes []string
	for _, ks := range kwa.DefaultKeystrokes {
		if ks == "" || strings.Contains(ks, " ") {
			continue
		}
		defaultKeystrokes = append(defaultKeystrokes, ks)
	}
	if len(defaultKeystrokes) > 1 {
		defaultKeystrokes = defaultKeystrokes[len(defaultKeystrokes)-1:]
	}

	kwa.DefaultKeystrokes = defaultKeystrokes
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

	for idx := range result {
		result[idx].fix()
	}

	return result, nil
}
