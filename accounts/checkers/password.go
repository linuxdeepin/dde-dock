/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package checkers

import (
	. "pkg.deepin.io/lib/gettext"
	"regexp"
	"strings"
)

const (
	passwordMinLength    = 8
	passwordSpecialChars = "~!@#$%^&*()[]{}\\|/?,.<>"
)

var passwordNumberRegexp = regexp.MustCompile("[0-9]")
var passwordUpperAlphabetRegexp = regexp.MustCompile("[A-Z]")
var passwordLowerAlphabetRegexp = regexp.MustCompile("[a-z]")

type passwordErrorCode int32

const (
	passwordOK passwordErrorCode = iota
	passwordErrCodeShort
	passwordErrCodeSimple
)

func (code passwordErrorCode) IsOk() bool {
	return code == passwordOK
}

func (code passwordErrorCode) Prompt() string {
	switch code {
	case passwordOK:
		return ""
	case passwordErrCodeShort:
		return Tr("Please enter a password not less than 8 characters")
	case passwordErrCodeSimple:
		return Tr("The password must contain English letters (case-sensitive), numbers or special symbols (~!@#$%^&*()[]{}\\|/?,.<>)")
	default:
		return ""
	}
}

type password string

func (p password) hasAnyNumber() bool {
	str := string(p)
	return passwordNumberRegexp.MatchString(str)
}

func (p password) hasAnySpecialChar() bool {
	str := string(p)
	return strings.ContainsAny(str, passwordSpecialChars)
}

func (p password) hasUpperAndLowerAlphabet() bool {
	str := string(p)
	return passwordUpperAlphabetRegexp.MatchString(str) &&
		passwordLowerAlphabetRegexp.MatchString(str)
}

func CheckPasswordValid(releaseType, passwd string) passwordErrorCode {
	if releaseType != "Server" {
		return passwordOK
	}

	if len(passwd) < passwordMinLength {
		return passwordErrCodeShort
	}

	p := password(passwd)
	if !p.hasAnyNumber() {
		return passwordErrCodeSimple
	}

	if !p.hasAnySpecialChar() {
		return passwordErrCodeSimple
	}

	if !p.hasUpperAndLowerAlphabet() {
		return passwordErrCodeSimple
	}

	return passwordOK
}
