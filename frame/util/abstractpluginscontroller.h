/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

class AbstractPluginsController : public QObject, PluginProxyInterface
{
    Q_OBJECT

public:
    explicit AbstractPluginsController(QObject *parent = Q_NULLPTR);
    ~ AbstractPluginsController() override;

    void updateDockInfo(PluginsItemInterface *const, const DockPart &) override {}

    virtual bool needLoad(PluginsItemInterface *) { return true; }
    QMap<PluginsItemInterface *, QMap<QString, QObject *>> &pluginsMap();

    virtual void savePluginValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant &value);
    virtual const QVariant getPluginValue(PluginsItemInterface *const itemInter, const QString &key, const QVariant& fallback = QVariant());
    virtual void removePluginValue(PluginsItemInterface * const itemInter, const QStringList &keyList);

    virtual void pluginItemAdded(PluginsItemInterface * const, const QString &) = 0;
    virtual void pluginItemUpdate(PluginsItemInterface * const, const QString &) = 0;
    virtual void pluginItemRemoved(PluginsItemInterface * const, const QString &) = 0;
    virtual void requestPluginWindowAutoHide(PluginsItemInterface * const, const QString &, const bool) {}
    virtual void requestRefreshPluginWindowVisible(PluginsItemInterface * const, const QString &) {}
    virtual void requestSetPluginAppletVisible(PluginsItemInterface * const, const QString &, const bool) {}

Q_SIGNALS:
    void pluginLoaderFinished();

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

protected:
    QObject *pluginItemAt(PluginsItemInterface * const itemInter, const QString &itemKey) const;
    PluginsItemInterface *pluginInterAt(const QString &itemKey);
    PluginsItemInterface *pluginInterAt(QObject *destItem);
    bool eventFilter(QObject *o, QEvent *e) override;

protected Q_SLOTS:
    void startLoader(PluginLoader *loader);

private slots:
    void displayModeChanged();
    void positionChanged();
    void loadPlugin(const QString &pluginFile);
    void initPlugin(PluginsItemInterface *interface);
    void refreshPluginSettings();

private:
    QDBusConnectionInterface *m_dbusDaemonInterface;
    DockInter *m_dockDaemonInter;

    // interface,  "pluginloader", PluginLoader指针对象
    QMap<PluginsItemInterface *, QMap<QString, QObject *>> m_pluginsMap;

    // filepath, interface, loaded
    QMap<QPair<QString, PluginsItemInterface *>, bool> m_pluginLoadMap;

    QJsonObject m_pluginSettingsObject;
    QMap<qulonglong, PluginAdapter *> m_pluginAdapterMap;
};

#endif // ABSTRACTPLUGINSCONTROLLER_H
