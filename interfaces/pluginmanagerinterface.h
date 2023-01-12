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
#ifndef PLUGINMANAGERINTERFACE_H
#define PLUGINMANAGERINTERFACE_H

#include <QObject>
#include <QJsonObject>

class PluginsItemInterface;

class PluginManagerInterface : public QObject
{
    Q_OBJECT

public:
    virtual QList<PluginsItemInterface *> plugins() const = 0;
    virtual QList<PluginsItemInterface *> pluginsInSetting() const = 0;
    virtual QList<PluginsItemInterface *> currentPlugins() const = 0;
    virtual QString itemKey(PluginsItemInterface *itemInter) const = 0;
    virtual QJsonObject metaData(PluginsItemInterface *itemInter) const = 0;

Q_SIGNALS:
    void pluginLoadFinished();
};

#endif // PLUGINMANAGERINTERFACE_H
