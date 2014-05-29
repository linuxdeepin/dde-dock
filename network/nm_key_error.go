package network

import (
	"strings"
)

const (
	NM_KEY_ERROR_MISSING_VALUE        = "missing value"
	NM_KEY_ERROR_EMPTY_VALUE          = "value is empty"
	NM_KEY_ERROR_INVALID_VALUE        = "invalid value"
	NM_KEY_ERROR_IP4_METHOD_CONFLICT  = `%s cannot be used with the 'shared', 'link-local', or 'disabled' methods`
	NM_KEY_ERROR_IP4_ADDRESSES_STRUCT = "echo IPv4 address structure is composed of 3 32-bit values, address, prefix and gateway"
	// NM_KEY_ERROR_IP4_ADDRESSES_PREFIX = "IPv4 prefix's value should be 1-32"
	NM_KEY_ERROR_IP6_METHOD_CONFLICT = `%s cannot be used with the 'shared', 'link-local', or 'ignore' methods`
	NM_KEY_ERROR_MISSING_SECTION     = "missing section %s"
	NM_KEY_ERROR_EMPTY_SECTION       = "section %s is empty"
	NM_KEY_ERROR_MISSING_DEPENDS_KEY = "missing depends key %s"
)

func rememberError(errs sectionErrors, section, key, errMsg string) {
	relatedVks := getRelatedVkeys(section, key)
	if len(relatedVks) > 0 {
		rememberVkError(errs, section, key, errMsg)
		return
	}
	doRememberError(errs, key, errMsg)
}

func rememberVkError(errs sectionErrors, section, key, errMsg string) {
	vks := getRelatedVkeys(section, key)
	for _, vk := range vks {
		if !isOptionalVkey(section, vk) {
			doRememberError(errs, vk, errMsg)
		}
	}
}

func doRememberError(errs sectionErrors, key, errMsg string) {
	if _, ok := errs[key]; ok {
		// error already exists for this key
		return
	}
	errs[key] = errMsg
}

// start with "file://", end with null byte
func ensureByteArrayUriPathExists(errs sectionErrors, section, key string, bytePath []byte, limitedExts ...string) {
	path := byteArrayToStrPath(bytePath)
	if !isUriPath(path) {
		rememberError(errs, section, key, NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	ensureFileExists(errs, section, key, path, limitedExts...)
}

func ensureFileExists(errs sectionErrors, section, key, file string, limitedExts ...string) {
	file = toLocalPath(file)
	// ensure file suffix with target extension
	if len(limitedExts) > 0 {
		match := false
		for _, ext := range limitedExts {
			if strings.HasSuffix(strings.ToLower(file), strings.ToLower(ext)) {
				match = true
				break
			}
		}
		if !match {
			// TODO
			// rememberError(errs, section, key, NM_KEY_ERROR_INVALID_VALUE)
		}
	}
	if !isFileExists(file) {
		rememberError(errs, section, key, NM_KEY_ERROR_INVALID_VALUE)
	}
}
