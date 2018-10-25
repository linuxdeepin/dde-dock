/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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

#include "powertrayloader.h"
#include "powertraywidget.h"

#define PowerItemKey "system-tray-power"
#define PowerService "com.deepin.daemon.Power"

PowerTrayLoader::PowerTrayLoader(QObject *parent)
    : AbstractTrayLoader(PowerService, parent),
      m_powerInter(new DBusPower(this))
{
}

void PowerTrayLoader::load()
{
    if (!m_powerInter->batteryState().isEmpty()) {
        emit systemTrayAdded(PowerItemKey, new PowerTrayWidget);
    }
}
