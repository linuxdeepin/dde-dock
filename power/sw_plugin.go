// +build sw

/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

func init() {
	submoduleList = append(submoduleList, newSwPlugin)
}

type swPlugin struct {
	manager *Manager
}

func newSwPlugin(m *Manager) (string, submodule, error) {
	return "swPlugin", &swPlugin{manager: m}, nil
}

func (p *swPlugin) Start() error {
	logger.Debug("swPlugin: force use percentage for policy")
	m := p.manager
	m.usePercentageForPolicy = true
	m.powerSupplyDataBackend = powerSupplyDataBackendPoll
	return nil
}

func (p *swPlugin) Destroy() {
	return
}
