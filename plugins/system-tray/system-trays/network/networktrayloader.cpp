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

#include "networktrayloader.h"
#include "item/wiredtraywidget.h"
#include "item/wirelesstraywidget.h"

using namespace dde::network;

#define NetworkItemKeyPrefix "system-tray-network-"
#define NetworkService "com.deepin.daemon.Network"

NetworkTrayLoader::NetworkTrayLoader(QObject *parent)
    : AbstractTrayLoader(NetworkService, parent),
      m_networkModel(nullptr),
      m_networkWorker(nullptr),
      m_delayRefreshTimer(new QTimer)
{
    m_delayRefreshTimer->setSingleShot(true);
    m_delayRefreshTimer->setInterval(2000);

    connect(m_delayRefreshTimer, &QTimer::timeout, this, &NetworkTrayLoader::refreshWiredItemVisible);
}

void NetworkTrayLoader::load()
{
    m_networkModel = new NetworkModel;
    m_networkWorker = new NetworkWorker(m_networkModel);

    connect(m_networkModel, &NetworkModel::deviceListChanged, this, &NetworkTrayLoader::onDeviceListChanged);

    m_networkModel->moveToThread(qApp->thread());
    m_networkWorker->moveToThread(qApp->thread());

    onDeviceListChanged(m_networkModel->devices());
}

AbstractNetworkTrayWidget *NetworkTrayLoader::trayWidgetByPath(const QString &path)
{
    for (auto trayWidget : m_trayWidgetsMap.values()) {
        if (trayWidget->path() == path) {
            return trayWidget;
        }
    }

    Q_UNREACHABLE();
    return nullptr;
}

void NetworkTrayLoader::onDeviceListChanged(const QList<dde::network::NetworkDevice *> devices)
{
    QList<QString> mPaths = m_trayWidgetsMap.keys();
    QList<QString> newPaths;

    QList<WirelessTrayWidget *> wirelessTrayList;

    for (auto device : devices) {
        const QString &path = device->path();
        newPaths << path;
        // new device
        if (!mPaths.contains(path)) {
            AbstractNetworkTrayWidget *networkTray = nullptr;
            switch (device->type()) {
                case NetworkDevice::Wired:
                    networkTray = new WiredTrayWidget(static_cast<WiredDevice *>(device));
                    break;
                case NetworkDevice::Wireless:
                    networkTray = new WirelessTrayWidget(static_cast<WirelessDevice *>(device));
                    wirelessTrayList.append(static_cast<WirelessTrayWidget *>(networkTray));

                    connect(static_cast<WirelessTrayWidget *>(networkTray), &WirelessTrayWidget::queryActiveConnInfo,
                            m_networkWorker, &NetworkWorker::queryActiveConnInfo);
                    connect(static_cast<WirelessTrayWidget *>(networkTray), &WirelessTrayWidget::requestActiveAP,
                            m_networkWorker, &NetworkWorker::activateAccessPoint);
                    connect(static_cast<WirelessTrayWidget *>(networkTray), &WirelessTrayWidget::requestDeactiveAP,
                            m_networkWorker, &NetworkWorker::disconnectDevice);
                    connect(static_cast<WirelessTrayWidget *>(networkTray), &WirelessTrayWidget::feedSecret,
                            m_networkWorker, &NetworkWorker::feedSecret);
                    connect(static_cast<WirelessTrayWidget *>(networkTray), &WirelessTrayWidget::cancelSecret,
                            m_networkWorker, &NetworkWorker::cancelSecret);
                    connect(static_cast<WirelessTrayWidget *>(networkTray), &WirelessTrayWidget::requestWirelessScan,
                            m_networkWorker, &NetworkWorker::requestWirelessScan);
                    connect(static_cast<WirelessTrayWidget *>(networkTray), &WirelessTrayWidget::createApConfig,
                            m_networkWorker, &NetworkWorker::createApConfig);
                    connect(static_cast<WirelessTrayWidget *>(networkTray), &WirelessTrayWidget::queryConnectionSession,
                            m_networkWorker, &NetworkWorker::queryConnectionSession);

                    connect(m_networkModel, &NetworkModel::needSecrets,
                            static_cast<WirelessTrayWidget *>(networkTray), &WirelessTrayWidget::onNeedSecrets);
                    connect(m_networkModel, &NetworkModel::needSecretsFinished,
                            static_cast<WirelessTrayWidget *>(networkTray), &WirelessTrayWidget::onNeedSecretsFinished);

                    m_networkWorker->queryAccessPoints(path);
                    m_networkWorker->requestWirelessScan();
                    break;
                default:
                    Q_UNREACHABLE();
            }

            mPaths << path;
            m_trayWidgetsMap.insert(path, networkTray);

            connect(device, &dde::network::NetworkDevice::enableChanged,
                    m_delayRefreshTimer, static_cast<void (QTimer:: *)()>(&QTimer::start));

//            connect(networkTray, &AbstractNetworkTrayWidget::requestContextMenu, this, &NetworkPlugin::contextMenuRequested);
            connect(networkTray, &AbstractNetworkTrayWidget::requestSetDeviceEnable, m_networkWorker, &NetworkWorker::setDeviceEnable);
        }
    }

    for (auto mPath : mPaths) {
        // removed device
        if (!newPaths.contains(mPath)) {
            Q_EMIT systemTrayRemoved(NetworkItemKeyPrefix + mPath);
            m_trayWidgetsMap.take(mPath)->deleteLater();
            break;
        }

        Q_EMIT systemTrayAdded(NetworkItemKeyPrefix + mPath, m_trayWidgetsMap.value(mPath));
    }

    int wirelessItemCount = wirelessTrayList.size();
    for (int i = 0; i < wirelessItemCount; ++i) {
        QTimer::singleShot(0, this, [=] {
            wirelessTrayList.at(i)->setDeviceInfo(wirelessItemCount == 1 ? -1 : i + 1);
        });
    }

    m_delayRefreshTimer->start();
}

void NetworkTrayLoader::refreshWiredItemVisible()
{
    bool hasWireless = false;
    QList<WiredTrayWidget *> wiredTrayList;

    for (auto trayWidget : m_trayWidgetsMap.values()) {
        if (trayWidget->device()->type() == NetworkDevice::Wireless) {
            hasWireless = true;
        } else {
            wiredTrayList.append(static_cast<WiredTrayWidget *>(trayWidget));
        }
    }

    if (!hasWireless) {
        return;
    }

    for (auto wiredTrayWidget : wiredTrayList) {
        if (!wiredTrayWidget->device()->enabled()) {
            Q_EMIT systemTrayRemoved(NetworkItemKeyPrefix + wiredTrayWidget->path());
        } else {
            Q_EMIT systemTrayAdded(NetworkItemKeyPrefix + wiredTrayWidget->path(), wiredTrayWidget);
        }
    }
}
