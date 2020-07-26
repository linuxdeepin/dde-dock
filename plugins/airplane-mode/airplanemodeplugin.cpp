/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     zhaolong<zhaolong@uniontech.com>
 *
 * Maintainer:  xiehui<xiehui@uniontech.com>
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
{
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

    m_item = new AirplaneModeItem;

    if (!pluginIsDisable())
        m_proxyInter->itemAdded(this, AIRPLANEMODE_KEY);
}

void AirplaneModePlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, STATE_KEY, pluginIsDisable());
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


