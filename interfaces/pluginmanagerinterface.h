// Copyright (C) 2023 ~ 2023 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
