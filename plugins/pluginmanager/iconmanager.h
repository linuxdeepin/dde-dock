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
    void updateSize(QSize size);
    void setPosition(Dock::Position position);
    void setDisplayMode(Dock::DisplayMode displayMode);
    QPixmap pixmap(DGuiApplicationHelper::ColorType colorType) const;
    bool isFixedPlugin(PluginsItemInterface *plugin) const;

private:
    PluginsItemInterface *findPlugin(const QString &pluginName) const;

private:
    DockPluginController 