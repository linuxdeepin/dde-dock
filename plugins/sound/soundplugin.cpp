// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "soundplugin.h"
#include "soundaccessible.h"

#include <QDebug>
#include <QAccessible>

#define STATE_KEY  "enable"

SoundPlugin::SoundPlugin(QObject *parent)
    : QObject(parent),
      m_soundItem(nullptr)
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
        return m_soundItem.data();
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

void SoundPlugin::refreshPluginItemsVisible()
{
    if (pluginIsDisable())
        m_proxyInter->itemRemoved(this, SOUND_KEY);
    else
        m_proxyInter->itemAdded(this, SOUND_KEY);
}
