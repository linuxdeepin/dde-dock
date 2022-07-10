/*
 * Copyright (C) 2020 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     weizhixiang <weizhixiang@uniontech.com>
 *
 * Maintainer: weizhixiang <weizhixiang@uniontech.com>
 *
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

#include "airplanemodeplugin.h"
#include "airplanemodeitem.h"

#define AIRPLANEMODE_KEY "airplane-mode-key"
#define STATE_KEY  "enable"

AirplaneModePlugin::AirplaneModePlugin(QObject *parent)
    : QObject(parent)
    , m_item(new AirplaneModeItem)
{
    connect(m_item, &AirplaneModeItem::airplaneEnableChanged, this, &AirplaneModePlugin::onAirplaneEnableChanged);
}

const QString AirplaneModePlugin::pluginName() const
{
    return "airplane-mode";
}

const QString AirplaneModePlugin::pluginDisplayName() const
{
    return tr("Airplane Mode");
}

void AirplaneModePlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    if (!pluginIsDisable())
        m_proxyInter->itemAdded(this, AIRPLANEMODE_KEY);

    bool hasBluetooth = false;
    QDBusInterface inter("com.deepin.daemon.Bluetooth",
                    "/com/deepin/daemon/Bluetooth",
                    "com.deepin.daemon.Bluetooth",
                    QDBusConnection::sessionBus());
    if (inter.isValid()) {
        QDBusReply<QString> reply = inter.call("GetAdapters");
        QString replyStr = reply.value();
        replyStr = replyStr.replace("[", "");
        replyStr = replyStr.replace("]", "");
        if (!replyStr.isEmpty()) {
            hasBluetooth = true;
        }
    }

    bool hasWireless = false;
    NetworkInter networkInter("com.deepin.daemon.Network", "/com/deepin/daemon/Network", QDBusConnection::sessionBus(), this);
    if (networkInter.isValid()) {
        if (!networkInter.wirelessAccessPoints().isEmpty()) {
            hasWireless = true;
        }
    }
    // 蓝牙和无线网络,只要有其中一个就允许显示飞行模式
    if (hasBluetooth || hasWireless) {
        refreshAirplaneEnableState();
    }
}

void AirplaneModePlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, STATE_KEY, pluginIsDisable());

    refreshAirplaneEnableState();
}

bool AirplaneModePlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, STATE_KEY, true).toBool();
}

QWidget *AirplaneModePlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == AIRPLANEMODE_KEY) {
        return m_item;
    }

    return nullptr;
}

QWidget *AirplaneModePlugin::itemTipsWidget(const QString &itemKey)
{
    if (itemKey == AIRPLANEMODE_KEY) {
        return m_item->tipsWidget();
    }

    return nullptr;
}

QWidget *AirplaneModePlugin::itemPopupApplet(const QString &itemKey)
{
    if (itemKey == AIRPLANEMODE_KEY) {
        return m_item->popupApplet();
    }

    return nullptr;
}

const QString AirplaneModePlugin::itemContextMenu(const QString &itemKey)
{
    if (itemKey == AIRPLANEMODE_KEY) {
        return m_item->contextMenu();
    }

    return QString();
}

void AirplaneModePlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    if (itemKey == AIRPLANEMODE_KEY) {
        m_item->invokeMenuItem(menuId, checked);
    }
}

int AirplaneModePlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    return m_proxyInter->getValue(this, key, 4).toInt();
}

void AirplaneModePlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    m_proxyInter->saveValue(this, key, order);
}

void AirplaneModePlugin::refreshIcon(const QString &itemKey)
{
    if (itemKey == AIRPLANEMODE_KEY) {
        m_item->refreshIcon();
    }
}

void AirplaneModePlugin::refreshAirplaneEnableState()
{
    onAirplaneEnableChanged(m_item->airplaneEnable());
}

void AirplaneModePlugin::onAirplaneEnableChanged(bool enable)
{
    if (!m_proxyInter)
        return;

    m_proxyInter->itemAdded(this, AIRPLANEMODE_KEY);
    if (enable) {
        m_proxyInter->saveValue(this, STATE_KEY, true);
    }
    else {
        m_proxyInter->saveValue(this, STATE_KEY, false);
    }
}


