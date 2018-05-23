/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     rekols <rekols@foxmail.com>
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

#include "dbusadaptors.h"
#include <QDebug>

DBusAdaptors::DBusAdaptors(QObject *parent)
    : QDBusAbstractAdaptor(parent),
      m_keyboard(new Keyboard("com.deepin.daemon.InputDevices",
                              "/com/deepin/daemon/InputDevice/Keyboard",
                              QDBusConnection::sessionBus(), this))
{
    connect(m_keyboard, &Keyboard::CurrentLayoutChanged, this, &DBusAdaptors::onLayoutChanged);
}

DBusAdaptors::~DBusAdaptors()
{
}

QString DBusAdaptors::layout() const
{
    const auto layouts = m_keyboard->currentLayout().split(';');

    if (!layouts.isEmpty())
        return layouts.first();

    qWarning() << Q_FUNC_INFO << "layouts is Empty!!";

    // re-fetch data.
    QTimer::singleShot(1000, this, &DBusAdaptors::onLayoutChanged);

    return QString();
}

void DBusAdaptors::onLayoutChanged()
{
    emit layoutChanged(layout());
}
