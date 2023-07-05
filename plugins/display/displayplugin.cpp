// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "displayplugin.h"
#include "brightnesswidget.h"
#include "displaysettingwidget.h"
#include "brightnesswidget.h"
#include "brightnessmodel.h"

#include "../../widgets/tipswidget.h"
#include "../../frame/util/utils.h"

#include <DDBusSender>

#include <QDebug>
#include <QDBusConnectionInterface>

#include <unistd.h>

using namespace Dock;
DisplayPlugin::DisplayPlugin(QObject *parent)
    : QObject(parent)
    , m_displayWidget(nullptr)
    , m_displaySettingWidget(nullptr)
    , m_displayTips(nullptr)
    , m_model(nullptr)
{
}

const QString DisplayPlugin::pluginName() const
{
    return "display";
}

const QString DisplayPlugin::pluginDisplayName() const
{
    return tr("Brightness");
}

void DisplayPlugin::init(PluginProxyInterface *proxyInter)
{
    if (m_proxyInter == proxyInter)
        return;

    m_proxyInter = proxyInter;
    m_displayTips.reset(new TipsWidget);
    m_model.reset(new BrightnessModel);
    m_displayWidget.reset(new BrightnessWidget(m_model.data()));
    m_displayWidget->setFixedHeight(60);
    m_displaySettingWidget.reset(new DisplaySettingWidget);

    if (m_model->monitors().size() > 0)
        m_proxyInter->itemAdded(this, pluginName());

    connect(m_displayWidget.data(), &BrightnessWidget::brightClicked, this, [ this ] {
        m_proxyInter->requestSetAppletVisible(this, QUICK_ITEM_KEY, true);
    });
    connect(m_displaySettingWidget.data(), &DisplaySettingWidget::requestHide, this, [ this ] {
        m_proxyInter->requestSetAppletVisible(this, QUICK_ITEM_KEY, false);
    });
    connect(m_model.data(), &BrightnessModel::screenVisibleChanged, this, [ this ](bool visible) {
        if (visible)
            m_proxyInter->itemAdded(this, pluginName());
        else
            m_proxyInter->itemRemoved(this, pluginName());
    });
}

QWidget *DisplayPlugin::itemWidget(const QString &itemKey)
{
    if (itemKey == QUICK_ITEM_KEY) {
        return m_displayWidget.data();
    }

    return nullptr;
}

QWidget *DisplayPlugin::itemTipsWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);

    return m_displayTips.data();
}

QWidget *DisplayPlugin::itemPopupApplet(const QString &itemKey)
{
    if (itemKey == QUICK_ITEM_KEY)
        return m_displaySettingWidget.data();

    return nullptr;
}

PluginFlags DisplayPlugin::flags() const
{
    return PluginFlag::Type_Common | PluginFlag::Quick_Full;
}
