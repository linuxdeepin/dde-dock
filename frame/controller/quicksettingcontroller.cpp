// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "quicksettingcontroller.h"
#include "pluginsitem.h"
#include "pluginmanagerinterface.h"

#include <QMetaObject>
#include <customevent.h>

QuickSettingController::QuickSettingController(QObject *parent)
    : AbstractPluginsController(parent)
{
    qApp->installEventFilter(this);
    // 只有在非安全模式下才加载插件，安全模式会在等退出安全模式后通过接受事件的方式来加载插件
    if (!qApp->property("safeMode").toBool())
        QMetaObject::invokeMethod(this, &QuickSettingController::startLoader, Qt::QueuedConnection);
}

QuickSettingController::~QuickSettingController()
{
}

bool QuickSettingController::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == qApp && event->type() == PluginLoadEvent::eventType()) {
        // 如果收到的是重新加载插件的消息（一般是在退出安全模式后），则直接加载插件即可
        startLoader();
    }

    return AbstractPluginsController::eventFilter(watched, event);
}

void QuickSettingController::startLoader()
{
#ifdef QT_DEBUG
        AbstractPluginsController::startLoader(new PluginLoader(QString("%1/..%2").arg(qApp->applicationDirPath()).arg("/plugins/loader"), this));
#else
        AbstractPluginsController::startLoader(new PluginLoader("/usr/lib/dde-dock/plugins/loader", this));
#endif
}

void QuickSettingController::itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    // 根据读取到的metaData数据获取当前插件的类型，提供给外部
    PluginAttribute pluginAttr = pluginAttribute(itemInter);
    m_quickPlugins[pluginAttr] << itemInter;

    emit pluginInserted(itemInter, pluginAttr);
}

void QuickSettingController::itemUpdate(PluginsItemInterface * const itemInter, const QString &)
{
    updateDockInfo(itemInter, DockPart::QuickPanel);
    updateDockInfo(itemInter, DockPart::QuickShow);
    updateDockInfo(itemInter, DockPart::SystemPanel);
}

void QuickSettingController::itemRemoved(PluginsItemInterface * const itemInter, const QString &)
{
    for (auto it = m_quickPlugins.begin(); it != m_quickPlugins.end(); it++) {
        QList<PluginsItemInterface *> &plugins = m_quickPlugins[it.key()];
        if (!plugins.contains(itemInter))
            continue;

        plugins.removeOne(itemInter);
        if (plugins.isEmpty()) {
            QuickSettingController::PluginAttribute pluginclass = it.key();
            m_quickPlugins.remove(pluginclass);
        }

        break;
    }

    Q_EMIT pluginRemoved(itemInter);
}

void QuickSettingController::requestSetAppletVisible(PluginsItemInterface * const itemInter, const QString &itemKey, const bool visible)
{
    // 设置插件列表可见事件
    Q_EMIT requestAppletVisible(itemInter, itemKey, visible);
}

void QuickSettingController::updateDockInfo(PluginsItemInterface * const itemInter, const DockPart &part)
{
    Q_EMIT pluginUpdated(itemInter, part);
}

QuickSettingController::PluginAttribute QuickSettingController::pluginAttribute(PluginsItemInterface * const itemInter) const
{
    // 工具插件，例如回收站
    if (itemInter->flags() & PluginFlag::Type_Tool)
        return PluginAttribute::Tool;

    // 系统插件，例如关机按钮
    if (itemInter->flags() & PluginFlag::Type_System)
        return PluginAttribute::System;

    // 托盘插件，例如磁盘图标
    if (itemInter->flags() & PluginFlag::Type_Tray)
        return PluginAttribute::Tray;

    // 固定插件，例如显示桌面和多任务试图
    if (itemInter->flags() & PluginFlag::Type_Fixed)
        return PluginAttribute::Fixed;

    // 通用插件，一般的插件都是通用插件，就是放在快捷插件区域的那些插件
    if (itemInter->flags() & PluginFlag::Type_Common)
        return PluginAttribute::Quick;

    // 基本插件，不在任务栏上显示的插件
    return PluginAttribute::None;
}

QString QuickSettingController::itemKey(PluginsItemInterface *pluginItem) const
{
    PluginManagerInterface *pManager = pluginManager();
    if (pManager)
        return pManager->itemKey(pluginItem);

    return QString();
}

QuickSettingController *QuickSettingController::instance()
{
    static QuickSettingController instance;
    return &instance;
}

QList<PluginsItemInterface *> QuickSettingController::pluginItems(const PluginAttribute &pluginClass) const
{
    return m_quickPlugins.value(pluginClass);
}

QJsonObject QuickSettingController::metaData(PluginsItemInterface *pluginItem)
{
    PluginManagerInterface *pManager = pluginManager();
    if (pManager)
        return pManager->metaData(pluginItem);

    return QJsonObject();
}

PluginsItem *QuickSettingController::pluginItemWidget(PluginsItemInterface *pluginItem)
{
    if (m_pluginItemWidgetMap.contains(pluginItem))
        return m_pluginItemWidgetMap[pluginItem];

    PluginsItem *widget = new PluginsItem(pluginItem, itemKey(pluginItem), metaData(pluginItem));
    m_pluginItemWidgetMap[pluginItem] = widget;
    return widget;
}

QList<PluginsItemInterface *> QuickSettingController::pluginInSettings()
{
    PluginManagerInterface *pManager = pluginManager();
    if (!pManager)
        return QList<PluginsItemInterface *>();

    // 返回可用于在控制中心显示的插件
    return pManager->pluginsInSetting();
}
