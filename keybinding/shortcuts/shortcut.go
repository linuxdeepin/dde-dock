/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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

func (sb *BaseShortcut) GetAccelsModifiable() bool {
	return sb.Type != ShortcutTypeFake
}

const (
	ShortcutTypeSystem int32 = iota
	ShortcutTypeCustom
	ShortcutTypeMedia
	ShortcutTypeWM
	ShortcutTypeFake
)

type Shortcut interface {
	GetId() string
	GetUid() string
	GetType() int32

	GetName() string

	GetAccelsModifiable() bool
	GetAccels() []ParsedAccel
	setAccels(newAccels []ParsedAccel)
	SaveAccels() error
	ReloadAccels() bool

	GetAction() *Action
}

// errors:
var ErrOpNotSupported = errors.New("operation is not supported")
var ErrTypeAssertionFail = errors.New("type assertion failed")
var ErrNilAction = errors.New("action is nil")
var ErrInvalidActionType = errors.New("invalid action type")

type FakeShortcut struct {
	BaseShortcut
	action *Action
}

func NewFakeShortcut(action *Action) *FakeShortcut {
	return &FakeShortcut{
		BaseShortcut: BaseShortcut{
			Type: ShortcutTypeFake,
		},
		action: action,
	}
}

func (s *FakeShortcut) GetAction() *Action {
	return s.action
}

func (s *FakeShortcut) SaveAccels() error {
	return ErrOpNotSupported
}

func (s *FakeShortcut) ReloadAccels() bool {
	return false
}
