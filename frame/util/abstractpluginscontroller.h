// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef ABSTRACTPLUGINSCONTROLLER_H
#define ABSTRACTPLUGINSCONTROLLER_H

#include "pluginproxyinterface.h"
#include "pluginloader.h"
#include "dbusutil.h"

#include <QPluginLoader>
#include <QList>
#include <QMap>
#include <QDBusConnectionInterface>

class PluginsItemInterface;
class PluginAdapter;
class PluginManagerInterface;

class AbstractPluginsController : public QObject, PluginProxyInterface
{
    Q_OBJECT

public:
    explicit AbstractPluginsController(QObject *parent = Q_NULLPTR);
    ~ AbstractPluginsController() override;

    void updateDockInfo(PluginsItemInterface *const, const DockPart &) override {}

Q_SIGNALS:
    void pluginLoaderFinished();

protected:
    bool eventFilter(QObject *object, QEvent *event) override;

    PluginManagerInterface *pluginManager() const;

private:
    // implements PluginProxyInterface
    void requestWindowAutoHide(PluginsItemInterface * const itemInter, const QString &itemKey, const bool autoHide) override {}
    void requestRefreshWindowVisible(PluginsItemInterface * const itemInter, const QString &itemKey) override {}
    void saveValue(PluginsItemInterface * const itemInter, const QString &key, const QVariant &value) override {}
    void removeValue(PluginsItemInterface *const itemInter, const QStringList &keyList) override {}
    const QVariant getValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant& fallback = QVariant()) override { return QVariant(); }

protected Q_SLOTS:
    void startLoader(PluginLoader *loader);

private slots:
    void displayModeChanged();
    void positionChanged();
    void loadPlugin(const QString &pluginFile);
    void initPlugin(PluginsItemInterface *interface);

private:
    QMap<PluginsItemInterface *, QMap<QString, QObject *>> m_pluginsMap;

    // filepath, interface, loaded
    QMap<QPair<QString, PluginsItemInterface *>, bool> m_pluginLoadMap;

    QJsonObject m_pluginSettingsObject;
    PluginManagerInterface *m_pluginManager;
};

#endif // ABSTRACTPLUGINSCONTROLLER_H
