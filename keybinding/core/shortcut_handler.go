/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package core

import (
	"sync"
)

var locker sync.Mutex

type Handler struct {
	Shortcut string
	// Hnadler (state, detail, pressed)
	Handler HandleType
}
type Handlers []*Handler

func NewHandler(s string, handler HandleType) *Handler {
	return &Handler{Shortcut: s, Handler: handler}
}

func (handlers Handlers) AddHandler(handler *Handler) Handlers {
	h := handlers.GetHandler(handler.Shortcut)
	if h != nil {
		handlers = handlers.DeleteHandler(h.Shortcut)
	}
	locker.Lock()
	defer locker.Unlock()
	handlers = append(handlers, handler)
	return handlers
}

func (handlers Handlers) DeleteHandler(s string) Handlers {
	locker.Lock()
	defer locker.Unlock()

	var (
		newHandlers Handlers
	)
	for _, h := range handlers {
		if IsAccelEqual(h.Shortcut, s) {
			continue
		}
		newHandlers = append(newHandlers, h)
	}
	return newHandlers
}

func (handlers Handlers) GetHandler(s string) *Handler {
	locker.Lock()
	defer locker.Unlock()

	for _, handler := range handlers {
		if IsAccelEqual(s, handler.Shortcut) {
			return handler
		}
	}
	return nil
}

func (handlers Handlers) GetHandlerByKeycode(mod uint16, keycode int) *Handler {
	locker.Lock()
	defer locker.Unlock()

	for _, handler := range handlers {
		if IsKeyMatch(handler.Shortcut, mod, keycode) {
			return handler
		}
	}
	return nil
}
