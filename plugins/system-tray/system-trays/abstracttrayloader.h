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

#ifndef ABSTRACTTRAYLOADER_H
#define ABSTRACTTRAYLOADER_H

#include "abstracttraywidget.h"

#include <QDBusConnectionInterface>
#include <QObject>

class AbstractTrayLoader : public QObject
{
    Q_OBJECT
public:
    explicit AbstractTrayLoader(const QString &waitService, QObject *parent = nullptr);

Q_SIGNALS:
    void systemTrayAdded(const QString &itemKey, AbstractTrayWidget *trayWidget);
    void systemTrayRemoved(const QString &itemKey);

public Q_SLOTS:
    virtual void load() = 0;

public:
    inline bool waitService() { return !m_waitingService.isEmpty(); }
    bool serviceExist();
    void waitServiceForLoad();

private Q_SLOTS:
    void onServiceOwnerChanged(const QString &name, const QString &oldOwner, const QString &newOwner);

private:
    QDBusConnectionInterface *m_dbusDaemonInterface;

    QString m_waitingService;
};

#endif // ABSTRACTTRAYLOADER_H
