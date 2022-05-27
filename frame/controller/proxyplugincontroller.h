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
#ifndef PROXYPLUGINCONTROLLER_H
#define PROXYPLUGINCONTROLLER_H

#include "abstractpluginscontroller.h"

class PluginsItemInterface;
// 加载的插件的类型
#define FIXEDSYSTEMPLUGIN 1

// 该类是一个底层用来加载系统插件的类，DockPluginsController和
// FixedPluginController类都是通过这个类来加载系统插件的
// 该类做成一个单例，因为理论上一个插件只允许被加载一次，但是对于电源插件来说，
// 电源插件在高效模式和特效模式下都需要显示，因此，此类用于加载插件，然后分发到不同的
// 上层控制器中
class ProxyPluginController : public AbstractPluginsController
{
    Q_OBJECT

public:
    static ProxyPluginController *instance(int instanceKey = FIXEDSYSTEMPLUGIN);
    void addProxyInterface(AbstractPluginsController *interface, const QStringList &pluginNames = QStringList());
    void removeProxyInterface(AbstractPluginsController *interface);
    void startLoader();
    QPluginLoader *pluginLoader(PluginsItemInterface * const itemInter);

protected:
    explicit ProxyPluginController(QObject *parent = nullptr);
    ~ProxyPluginController() override {}

    void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey) override;

    void requestWindowAutoHide(PluginsItemInterface * const, const QString &, const bool) override {}
    void requestRefreshWindowVisible(PluginsItemInterface * const, const QString &) override {}
    void requestSetAppletVisible(PluginsItemInterface * const, const QString &, const bool) override {}

private:
    static QMap<int, ProxyPluginController *> m_instances;
    QMap<AbstractPluginsController *, QStringList> m_interfaces;
    QList<QStringList> m_dirs;
};

#endif // PROXYPLUGINCONTROLLER_H
