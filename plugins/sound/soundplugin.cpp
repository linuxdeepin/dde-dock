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

    if (!pluginIsDisable())
        m_proxyInter->itemAdded(this, SOUND_KEY);
}

void SoundPlugin::pluginStateSwitched()
{
    m_proxyInter->saveValue(this, STATE_KEY, pluginIsDisable());

    refreshPluginItemsVisible();
}

bool SoundPlugin::pluginIsDisable()
{
    return !m_proxyInter->getValue(this, STATE_KEY, true).toBool();
}

QWidget *SoundPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == SOUND_KEY) {
        return m_soundItem;
    }

    return nullptr;
}

QWidget *SoundPlugin::itemTipsWidget(const QString &itemKey)
{
    if (itemKey == SOUND_KEY) {
        return m_soundItem->tipsWidget();
    }

    return nullptr;
}

QWidget *SoundPlugin::itemPopupApplet(const QString &itemKey)
{
    if (itemKey == SOUND_KEY) {
        return m_soundItem->popupApplet();
    }

    return nullptr;
}

const QString SoundPlugin::itemContextMenu(const QString &itemKey)
{
    if (itemKey == SOUND_KEY) {
        return m_soundItem->contextMenu();
    }

    return QString();
}

void SoundPlugin::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    if (itemKey == SOUND_KEY) {
        m_soundItem->invokeMenuItem(menuId, checked);
    }
}

int SoundPlugin::itemSortKey(const QString &itemKey)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    return m_proxyInter->getValue(this, key, 1).toInt();
}

void SoundPlugin::setSortKey(const QString &itemKey, const int order)
{
    const QString key = QString("pos_%1_%2").arg(itemKey).arg(Dock::Efficient);

    m_proxyInter->saveValue(this, key, order);
}

void SoundPlugin::refreshIcon(const QString &itemKey)
{
    if (itemKey == SOUND_KEY) {
        m_soundItem->refreshIcon();
    }
}

void SoundPlugin::pluginSettingsChanged()
{
    refreshPluginItemsVisible();
}

void SoundPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable())
        m_proxyInter->itemRemoved(this, SOUND_KEY);
    else
        m_proxyInter->itemAdded(this, SOUND_KEY);
}
