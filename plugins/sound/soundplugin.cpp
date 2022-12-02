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
#include "soundaccessible.h"
#include "soundwidget.h"
#include "sounddeviceswidget.h"

#include <QDebug>
#include <QAccessible>

#define STATE_KEY  "enable"

SoundPlugin::SoundPlugin(QObject *parent)
    : QObject(parent)
    , m_soundItem(nullptr)
    , m_soundWidget(nullptr)
{
    QAccessible::installFactory(soundAccessibleFactory);
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

    if (m_soundItem)
        return;

    m_soundItem.reset(new SoundItem);
    m_soundWidget.reset(new SoundWidget);
    m_soundWidget->setFixedHeight(60);

    m_soundDeviceWidget.reset(new SoundDevicesWidget);

    if (!pluginIsDisable()) {
        m_proxyInter->itemAdded(this, SOUND_KEY);
        connect(m_soundWidget.data(), &SoundWidget::rightIconClick, this, [ this, proxyInter ] {
            proxyInter->requestSetAppletVisible(this, QUICK_ITEM_KEY, true);
        });
    }

    connect(m_soundDeviceWidget.data(), &SoundDevicesWidget::enableChanged, m_soundWidget.data(), &SoundWidget::setEnabled);
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
    if (itemKey == SOUND_KEY)
        return m_soundItem.data();

    if (itemKey == QUICK_ITEM_KEY)
        return m_soundWidget.data();

    return nullptr;
}

QWidget *SoundPlugin::itemTipsWidget(const QString &itemKey)
{
    if (itemKey == SOUND_KEY)
        return m_soundItem->tipsWidget();

    return nullptr;
}

QWidget *SoundPlugin::itemPopupApplet(const QString &itemKey)
{
    if (itemKey == SOUND_KEY)
        return m_soundItem->popupApplet();

    if (itemKey == QUICK_ITEM_KEY)
        return m_soundDeviceWidget.data();

    return nullptr;
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

    return m_proxyInter->getValue(this, key, 2).toInt();
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

QIcon SoundPlugin::icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType)
{
    return m_soundItem->pixmap(themeType);
}

PluginsItemInterface::PluginMode SoundPlugin::status() const
{
    return SoundPlugin::Active;
}

PluginFlags SoundPlugin::flags() const
{
    return PluginFlag::Type_Common
            | PluginFlag::Quick_Full
            | PluginFlag::Attribute_CanDrag
            | PluginFlag::Attribute_CanInsert
            | PluginFlag::Attribute_CanSetting;
}

void SoundPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable())
        m_proxyInter->itemRemoved(this, SOUND_KEY);
    else
        m_proxyInter->itemAdded(this, SOUND_KEY);
}
