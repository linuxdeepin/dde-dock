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

#include "soundplugin.h"
#include <QDebug>

#define STATE_KEY  "enable"

SoundPlugin::SoundPlugin(QObject *parent)
    : QObject(parent),
      m_settings("deepin", "dde-dock-sound"),
      m_soundItem(nullptr)
{
}

const QString SoundPlugin::pluginName() const
{
    return "sound";
}

const QString SoundPlugin::pluginDisplayName() const
{
    return tr("Sound");
}

void SoundPlugin::init(PluginProxyInterface *proxyInter)
{
    m_proxyInter = proxyInter;

    m_soundItem = new SoundItem;
    connect(m_soundItem, &SoundItem::requestContextMenu, [this] { m_proxyInter->requestContextMenu(this, QString()); });

    if (m_settings.value(STATE_KEY, true).toBool())
        m_proxyInter->itemAdded(this, QString());
}

void SoundPlugin::pluginStateSwitched()
{
    m_settings.setValue(STATE_KEY, !m_settings.value(STATE_KEY, true).toBool());

    if (m_settings.value(STATE_KEY).toBool())
        m_proxyInter->itemAdded(this, QString());
    else
        m_proxyInter->itemRemoved(this, QString());
}

bool SoundPlugin::pluginIsDisable()
{
    return !m_settings.value(STATE_KEY, true).toBool();
}

QWidget *SoundPlugin::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_soundItem;
}

QWidget *SoundPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_soundItem->tipsWidget();
}

QWidget *SoundPlugin::itemPopupApplet(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_soundItem->popupApplet();
}

const QString SoundPlugin::itemContextMenu(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_soundItem->contextMenu();
}

void SoundPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    Q_UNUSED(itemKey);

    m_soundItem->invokeMenuItem(menuId, checked);
}

int SoundPlugin::itemSortKey(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    Dock::DisplayMode mode = displayMode();
    const QString key = QString("pos_%1").arg(mode);

    if (mode == Dock::DisplayMode::Fashion) {
        return m_settings.value(key, 1).toInt();
    } else {
        return m_settings.value(key, 2).toInt();
    }
}

void SoundPlugin::setSortKey(const QString &itemKey, const int order)
{
    Q_UNUSED(itemKey);

    const QString key = QString("pos_%1").arg(displayMode());
    m_settings.setValue(key, order);
}
