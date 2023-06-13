// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef DOCKPLUGINCONTROLLER_H
#define DOCKPLUGINCONTROLLER_H

#include "pluginproxyinterface.h"
#include "pluginloader.h"
#include "dbusutil.h"

#include <QList>
#include <QMap>
#include <QDBusConnectionInterface>

class PluginsItemInterface;
class PluginAdapter;

class DockPluginController : public QObject, protected PluginProxyInterface
{
    Q_OBJECT

public:
    explicit DockPluginController(PluginProxyInterface *proxyInter, QObject *parent = Q_NULLPTR);
    ~ DockPluginController() override;

    QList<PluginsItemInterface *> plugins() const;              // 所有的插件
    QList<PluginsItemInterface *> pluginsInSetting() const;     // 控制中心用于可以设置是否显示或隐藏的插件
    QList<PluginsItemInterface *> currentPlugins() const;       // 当前使用的插件

    QString itemKey(PluginsItemInterface *itemInter) const;
    QJsonObject metaData(PluginsItemInterface *pluginItem);

    virtual void savePluginValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant &value);
    virtual const QVariant getPluginValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant& fallback = QVariant());
    virtual void removePluginValue(PluginsItemInterface * const itemInter, const QStringList &keyList);
    void startLoadPlugin(const QStringList &dirs);

Q_SIGNALS:
    void pluginLoadFinished();
    void pluginInserted(PluginsItemInterface *itemInter, QString);
    void pluginRemoved(PluginsItemInterface *itemInter);
    void pluginUpdated(PluginsItemInterface *, const DockPart);
    void requestAppletVisible(PluginsItemInterface *, const QString &, bool);

protected:
    bool isPluginLoaded(PluginsItemInterface *itemInter);

    QObject *pluginItemAt(PluginsItemInterface * const itemInter, const QString &itemKey) const;
    PluginsItemInterface *pluginInterAt(const QString &itemKey);
    PluginsItemInterface *pluginInterAt(QObject *destItem);
    bool eventFilter(QObject *object, QEvent *event) override;

    bool pluginCanDock(PluginsItemInterface *plugin) const;
    bool pluginCanDock(const QStringList &config, PluginsItemInterface *plugin) const;
    void updateDockInfo(PluginsItemInterface * const itemInter, const DockPart &part) override;

private:
    // implements PluginProxyInterface
    void saveValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant &value) override;
    const QVariant getValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant& fallback = QVariant()) override;
    void removeValue(PluginsItemInterface * const itemInter, const QStringList &keyList) override;

    void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void requestWindowAutoHide(PluginsItemInterface * const itemInter, const QString &itemKey, const bool autoHide) override;
    void requestRefreshWindowVisible(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void requestSetAppletVisible(PluginsItemInterface * const itemInter, const QString &itemKey, const bool visible) override;
    PluginsItemInterface *getPluginInterface(PluginsItemInterface * const itemInter);

    void addPluginItem(PluginsItemInterface * const itemInter, const QString &itemKey);
    void removePluginItem(PluginsItemInterface * const itemInter, const QString &itemKey);

private Q_SLOTS:
    void startLoader(PluginLoader *loader);
    void displayModeChanged();
    void positionChanged();
    void loadPlugin(const QString &pluginFile);
    void initPlugin(PluginsItemInterface *interface);
    void refreshPluginSettings();
    void onConfigChanged(const QStringList &pluginNames);

private:
    QDBusConnectionInterface *m_dbusDaemonInterface;

    // interface,  "pluginloader", PluginLoader指针对象
    QMap<PluginsItemInterface *, QMap<QString, QObject *>> m_pluginsMap;

    // filepath, interface, loaded
    QMap<QPair<QString, PluginsItemInterface *>, bool> m_pluginLoadMap;

    QJsonObject m_pluginSettingsObject;
    QMap<qulonglong, PluginAdapter *> m_pluginAdapterMap;

    PluginProxyInterface *m_proxyInter;
};

#endif // ABSTRACTPLUGINSCONTROLLER_H
