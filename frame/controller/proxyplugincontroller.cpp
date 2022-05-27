/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
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
#include "proxyplugincontroller.h"
#include "pluginsiteminterface.h"

QMap<int, ProxyPluginController *> ProxyPluginController::m_instances = {};

// 该方法用来设置所有的需要加载的插件的路径信息
static QMap<int, QList<QStringList>> getPluginPaths()
{
    QList<QStringList> pluginPaths;
    pluginPaths << QStringList{ QString("%1/.local/lib/dde-dock/plugins/").arg(QDir::homePath()) }
                << QStringList{ QString(qApp->applicationDirPath() + "/../plugins"),
                                QString("/usr/lib/dde-dock/plugins") };
    QMap<int, QList<QStringList>> plugins;
    plugins[FIXEDSYSTEMPLUGIN] = pluginPaths;
    return plugins;
}

// 该方法根据当前加载插件的类型来生成对应的单例的类
ProxyPluginController *ProxyPluginController::instance(int instanceKey)
{
    static QMap<int, QList<QStringList>> pluginLoadInfos = getPluginPaths();

    if (m_instances.contains(instanceKey))
        return m_instances.value(instanceKey);

    // 生成单例类，获取加载插件的路径信息
    ProxyPluginController *controller = new ProxyPluginController();
    controller->m_dirs = (pluginLoadInfos.contains(instanceKey) ? pluginLoadInfos[instanceKey] : QList<QStringList>());
    m_instances[instanceKey] = controller;
    return controller;
}

// 新增要使用的控制器，第二个参数表示当前控制器需要加载的插件名称，为空表示加载所有插件
void ProxyPluginController::addProxyInterface(AbstractPluginsController *interface, const QStringList &pluginNames)
{
    if (!m_interfaces.contains(interface))
        m_interfaces[interface] = pluginNames;
}

void ProxyPluginController::removeProxyInterface(AbstractPluginsController *interface)
{
    if (m_interfaces.contains(interface))
        m_interfaces.remove(interface);
}

ProxyPluginController::ProxyPluginController(QObject *parent)
    : AbstractPluginsController(parent)
{
}

QPluginLoader *ProxyPluginController::pluginLoader(PluginsItemInterface * const itemInter)
{
    QMap<PluginsItemInterface *, QMap<QString, QObject *> > &plugin = pluginsMap();
    if (plugin.contains(itemInter))
        return qobject_cast<QPluginLoader *>(plugin[itemInter].value("pluginloader"));

    return nullptr;
}

void ProxyPluginController::itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    // 只有当前的controll设置的过滤名称包含当前插件的名称或者过滤名称为空，才新增当前插件
    QList<AbstractPluginsController *> pluginKeys = m_interfaces.keys();
    for (AbstractPluginsController *interface: pluginKeys) {
        const QStringList &filterNames = m_interfaces[interface];
        if (filterNames.isEmpty() || filterNames.contains(itemInter->pluginName()))
            interface->itemAdded(itemInter, itemKey);
    }
}

void ProxyPluginController::itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    QList<AbstractPluginsController *> pluginKeys = m_interfaces.keys();
    for (AbstractPluginsController *interface: pluginKeys) {
        const QStringList &filterNames = m_interfaces[interface];
        if (filterNames.isEmpty() || filterNames.contains(itemInter->pluginName()))
            interface->itemUpdate(itemInter, itemKey);
    }
}

void ProxyPluginController::itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    QList<AbstractPluginsController *> pluginKeys = m_interfaces.keys();
    for (AbstractPluginsController *interface: pluginKeys) {
        const QStringList &filterNames = m_interfaces[interface];
        if (filterNames.isEmpty() || filterNames.contains(itemInter->pluginName()))
            interface->itemRemoved(itemInter, itemKey);
    }
}

void ProxyPluginController::startLoader()
{
    QDir dir;
    for (const QStringList &pluginPaths : m_dirs) {
        QString loadPath;
        for (const QString &pluginPath : pluginPaths) {
            if (!dir.exists(pluginPath))
                continue;

            loadPath = pluginPath;
            break;
        }

        if (!loadPath.isEmpty())
            AbstractPluginsController::startLoader(new PluginLoader(loadPath, this));
    }
}
