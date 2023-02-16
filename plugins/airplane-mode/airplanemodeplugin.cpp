// Copyright (C) 2020 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
    if (m_proxyInter == proxyInter)
        return;

    m_proxyInter = proxyInter;

    m_proxyInter->itemAdded(this, AIRPLANEMODE_KEY);

    refreshAirplaneEnableState();
}

void AirplaneModePlugin::pluginStateSwitched()
{
    refreshAirplaneEnableState();
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

    if (enable) {
        m_proxyInter->itemAdded(this, AIRPLANEMODE_KEY);
        m_proxyInter->saveValue(this, STATE_KEY, true);
    }
    else {
        m_proxyInter->itemRemoved(this, AIRPLANEMODE_KEY);
        m_proxyInter->saveValue(this, STATE_KEY, false);
    }
}


