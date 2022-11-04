/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#include "tray_monitor.h"
#include "quicksettingcontroller.h"
#include "pluginsiteminterface.h"

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

    QMetaObject::invokeMethod(this, [ = ] {
        QList<PluginsItemInterface *> trayPlugins = quickController->pluginItems(QuickSettingController::PluginAttribute::Tray);
        for (PluginsItemInterface *plugin : trayPlugins) {
            m_systemTrays << plugin;
            Q_EMIT systemTrayAdded(plugin);
        }
    }, Qt::QueuedConnection);

    //-------------------------------Tray Indicator---------------------------------------------//
    QMetaObject::invokeMethod(this, "startLoadIndicators", Qt::QueuedConnection);
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
    if (m_sniServices == sniServices)
        return;

    for (auto s : sniServices) {
        if (!m_sniServices.contains(s)) {
            if (s.startsWith("/") || !s.contains("/")) {
                qWarning() << __FUNCTION__ << "invalid sni service" << s;
                continue;
            }
            Q_EMIT sniTrayAdded(s);
        }
    }

    for (auto s : m_sniServices) {
        if (!sniServices.contains(s)) {
            Q_EMIT sniTrayRemoved(s);
        }
    }

    m_sniServices = sniServices;
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
