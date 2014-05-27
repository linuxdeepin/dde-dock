package main

// TODO doc

const NM_SETTING_SERIAL_SETTING_NAME = "serial"

const (
	NM_SETTING_SERIAL_BAUD       = "baud"
	NM_SETTING_SERIAL_BITS       = "bits"
	NM_SETTING_SERIAL_PARITY     = "parity"
	NM_SETTING_SERIAL_STOPBITS   = "stopbits"
	NM_SETTING_SERIAL_SEND_DELAY = "send-delay"
)

// Get available keys
func getSettingSerialAvailableKeys(data connectionData) (keys []string) {
	return
}

// Get available values
func getSettingSerialAvailableValues(data connectionData, key string) (values []kvalue) {
	return
}

// Check whether the values are correct
func checkSettingSerialValues(data connectionData) (errs sectionErrors) {
	errs = make(map[string]string)
	return
}
