// Copyright (C) 2023 ~ 2023 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "iconmanager.h"
#include "dockplugincontroller.h"
#include "pluginsiteminterface.h"

#include <DDciIcon>
#include <DWindowManagerHelper>
#include <DSysInfo>
#include <DPlatformTheme>

#include <QPainter>
#include <QPainterPath>

#define ITEMSPACE 6
#define IMAGESIZE 12
#define ITEMSIZE 18
#define MINISIZE 1
#define STARTPOS 2

DGUI_USE_NAMESPACE

IconManager::IconManager(DockPluginController *pluginController, QObject *parent)
    : QObject{parent}
    , m_pluginController(pluginController)
    , m_position(Dock::Position::Bottom)
    , m_displayMode(Dock::DisplayMode::Efficient)
{
}

void IconManager::setPosition(Dock::Position position)
{
    m_position = position;
}

void IconManager::setDisplayMode(Dock::DisplayMode displayMode)
{
    m_displayMode = displayMode;
}

QPixmap IconManager::pixmap(DGuiApplicationHelper::ColorType colorType) const
{
    // 缺省图标
    return QIcon::fromTheme("dock-control-panel").pixmap(ITEMSIZE, ITEMSIZE);
}

PluginsItemInterface *IconManager::findPlugin(const QString &pluginName) const
{
    QList<PluginsItemInterface *> plugins = m_pluginController->currentPlugins();
    for (PluginsItemInterface *plugin : plugins) {
        if (plugin->pluginName() == pluginName)
            return plugin;
    }

    return nullptr;
}
