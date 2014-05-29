package bluetooth

import (
	"encoding/json"
)

func isStringInArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func marshalJSON(v interface{}) (strJSON string) {
	byteJSON, err := json.Marshal(v)
	if err != nil {
		logger.Error(err)
		return
	}
	strJSON = string(byteJSON)
	return
}

func isDBusObjectKeyExists(data dbusObjectData, key string) (ok bool) {
	_, ok = data[key]
	return
}

func getDBusObjectValueString(data dbusObjectData, key string) (r string) {
	v, ok := data[key]
	if ok {
		r = interfaceToString(v.Value())
	}
	return
}

func getDBusObjectValueInt16(data dbusObjectData, key string) (r int16) {
	v, ok := data[key]
	if ok {
		r = interfaceToInt16(v.Value())
	}
	return
}

func interfaceToString(v interface{}) (r string) {
	r, _ = v.(string)
	return
}

func interfaceToInt16(v interface{}) (r int16) {
	r, _ = v.(int16)
	return
}
