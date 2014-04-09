package main

const (
	NM_KEY_ERROR_MISSING_VALUE        = "missing value"
	NM_KEY_ERROR_EMPTY_VALUE          = "value is empty"
	NM_KEY_ERROR_INVALID_VALUE        = "invalid value"
	NM_KEY_ERROR_IP4_METHOD_CONFLICT  = `%s cannot be used with the 'shared', 'link-local', or 'disabled' methods`
	NM_KEY_ERROR_IP4_ADDRESSES_STRUCT = "echo IPv4 address structure is composed of 3 32-bit values, address, prefix and gateway"
	// NM_KEY_ERROR_IP4_ADDRESSES_PREFIX = "IPv4 prefix's value should be 1-32"
	NM_KEY_ERROR_IP6_METHOD_CONFLICT = `%s cannot be used with the 'shared', 'link-local', or 'ignore' methods`
	NM_KEY_ERROR_MISSING_SECTION     = "missing %s field section"
)

func rememberError(errs map[string]string, key, errMsg string) {
	if _, ok := errs[key]; ok {
		// error already exists for this key
		return
	}
	errs[key] = errMsg
}

func rememberErrorForVirtualKey(errs map[string]string, field, key, errMsg string) {
	vks := getRelatedVirtualKeys(field, key)
	for _, vk := range vks {
		rememberError(errs, vk, errMsg)
	}
}
