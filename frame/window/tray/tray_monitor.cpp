// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "tray_monitor.h"
#include "quicksettingcontroller.h"
#include "pluginsiteminterface.h"
#include "xembedtrayitemwidget.h"
#include "snitrayitemwidget.h"

#define REGISTERTED_WAY_IS_SNI 1
#define REGISTERTED_WAY_IS_XEMBED 2

TrayMonitor::TrayMonitor(QObject *parent)
    : QObject(parent)
    , m_trayInter(new DBusTrayManager(this))
    , m_sniWatcher(new StatusNotifierWatcher("org.kde.StatusNotifierWatcher", "/StatusNotifierWatcher", QDBusConnection::sessionBus(), this))
{
    //-------------------------------Tray Embed---------------------------------------------//
    connect(m_trayInter, &DBusTrayManager::TrayIconsChanged, this, &TrayMonitor::onTrayIconsChanged, Qt::QueuedConnection);
    connect(m_trayInter, &DBusTrayManager::Changed, this, &TrayMonitor::requestUpdateIcon, Qt::QueuedConnection);
    m_trayInter->Manage();
    QMetaObject::invokeMethod(this, "onTrayIconsChanged", Qt::QueuedConnection);

    //-------------------------------Tray SNI---------------------------------------------//
    connect(m_sniWatcher, &StatusNotifierWatcher::StatusNotifierItemRegistered, this, &TrayMonitor::onSniItemsChanged, Qt::QueuedConnection);
    connect(m_sniWatcher, &StatusNotifierWatcher::StatusNotifierItemUnregistered, this, &TrayMonitor::onSniItemsChanged, Qt::QueuedConnection);
    QMetaObject::invokeMethod(this, "onSniItemsChanged", Qt::QueuedConnection);

    //-------------------------------System Tray------------------------------------------//
    QuickSettingController *quickController = QuickSettingController::instance();
    connect(quickController, &QuickSettingController::pluginInserted, this, [ = ](PluginsItemInterface *itemInter, const QuickSettingController::PluginAttribute pluginAttr) {
        if (pluginAttr != QuickSettingController::PluginAttribute::Tray)
            return;

        m_systemTrays << itemInter;
        Q_EMIT systemTrayAdded(itemInter);
    });

    connect(quickController, &QuickSettingController::pluginRemoved, this, [ = ](PluginsItemInterface *itemInter) {
        if (!m_systemTrays.contains(itemInter))
            return;

        m_systemTrays.removeOne(itemInter);
        Q_EMIT systemTrayRemoved(itemInter);
    });

    //-------------------------------Tray Indicator---------------------------------------------//
    // Indicators服务是集成在插件中的，因此需要在所有的插件加载完成后再加载Indicators服务
    connect(quickController, &QuickSettingController::pluginLoaderFinished, this, [ this ] {
        startLoadIndicators();
    });

    QMetaObject::invokeMethod(this, [ = ] {
        QList<PluginsItemInterface *> trayPlugins = quickController->pluginItems(QuickSettingController::PluginAttribute::Tray);
        for (PluginsItemInterface *plugin : trayPlugins) {
            m_systemTrays << plugin;
            Q_EMIT systemTrayAdded(plugin);
        }
    }, Qt::QueuedConnection);
}

QList<quint32> TrayMonitor::trayWinIds() const
{
    return m_trayWids;
}

QStringList TrayMonitor::sniServices() const
{
    return m_sniServices;
}

QStringList TrayMonitor::indicatorNames() const
{
    return m_indicatorNames;
}

QList<PluginsItemInterface *> TrayMonitor::systemTrays() const
{
    return m_systemTrays;
}

void TrayMonitor::onTrayIconsChanged()
{
    QList<quint32> wids = m_trayInter->trayIcons();

    // Filter xembed by tray pid
    QList<quint32> filteredWids;
    for (auto wid : wids) {
        uint pid = XEmbedTrayItemWidget::getWindowPID(wid);
        // filter registered sni tray
        if (m_trayPids.value(pid, REGISTERTED_WAY_IS_XEMBED) == REGISTERTED_WAY_IS_XEMBED) {
            m_trayPids.insert(pid, REGISTERTED_WAY_IS_XEMBED);
            filteredWids.append(wid);
        }
    }
    wids = filteredWids;

    if (m_trayWids == wids)
        return;

    for (auto wid : wids) {
        if (!m_trayWids.contains(wid)) {
            Q_EMIT xEmbedTrayAdded(wid);
        }
    }

    for (auto wid : m_trayWids) {
        if (!wids.contains(wid)) {
            Q_EMIT xEmbedTrayRemoved(wid);
        }
    }

    m_trayWids = wids;
}

void TrayMonitor::onSniItemsChanged()
{
    //TODO 防止同一个进程注册多个sni服务
    const QStringList &sniServices = m_sniWatcher->registeredStatusNotifierItems();

    // Filter sni by pid
    QStringList filteredSniServices;
    for (auto s : sniServices) {
        uint pid = SNITrayItemWidget::servicePID(s);
        // filter registered xembed tray
        // TODO: Priority use of SNI, need remove registered xembed tray?
        if (m_trayPids.value(pid, REGISTERTED_WAY_IS_SNI) == REGISTERTED_WAY_IS_SNI) {
            m_trayPids.insert(pid, REGISTERTED_WAY_IS_SNI);
            filteredSniServices.append(s);
        }
    }

    if (m_sniServices == filteredSniServices)
        return;

    for (auto s : filteredSniServices) {
        if (!m_sniServices.contains(s)) {
            if (s.startsWith("/") || !s.contains("/")) {
                qWarning() << __FUNCTION__ << "invalid sni service" << s;
                continue;
            }
            Q_EMIT sniTrayAdded(s);
        }
    }

    for (auto s : m_sniServices) {
        if (!filteredSniServices.contains(s)) {
            Q_EMIT sniTrayRemoved(s);
        }
    }

    m_sniServices = filteredSniServices;
}

void TrayMonitor::startLoadIndicators()
{
    QDir indicatorConfDir("/etc/dde-dock/indicator");

    for (const QFileInfo &fileInfo : indicatorConfDir.entryInfoList({"*.json"}, QDir::Files | QDir::NoDotAndDotDot)) {
        const QString &indicatorName = fileInfo.baseName();
        m_indicatorNames << indicatorName;
        Q_EMIT indicatorFounded(indicatorName);
    }
}
