// Copyright (C) 2023 ~ 2023 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef ICONMANAGER_H
#define ICONMANAGER_H

#include "pluginsiteminterface.h"

#include <QObject>
#include <QPixmap>

class DockPluginController;
class PluginsItemInterface;

class IconManager : public QObject
{
    Q_OBJECT

public:
    explicit IconManager(DockPluginController *pluginController, QObject *parent = nullptr);
    void setPosition(Dock::Position position);
    void setDisplayMode(Dock::DisplayMode displayMode);
    QPixmap pixmap(DGuiApplicationHelper::ColorType colorType) const;

private:
    PluginsItemInterface *findPlugin(const QString &pluginName) const;

private:
    DockPluginController *m_pluginController;
    Dock::Position m_position;
    Dock::DisplayMode m_displayMode;
};

#endif // ICONMANAGER_H
