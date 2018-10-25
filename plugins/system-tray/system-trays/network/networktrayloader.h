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

#ifndef NETWORKTRAYLOADER_H
#define NETWORKTRAYLOADER_H

#include "../abstracttrayloader.h"
#include "item/abstractnetworktraywidget.h"

#include <QObject>

#include <NetworkWorker>
#include <NetworkModel>

class NetworkTrayLoader : public AbstractTrayLoader
{
    Q_OBJECT
public:
    explicit NetworkTrayLoader(QObject *parent = nullptr);

public Q_SLOTS:
    void load() Q_DECL_OVERRIDE;

private:
    AbstractNetworkTrayWidget *trayWidgetByPath(const QString &path);

private Q_SLOTS:
    void onDeviceListChanged(const QList<dde::network::NetworkDevice *> devices);
    void refreshWiredItemVisible();

private:
    dde::network::NetworkModel *m_networkModel;
    dde::network::NetworkWorker *m_networkWorker;

    QMap<QString, AbstractNetworkTrayWidget *> m_trayWidgetsMap;
    QTimer *m_delayRefreshTimer;
};

#endif // NETWORKTRAYLOADER_H
