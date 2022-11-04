/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
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
#include "pluginadapter.h"

#include <QWidget>

#define ICONWIDTH 24
#define ICONHEIGHT 24

PluginAdapter::PluginAdapter(PluginsItemInterface_V20 *pluginInter)
    : m_pluginInter(pluginInter)
{
}

PluginAdapter::~PluginAdapter()
{
    delete m_pluginInter;
}

const QString PluginAdapter::pluginName() const
{
    return m_pluginInter->pluginName();
}

const QString PluginAdapter::pluginDisplayName() const
{
    return m_pluginInter->pluginDisplayName();
}

void PluginAdapter::init(PluginProxyInterface *proxyInter)
{
    m_pluginInter->init(proxyInter);
}

QWidget *PluginAdapter::itemWidget(const QString &itemKey)
{
    return m_pluginInter->itemWidget(itemKey);
}

QWidget *PluginAdapter::itemTipsWidget(const QString &itemKey)
{
    return m_pluginInter->itemTipsWidget(itemKey);
}

QWidget *PluginAdapter::itemPopupApplet(const QString &itemKey)
{
    return m_pluginInter->itemPopupApplet(itemKey);
}

const QString PluginAdapter::itemCommand(const QString &itemKey)
{
    return m_pluginInter->itemCommand(itemKey);
}

const QString PluginAdapter::itemContextMenu(const QString &itemKey)
{
    return m_pluginInter->itemContextMenu(itemKey);
}

void PluginAdapter::invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked)
{
    m_pluginInter->invokedMenuItem(itemKey, menuId, checked);
}

int PluginAdapter::itemSortKey(const QString &itemKey)
{
    return m_pluginInter->itemSortKey(itemKey);
}

void PluginAdapter::setSortKey(const QString &itemKey, const int order)
{
    m_pluginInter->setSortKey(itemKey, order);
}

bool PluginAdapter::itemAllowContainer(const QString &itemKey)
{
    return m_pluginInter->itemAllowContainer(itemKey);
}

bool PluginAdapter::itemIsInContainer(const QString &itemKey)
{
    return m_pluginInter->itemIsInContainer(itemKey);
}

void PluginAdapter::setItemIsInContainer(const QString &itemKey, const bool container)
{
    m_pluginInter->setItemIsInContainer(itemKey, container);
}

bool PluginAdapter::pluginIsAllowDisable()
{
    return m_pluginInter->pluginIsAllowDisable();
}

bool PluginAdapter::pluginIsDisable()
{
    return m_pluginInter->pluginIsDisable();
}

void PluginAdapter::pluginStateSwitched()
{
    m_pluginInter->pluginStateSwitched();
}

void PluginAdapter::displayModeChanged(const Dock::DisplayMode displayMode)
{
    m_pluginInter->displayModeChanged(displayMode);
}

void PluginAdapter::positionChanged(const Dock::Position position)
{
    m_pluginInter->positionChanged(position);
}

void PluginAdapter::refreshIcon(const QString &itemKey)
{
    m_pluginInter->refreshIcon(itemKey);
}

void PluginAdapter::pluginSettingsChanged()
{
    m_pluginInter->pluginSettingsChanged();
}

PluginsItemInterface::PluginType PluginAdapter::type()
{
    switch (m_pluginInter->type()) {
    case PluginsItemInterface_V20::PluginType::Fixed:
        return PluginsItemInterface::PluginType::Fixed;
    case PluginsItemInterface_V20::PluginType::Normal:
        return PluginsItemInterface::PluginType::Normal;
    }
    return PluginsItemInterface::PluginType::Normal;
}

PluginsItemInterface::PluginSizePolicy PluginAdapter::pluginSizePolicy() const
{
    switch (m_pluginInter->pluginSizePolicy()) {
    case PluginsItemInterface_V20::PluginSizePolicy::Custom:
        return PluginsItemInterface::PluginSizePolicy::Custom;
    case PluginsItemInterface_V20::PluginSizePolicy::System:
        return PluginsItemInterface::PluginSizePolicy::System;
    }
    return PluginsItemInterface::PluginSizePolicy::Custom;
}

QIcon PluginAdapter::icon(const DockPart &dockPart)
{
    QWidget *itemWidget = m_pluginInter->itemWidget(m_itemKey);
    if (!itemWidget)
        return QIcon();

    switch (dockPart) {
    case DockPart::QuickPanel: {
        // 如果图标为空，就使用itemWidget的截图作为它的图标，这种一般是适用于老版本插件或者没有实现v23接口的插件
        itemWidget->setFixedSize(ICONWIDTH, ICONHEIGHT);
        return itemWidget->grab();
    }
    case DockPart::SystemPanel: {
        itemWidget->setFixedSize(16, 16);
        return itemWidget->grab();
    }
    default: break;
    }

    return QIcon();
}

PluginsItemInterface::PluginStatus PluginAdapter::status() const
{
    return PluginStatus::Active;
}

QString PluginAdapter::description() const
{
    return tr("actived");
}

void PluginAdapter::setItemKey(const QString &itemKey)
{
    m_itemKey = itemKey;
}
