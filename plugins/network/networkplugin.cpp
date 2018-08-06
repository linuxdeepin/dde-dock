/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#include "networkplugin.h"
#include "item/wireditem.h"
#include "item/wirelessitem.h"

using namespace dde::network;

#define WIRED_ITEM      "wired"
#define WIRELESS_ITEM   "wireless"
#define STATE_KEY       "enabled"

NetworkPlugin::NetworkPlugin(QObject *parent)
    : QObject(parent),

      m_networkModel(nullptr),
      m_networkWorker(nullptr),
      m_settings("deepin", "dde-dock-network"),
      m_pluginLoaded(false)
{
}

const QString NetworkPlugin::pluginName() const
{
    return "network";
}

const QString NetworkPlugin::pluginDisplayName() const
{
    return tr("Network");
}

void NetworkPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void NetworkPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(checked)

    DeviceItem *item = itemByPath(itemKey);
    if (item) {
        return item->invokeMenuItem(menuId);
    }

    Q_UNREACHABLE();
}

void NetworkPlugin::refershIcon(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    for (auto item : m_itemsMap.values()) {
        item->refreshIcon();
    }
}

void NetworkPlugin::pluginStateSwitched()
{
    m_settings.setValue(STATE_KEY, !m_settings.value(STATE_KEY, true).toBool());

    onDeviceListChanged(m_networkModel->devices());
}

bool NetworkPlugin::pluginIsDisable()
{
    return !m_settings.value(STATE_KEY, true).toBool();
}

const QString NetworkPlugin::itemCommand(const QString &itemKey)
{
    DeviceItem *item = itemByPath(itemKey);
    if (item) {
        return item->itemCommand();
    }

    Q_UNREACHABLE();
    return QString();
}

const QString NetworkPlugin::itemContextMenu(const QString &itemKey)
{
    DeviceItem *item = itemByPath(itemKey);
    if (item) {
        return item->itemContextMenu();
    }

    Q_UNREACHABLE();
    return QString();
}

QWidget *NetworkPlugin::itemWidget(const QString &itemKey)
{
    return itemByPath(itemKey);
}

QWidget *NetworkPlugin::itemTipsWidget(const QString &itemKey)
{
    DeviceItem *item = itemByPath(itemKey);
    if (item) {
        return item->itemTips();
    }

    Q_UNREACHABLE();
    return nullptr;
}

QWidget *NetworkPlugin::itemPopupApplet(const QString &itemKey)
{
    DeviceItem *item = itemByPath(itemKey);
    if (item) {
        return item->itemApplet();
    }

    Q_UNREACHABLE();
    return nullptr;
}

int NetworkPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());
    return m_settings.value(key, 0).toInt();
}

void NetworkPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(displayMode());
    m_settings.setValue(key, order);
}

void NetworkPlugin::onDeviceListChanged(const QList<NetworkDevice *> devices)
{
    QList<QString> mPaths = m_itemsMap.keys();
    QList<QString> newPaths;

    QList<WirelessItem *> wirelessItems;

    for (auto device : devices) {
        const QString &path = device->path();
        newPaths << path;
        // new device
        if (!mPaths.contains(path)) {
            DeviceItem *item = nullptr;
            switch (device->type()) {
                case NetworkDevice::Wired:
                    item = new WiredItem(static_cast<WiredDevice *>(device));
                    break;
                case NetworkDevice::Wireless:
                    item = new WirelessItem(static_cast<WirelessDevice *>(device));
                    wirelessItems.append(static_cast<WirelessItem *>(item));

                    connect(static_cast<WirelessItem *>(item), &WirelessItem::queryActiveConnInfo,
                            m_networkWorker, &NetworkWorker::queryActiveConnInfo);
                    connect(static_cast<WirelessItem *>(item), &WirelessItem::requestActiveAP,
                            m_networkWorker, &NetworkWorker::activateAccessPoint);
                    connect(static_cast<WirelessItem *>(item), &WirelessItem::requestDeactiveAP,
                            m_networkWorker, &NetworkWorker::disconnectDevice);
                    connect(static_cast<WirelessItem *>(item), &WirelessItem::feedSecret,
                            m_networkWorker, &NetworkWorker::feedSecret);
                    connect(static_cast<WirelessItem *>(item), &WirelessItem::cancelSecret,
                            m_networkWorker, &NetworkWorker::cancelSecret);
                    connect(static_cast<WirelessItem *>(item), &WirelessItem::requestWirelessScan,
                            m_networkWorker, &NetworkWorker::requestWirelessScan);

                    connect(m_networkModel, &NetworkModel::needSecrets,
                            static_cast<WirelessItem *>(item), &WirelessItem::onNeedSecrets);
                    connect(m_networkModel, &NetworkModel::needSecretsFinished,
                            static_cast<WirelessItem *>(item), &WirelessItem::onNeedSecretsFinished);

                    m_networkWorker->queryAccessPoints(path);
                    m_networkWorker->requestWirelessScan();
                    break;
                default:
                    Q_UNREACHABLE();
            }

            mPaths << path;
            m_itemsMap.insert(path, item);

            connect(item, &DeviceItem::requestContextMenu, this, &NetworkPlugin::contextMenuRequested);
            connect(item, &DeviceItem::requestSetDeviceEnable, m_networkWorker, &NetworkWorker::setDeviceEnable);
        }
    }

    for (auto mPath : mPaths) {
        // removed device
        if (!newPaths.contains(mPath)) {
            m_proxyInter->itemRemoved(this, mPath);
            m_itemsMap.take(mPath)->deleteLater();
            break;
        }

        if (m_settings.value(STATE_KEY, true).toBool()) {
            m_proxyInter->itemAdded(this, mPath);
        } else {
            m_proxyInter->itemRemoved(this, mPath);
        }
    }

    int wirelessItemCount = wirelessItems.size();
    for (int i = 0; i < wirelessItemCount; ++i) {
        QTimer::singleShot(1, [=] {
            wirelessItems.at(i)->setDeviceInfo(wirelessItemCount == 1 ? -1 : i + 1);
        });
    }
}

DeviceItem *NetworkPlugin::itemByPath(const QString &path)
{
    for (auto item : m_itemsMap.values()) {
        if (item->path() == path) {
            return item;
        }
    }

    Q_UNREACHABLE();
    return nullptr;
}

void NetworkPlugin::loadPlugin()
{
    if (m_pluginLoaded) {
        qDebug() << "network plugin has been loaded! return";
        return;
    }

    m_pluginLoaded = true;

    m_networkModel = new NetworkModel;
    m_networkWorker = new NetworkWorker(m_networkModel);

    connect(m_networkModel, &NetworkModel::deviceListChanged, this, &NetworkPlugin::onDeviceListChanged);

    m_networkModel->moveToThread(qApp->thread());
    m_networkWorker->moveToThread(qApp->thread());

    onDeviceListChanged(m_networkModel->devices());
}

void NetworkPlugin::contextMenuRequested()
{
    DeviceItem *item = qobject_cast<DeviceItem *>(sender());
    Q_ASSERT(item);

    m_proxyInter->requestContextMenu(this, item->path());
}
