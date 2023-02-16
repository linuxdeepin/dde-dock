// Copyright (C) 2023 ~ 2023 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
    m_iconManager->setDisplayMode(displayMode());

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
        return m_iconManager->pixmap(themeType);
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

void PluginManager::positionChanged(const Dock::Position position)
{
    m_iconManager->setPosition(position);
    m_proxyInter->itemUpdate(this, pluginName());
}

void PluginManager::displayModeChanged(const Dock::DisplayMode displayMode)
{
    m_iconManager->setDisplayMode(displayMode);
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
