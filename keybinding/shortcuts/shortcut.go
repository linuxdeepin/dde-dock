/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package shortcuts

import (
	"errors"
	"strconv"
)

type BaseShortcut struct {
	Id     string
	Type   int32
	Accels []ParsedAccel
	Name   string
}

func (sb *BaseShortcut) GetId() string {
	return sb.Id
}

func idType2Uid(id string, _type int32) string {
	return strconv.Itoa(int(_type)) + id
}

func (sb *BaseShortcut) GetUid() string {
	return idType2Uid(sb.Id, sb.Type)
}

func (sb *BaseShortcut) GetAccels() []ParsedAccel {
	return sb.Accels
}

func (sb *BaseShortcut) setAccels(newAccels []ParsedAccel) {
	sb.Accels = newAccels
}

func (sb *BaseShortcut) GetType() int32 {
	return sb.Type
}

func (sb *BaseShortcut) GetName() string {
	return sb.Name
}

func (sb *BaseShortcut) GetAction() *Action {
	return ActionNoOp
}

func (sb *BaseShortcut) SetAction(newAction *Action) error {
	return ErrOpNotSupported
}

func (sb *BaseShortcut) SetName(name string) error {
	return ErrOpNotSupported
}

const (
	ShortcutTypeSystem int32 = iota
	ShortcutTypeCustom
	ShortcutTypeMedia
	ShortcutTypeWM
	ShortcutTypeMetacity
)

type Shortcut interface {
	GetId() string
	GetUid() string
	GetType() int32

	GetName() string
	SetName(name string) error

	GetAccels() []ParsedAccel
	setAccels(newAccels []ParsedAccel)
	SaveAccels() error
	ReloadAccels() bool

	GetAction() *Action
	SetAction(newAction *Action) error
}

// errors:
var ErrOpNotSupported = errors.New("operation is not supported")
var ErrTypeAssertionFail = errors.New("type assertion failed")
var ErrNilAction = errors.New("action is nil")
var ErrInvalidActionType = errors.New("invalid action type")
