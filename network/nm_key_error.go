package main

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
	NM_KEY_ERROR_MISSING_SECTION     = "missing %s field section"
	NM_KEY_ERROR_EMPTY_SECTION       = "field section %s is empty"
	NM_KEY_ERROR_MISSING_DEPENDS_KEY = "missing depends key %s"
)

func rememberError(errs fieldErrors, field, key, errMsg string) {
	relatedVks := getRelatedVirtualKeys(field, key)
	if len(relatedVks) > 0 {
		rememberVkError(errs, field, key, errMsg)
		return
	}
	doRememberError(errs, key, errMsg)
}

func rememberVkError(errs fieldErrors, field, key, errMsg string) {
	vks := getRelatedVirtualKeys(field, key)
	for _, vk := range vks {
		if !isOptionalChildVirtualKeys(field, vk) {
			doRememberError(errs, vk, errMsg)
		}
	}
}

func doRememberError(errs fieldErrors, key, errMsg string) {
	if _, ok := errs[key]; ok {
		// error already exists for this key
		return
	}
	errs[key] = errMsg
}

// start with "file://", end with null byte
func ensureByteArrayUriPathExists(errs fieldErrors, field, key string, bytePath []byte, limitedExts ...string) {
	path := byteArrayToStrPath(bytePath)
	if !isUriPath(path) {
		rememberError(errs, field, key, NM_KEY_ERROR_INVALID_VALUE)
		return
	}
	ensureFileExists(errs, field, key, path, limitedExts...)
}

func ensureFileExists(errs fieldErrors, field, key, file string, limitedExts ...string) {
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
			// rememberError(errs, field, key, NM_KEY_ERROR_INVALID_VALUE)
		}
	}
	if !isFileExists(file) {
		rememberError(errs, field, key, NM_KEY_ERROR_INVALID_VALUE)
	}
}
