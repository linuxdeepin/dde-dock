/*
 * Copyright (C) 2023 ~ 2023 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
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
#include "pluginmanager.h"
#include "dockplugincontroller.h"
#include "quicksettingcontainer.h"
#include "iconmanager.h"

#include <QResizeEvent>

PluginManager::PluginManager(QObject *parent)
    : m_dockController(nullptr)
    , m_iconManager(nullptr)
{
}

const QString PluginManager::pluginName() const
{
    return "pluginManager";
}

const QString PluginManager::pluginDisplayName() const
{
    return "pluginManager";
}

void PluginManager::init(PluginProxyInterface *proxyInter)
{
    if (m_proxyInter == proxyInter)
        return;

    m_proxyInter = proxyInter;

    m_dockController.reset(new DockPluginController(proxyInter));
    m_quickContainer.reset(new QuickSettingContainer(m_dockController.data()));
    m_iconManager.reset(new IconManager(m_dockController.data()));
    m_iconManager->setPosition(position());

    connect(m_dockController.data(), &DockPluginController::pluginInserted, this, [ this ](PluginsItemInterface *itemInter) {
        if (m_iconManager->isFixedPlugin(itemInter)) {
            m_proxyInter->itemUpdate(this, pluginName());
        }
    });
    connect(m_dockController.data(), &DockPluginController::pluginUpdated, this, [ this ](PluginsItemInterface *itemInter) {
        if (m_iconManager->isFixedPlugin(itemInter)) {
            m_proxyInter->itemUpdate(this, pluginName());
        }
    });
    connect(m_dockController.data(), &DockPluginController::pluginRemoved, this, [ this ](PluginsItemInterface *itemInter) {
        if (m_iconManager->isFixedPlugin(itemInter)) {
            m_proxyInter->itemUpdate(this, pluginName());
        }
    });
    connect(m_dockController.data(), &DockPluginController::requestAppletVisible, this, [ this ](PluginsItemInterface *itemInter, const QString &itemKey, bool visible) {
        if (visible) {
            QWidget *appletWidget = itemInter->itemPopupApplet(itemKey);
            if (appletWidget)
                m_quickContainer->showPage(appletWidget, itemInter);
        } else {
            // 显示主面板
            m_quickContainer->topLevelWidget()->hide();
        }
    });
    connect(m_dockController.data(), &DockPluginController::pluginLoadFinished, this, &PluginManager::pluginLoadFinished);

    // 开始加载插件
    m_dockController->startLoadPlugin(getPluginPaths());

    m_proxyInter->itemAdded(this, pluginName());
}

QWidget *PluginManager::itemWidget(const QString &itemKey)
{
    Q_UNUSED(itemKey);
    return nullptr;
}

QWidget *PluginManager::itemPopupApplet(const QString &itemKey)
{
    if (itemKey == QUICK_ITEM_KEY) {
        // 返回快捷面板
        return m_quickContainer.data();
    }

    return nullptr;
}

QIcon PluginManager::icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType)
{
    if (dockPart == DockPart::QuickShow) {
        return m_iconManager->pixmap();
    }

    return QIcon();
}

PluginFlags PluginManager::flags() const
{
    // 当前快捷插件组合区域只支持在快捷区域显示
    return PluginFlag::Type_Common | PluginFlag::Attribute_ForceDock;
}

PluginsItemInterface::PluginSizePolicy PluginManager::pluginSizePolicy() const
{
    return PluginSizePolicy::Custom;
}

bool PluginManager::eventHandler(QEvent *event)
{
    if (event->type() == QEvent::Resize) {
        QResizeEvent *resizeEvent = static_cast<QResizeEvent *>(event);
        m_iconManager->updateSize(resizeEvent->size());
    }
    return PluginsItemInterface::eventHandler(event);
}

void PluginManager::positionChanged(const Dock::Position position)
{
    m_iconManager->setPosition(position);
    m_proxyInter->itemUpdate(this, pluginName());
}

QList<PluginsItemInterface *> PluginManager::plugins() const
{
    return m_dockController->plugins();
}

QList<PluginsItemInterface *> PluginManager::pluginsInSetting() const
{
    return m_dockController->pluginsInSetting();
}

QList<PluginsItemInterface *> PluginManager::currentPlugins() const
{
    return m_dockController->currentPlugins();
}

QString PluginManager::itemKey(PluginsItemInterface *itemInter) const
{
    return m_dockController->itemKey(itemInter);
}

QJsonObject PluginManager::metaData(PluginsItemInterface *itemInter) const
{
    return m_dockController->metaData(itemInter);
}

#ifndef QT_DEBUG
static QStringList getPathFromConf(const QString &key) {
    QSettings set("/etc/deepin/dde-dock.conf", QSettings::IniFormat);
    auto value = set.value(key).toString();
    if (!value.isEmpty()) {
        return value.split(':');
    }

    return QStringList();
}
#endif

QStringList PluginManager::getPluginPaths() const
{
    QStringList pluginPaths;
#ifdef QT_DEBUG
    pluginPaths << QString("%1/..%2").arg(qApp->applicationDirPath()).arg(QUICK_PATH)
                << QString("%1/..%2").arg(qApp->applicationDirPath()).arg(PLUGIN_PATH)
                << QString("%1/..%2").arg(qApp->applicationDirPath()).arg(TRAY_PATH);
#else
    pluginPaths << QString("/usr/lib/dde-dock%1").arg(QUICK_PATH)
                << QString("/usr/lib/dde-dock%1").arg(PLUGIN_PATH)
                << QString("/usr/lib/dde-dock%1").arg(TRAY_PATH);

    const QStringList pluginsDirs = (getPathFromConf("PATH") << getPathFromConf("SYSTEM_TRAY_PATH"));
    if (!pluginsDirs.isEmpty())
        pluginPaths << pluginsDirs;
#endif

    return pluginPaths;
}
