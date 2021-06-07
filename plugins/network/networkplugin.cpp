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
#include "networkitem.h"
#include "item/wireditem.h"
#include "item/wirelessitem.h"

#include <DDBusSender>

using namespace dde::network;

#define STATE_KEY       "enabled"

NetworkPlugin::NetworkPlugin(QObject *parent)
    : QObject(parent)
    , m_networkModel(nullptr)
    , m_networkWorker(nullptr)
    , m_networkItem(nullptr)
    , m_hasDevice(false)
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

    if (m_networkItem)
        return;

    m_networkItem = new NetworkItem;

    if (!pluginIsDisable()) {
        loadPlugin();
    }
}

void NetworkPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    if (itemKey == NETWORK_KEY) {
        return m_networkItem->invokeMenuItem(menuId, checked);
    }
}

void NetworkPlugin::refreshIcon(const QString &itemKey)
{
    if (itemKey == NETWORK_KEY) {
        m_networkItem->refreshIcon();
    }
}

void NetworkPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, STATE_KEY, pluginIsDisable());

    refreshPluginItemsVisible();
}

bool NetworkPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, STATE_KEY, true).toBool();
}

const QString NetworkPlugin::itemCommand(const QString &itemKey)
{
    Q_UNUSED(itemKey)
    return m_hasDevice ? QString() : QString("dbus-send --print-reply "
                                              "--dest=com.deepin.dde.ControlCenter "
                                              "/com/deepin/dde/ControlCenter "
                                              "com.deepin.dde.ControlCenter.ShowModule "
                                              "\"string:network\"");
}

const QString NetworkPlugin::itemContextMenu(const QString &itemKey)
{
    if (itemKey == NETWORK_KEY) {
        return m_networkItem->contextMenu();
    }

    return QString();
}

QWidget *NetworkPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == NETWORK_KEY) {
        return m_networkItem;
    }

    return nullptr;
}

QWidget *NetworkPlugin::itemTipsWidget(const QString &itemKey)
{
    if (itemKey == NETWORK_KEY) {
        return m_networkItem->itemTips();
    }

    return nullptr;
}

QWidget *NetworkPlugin::itemPopupApplet(const QString &itemKey)
{
    if (itemKey == NETWORK_KEY && m_hasDevice) {
        return m_networkItem->itemApplet();
    }

    return nullptr;
}

int NetworkPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    return m_proxyInter->getValue(this, key, 3).toInt();
}

void NetworkPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    m_proxyInter->saveValue(this, key, order);
}

void NetworkPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

//bool NetworkPlugin::isConnectivity()
//{
//    return NetworkModel::connectivity() == Connectivity::Full;
//}

void NetworkPlugin::onDeviceListChanged(const QList<NetworkDevice *> devices)
{
    QMap<QString, WirelessItem*> wirelessItems;
    QMap<QString, WiredItem *> wiredItems;

    int wiredDeviceCnt = 0;
    int wirelessDeviceCnt = 0;
    for (auto device : devices) {
        if (device && device->type() == NetworkDevice::Wired)
            wiredDeviceCnt++;
        else
            wirelessDeviceCnt++;
    }

    // 编号 (与控制中心有线设备保持一致命名)
    int wiredNum = 0;
    int wirelessNum = 0;
    QString text;

    for (auto device : devices) {
        const QString &path = device->path();
        // new device
        DeviceItem *item = nullptr;
        switch (device->type()) {
        case NetworkDevice::Wired:
            wiredNum++;
            if (wiredDeviceCnt == 1)
                text = tr("Wired Network");
            else
                text = tr("Wired Network %1").arg(wiredNum);
            item = new WiredItem(static_cast<WiredDevice *>(device), text);
            wiredItems.insert(path, static_cast<WiredItem *>(item));

            connect(static_cast<WiredItem *>(item), &WiredItem::wiredStateChanged,
                    m_networkItem, &NetworkItem::updateSelf);
            connect(static_cast<WiredItem *>(item), &WiredItem::enableChanged,
                    m_networkItem, &NetworkItem::updateSelf);
            connect(static_cast<WiredItem *>(item), &WiredItem::activeConnectionChanged,
                    m_networkItem, &NetworkItem::updateSelf);
            break;
        case NetworkDevice::Wireless:
            item = new WirelessItem(static_cast<WirelessDevice *>(device));
            static_cast<WirelessItem *>(item)->setDeviceInfo(wirelessDeviceCnt == 1 ? -1 : ++wirelessNum);
            wirelessItems.insert(path, static_cast<WirelessItem *>(item));

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
            connect(static_cast<WirelessItem *>(item), &WirelessItem::createApConfig,
                    m_networkWorker, &NetworkWorker::createApConfig);
            connect(static_cast<WirelessItem *>(item), &WirelessItem::queryConnectionSession,
                    m_networkWorker, &NetworkWorker::queryConnectionSession);

            connect(static_cast<WirelessItem *>(item), &WirelessItem::deviceStateChanged,
                    m_networkItem, &NetworkItem::updateSelf);
            connect(static_cast<WirelessItem *>(item), &WirelessItem::requestWirelessScan,
                    m_networkItem, &NetworkItem::wirelessScan);

            m_networkWorker->queryAccessPoints(path);
            m_networkWorker->requestWirelessScan();
            break;
        default:
            Q_UNREACHABLE();
        }

        connect(item, &DeviceItem::requestSetDeviceEnable, m_networkWorker, &NetworkWorker::setDeviceEnable);
        connect(m_networkModel, &NetworkModel::connectivityChanged, item, &DeviceItem::refreshConnectivity);
        connect(m_networkModel, &NetworkModel::connectivityChanged, m_networkItem, &NetworkItem::updateSelf);
    }

    m_hasDevice = wiredItems.size() || wirelessItems.size();
    m_networkItem->updateDeviceItems(wiredItems, wirelessItems);
}

void NetworkPlugin::loadPlugin()
{
    m_networkModel = new NetworkModel;
    m_networkWorker = new NetworkWorker(m_networkModel);

    connect(m_networkModel, &NetworkModel::deviceListChanged, this, &NetworkPlugin::onDeviceListChanged);

    m_networkModel->moveToThread(qApp->thread());
    m_networkWorker->moveToThread(qApp->thread());

    onDeviceListChanged(m_networkModel->devices());

    m_proxyInter->itemAdded(this, NETWORK_KEY);
}

void NetworkPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable()) {
        m_proxyInter->itemRemoved(this, NETWORK_KEY);
    } else {
        m_proxyInter->itemAdded(this, NETWORK_KEY);
    }
}
