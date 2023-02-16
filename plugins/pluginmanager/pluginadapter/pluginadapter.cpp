// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "pluginadapter.h"

#include <QWidget>

#define ICONWIDTH 24
#define ICONHEIGHT 24

PluginAdapter::PluginAdapter(PluginsItemInterface_V20 *pluginInter, QPluginLoader *pluginLoader)
    : m_pluginInter(pluginInter)
    , m_pluginLoader(pluginLoader)
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

QIcon PluginAdapter::icon(const DockPart &dockPart, DGuiApplicationHelper::ColorType themeType)
{
    QWidget *itemWidget = m_pluginInter->itemWidget(m_itemKey);
    if (!itemWidget)
        return QIcon();

    switch (dockPart) {
    case DockPart::QuickPanel:
    case DockPart::SystemPanel: {
        // 如果图标为空，就使用itemWidget的截图作为它的图标，这种一般是适用于老版本插件或者没有实现v23接口的插件
        QSize oldSize = itemWidget->size();
        itemWidget->setFixedSize(ICONWIDTH, ICONHEIGHT);
        QPixmap pixmap = itemWidget->grab();
        itemWidget->setFixedSize(oldSize);
        return pixmap;
    }
    default: break;
    }

    return QIcon();
}

PluginsItemInterface::PluginMode PluginAdapter::status() const
{
    return PluginMode::Active;
}

QString PluginAdapter::description() const
{
    return m_pluginInter->pluginDisplayName();
}

PluginFlags PluginAdapter::flags() const
{
    if (m_pluginLoader->fileName().contains(TRAY_PATH))
        return PluginFlag::Type_Tray | PluginFlag::Attribute_CanDrag | PluginFlag::Attribute_CanInsert;

    return PluginsItemInterface::flags();
}

void PluginAdapter::setItemKey(const QString &itemKey)
{
    m_itemKey = itemKey;
}
