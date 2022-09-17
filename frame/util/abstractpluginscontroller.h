// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef ABSTRACTPLUGINSCONTROLLER_H
#define ABSTRACTPLUGINSCONTROLLER_H

#include "pluginproxyinterface.h"
#include "pluginloader.h"

#include <com_deepin_dde_daemon_dock.h>

#include <QPluginLoader>
#include <QList>
#include <QMap>
#include <QDBusConnectionInterface>

using DockDaemonInter = com::deepin::dde::daemon::Dock;

class PluginsItemInterface;
class AbstractPluginsController : public QObject, PluginProxyInterface
{
    Q_OBJECT

public:
    explicit AbstractPluginsController(QObject *parent = 0);
    ~ AbstractPluginsController() override;

    // implements PluginProxyInterface
    void saveValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant &value) override;
    const QVariant getValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant& fallback = QVariant()) override;
    void removeValue(PluginsItemInterface * const itemInter, const QStringList &keyList) override;

signals:
    void pluginLoaderFinished();

protected:
    QMap<PluginsItemInterface *, QMap<QString, QObject *>> &pluginsMap();
    QObject *pluginItemAt(PluginsItemInterface * const itemInter, const QString &itemKey) const;
    PluginsItemInterface *pluginInterAt(const QString &itemKey);
    PluginsItemInterface *pluginInterAt(QObject *destItem);

protected Q_SLOTS:
    void startLoader(PluginLoader *loader);

private slots:
    void displayModeChanged();
    void positionChanged();
    void loadPlugin(const QString &pluginFile);
    void initPlugin(PluginsItemInterface *interface);
    void refreshPluginSettings();

private:
    bool eventFilter(QObject *o, QEvent *e) override;

private:
    QDBusConnectionInterface *m_dbusDaemonInterface;
    DockDaemonInter *m_dockDaemonInter;

    // interface,  "pluginloader", PluginLoader指针对象
    QMap<PluginsItemInterface *, QMap<QString, QObject *>> m_pluginsMap;

    // filepath, interface, loaded
    QMap<QPair<QString, PluginsItemInterface *>, bool> m_pluginLoadMap;

    QJsonObject m_pluginSettingsObject;
};

#endif // ABSTRACTPLUGINSCONTROLLER_H
