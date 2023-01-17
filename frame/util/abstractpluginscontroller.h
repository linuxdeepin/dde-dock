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
