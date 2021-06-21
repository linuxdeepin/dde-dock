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
#include "networkpanel.h"

#include <DDBusSender>

#include <unetworkcontroller.h>

#define STATE_KEY       "enabled"

NetworkPlugin::NetworkPlugin(QObject *parent)
    : QObject(parent)
    , m_networkPanel(Q_NULLPTR)
{
}

NetworkPlugin::~NetworkPlugin()
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
    if(m_networkPanel)
        return;

    m_networkPanel.reset(new NetworkPanel);

    if (!pluginIsDisable())
        loadPlugin();
}

void NetworkPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(checked)

    if (itemKey == NETWORK_KEY)
         m_networkPanel->invokeMenuItem(menuId);
}

void NetworkPlugin::refreshIcon(const QString &itemKey)
{
    if (itemKey == NETWORK_KEY)
        m_networkPanel->refreshIcon();
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
    if (!hasDevice() || m_networkPanel->needShowControlCenter()) {
        return QString("dbus-send --print-reply "
                       "--dest=com.deepin.dde.ControlCenter "
                       "/com/deepin/dde/ControlCenter "
                       "com.deepin.dde.ControlCenter.ShowModule "
                       "\"string:network\"");
    }

    return QString();
}

const QString NetworkPlugin::itemContextMenu(const QString &itemKey)
{
    if (itemKey == NETWORK_KEY)
        return m_networkPanel->contextMenu();

    return QString();
}

QWidget *NetworkPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == NETWORK_KEY)
        return m_networkPanel.data();

    return Q_NULLPTR;
}

QWidget *NetworkPlugin::itemTipsWidget(const QString &itemKey)
{
    if (itemKey == NETWORK_KEY)
        return m_networkPanel->itemTips();

    return Q_NULLPTR;
}

QWidget *NetworkPlugin::itemPopupApplet(const QString &itemKey)
{
    if (itemKey == NETWORK_KEY && hasDevice() && !m_networkPanel->needShowControlCenter())
        return m_networkPanel->itemApplet();

    return Q_NULLPTR;
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

void NetworkPlugin::loadPlugin()
{
    m_proxyInter->itemAdded(this, NETWORK_KEY);
}

void NetworkPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable())
        m_proxyInter->itemRemoved(this, NETWORK_KEY);
    else
        m_proxyInter->itemAdded(this, NETWORK_KEY);
}

bool NetworkPlugin::hasDevice()
{
    if(m_networkPanel.isNull())
        return false;

    return m_networkPanel.data()->hasDevice();
}
