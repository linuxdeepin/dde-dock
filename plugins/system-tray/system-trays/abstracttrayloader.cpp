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

#include "abstracttrayloader.h"

#include <QDebug>

AbstractTrayLoader::AbstractTrayLoader(const QString &waitService, QObject *parent)
    : QObject(parent),
      m_dbusDaemonInterface(QDBusConnection::sessionBus().interface()),
      m_waitingService(waitService)
{
}

bool AbstractTrayLoader::serviceExist()
{
    bool exist = m_dbusDaemonInterface->isServiceRegistered(m_waitingService).value();

    if (!exist) {
        qDebug() << m_waitingService << "daemon has not started";
    }

    return exist;
}

void AbstractTrayLoader::waitServiceForLoad()
{
    connect(m_dbusDaemonInterface, &QDBusConnectionInterface::serviceOwnerChanged, this, &AbstractTrayLoader::onServiceOwnerChanged, Qt::UniqueConnection);
}

void AbstractTrayLoader::onServiceOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner)
{
    Q_UNUSED(oldOwner);

    if (m_waitingService.isEmpty() || newOwner.isEmpty()) {
        return;
    }

    if (m_waitingService == name) {
        qDebug() << m_waitingService << "daemon started, load tray and disconnect";
        load();
        disconnect(m_dbusDaemonInterface);
        m_waitingService = QString();
    }
}
