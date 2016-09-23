/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package network

import (
	"pkg.deepin.io/lib/utils"
	"strings"
)

const (
	nmKeyErrorMissingValue       = "missing value"
	nmKeyErrorEmptyValue         = "value is empty"
	nmKeyErrorInvalidValue       = "invalid value"
	nmKeyErrorIp4MethodConflict  = `%s cannot be used with the 'shared', 'link-local', or 'disabled' methods`
	nmKeyErrorIp4AddressesStruct = "echo IPv4 address structure is composed of 3 32-bit values, address, prefix and gateway"
	// nm.NM_KEY_ERROR_IP4_ADDRESSES_PREFIX = "IPv4 prefix's value should be 1-32"
	nmKeyErrorIp6MethodConflict     = `%s cannot be used with the 'shared', 'link-local', or 'ignore' methods`
	nmKeyErrorMissingSection        = "missing section %s"
	nmKeyErrorEmptySection          = "section %s is empty"
	nmKeyErrorMissingDependsKey     = "missing depends key %s"
	nmKeyErrorMissingDependsPackage = "missing depends package %s"
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
func ensureByteArrayUriPathExistsFor8021x(errs sectionErrors, section, key string, bytePath []byte, limitedExts ...string) {
	path := byteArrayToStrPath(bytePath)
	if !utils.IsURI(path) {
		rememberError(errs, section, key, nmKeyErrorInvalidValue)
		return
	}
	ensureFileExists(errs, section, key, toLocalPathFor8021x(path), limitedExts...)
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
			// TODO dispatch filter when select files
			// rememberError(errs, section, key, nmKeyErrorInvalidValue)
		}
	}
	if !utils.IsFileExist(file) {
		rememberError(errs, section, key, nmKeyErrorInvalidValue)
	}
}
