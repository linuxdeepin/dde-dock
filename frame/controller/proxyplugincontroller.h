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
// 加载的插件的类型(1 根目录下的插件 2 快捷设置插件 3 系统插件)
enum class PluginType {
    QuickPlugin = 0,
    SystemTrays
};

// 该类是一个底层用来加载系统插件的类，DockPluginsController和
// FixedPluginController类都是通过这个类来加载系统插件的
// 该类做成一个多例，因为理论上一个插件只允许被加载一次，但是对于电源插件来说，
// 电源插件在高效模式和特效模式下都需要显示，因此，此类用于加载插件，然后分发到不同的
// 上层控制器中
class ProxyPluginController : public AbstractPluginsController
{
    Q_OBJECT

public:
    static ProxyPluginController *instance(PluginType instanceKey);
    static ProxyPluginController *instance(PluginsItemInterface *itemInter);
    void addProxyInterface(AbstractPluginsController *interface);
    void removeProxyInterface(AbstractPluginsController *interface);
    QPluginLoader *pluginLoader(PluginsItemInterface * const itemInter);
    QList<PluginsItemInterface *> pluginsItems() const;
    QString itemKey(PluginsItemInterface *itemInter) const;

protected:
    explicit ProxyPluginController(QObject *parent = nullptr);
    ~ProxyPluginController() override {}

    void pluginItemAdded(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void pluginItemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void pluginItemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey) override;

    void requestPluginWindowAutoHide(PluginsItemInterface * const itemInter, const QString &itemKey, const bool autoHide) override;
    void requestRefreshPluginWindowVisible(PluginsItemInterface * const itemInter, const QString &itemKey) override;
    void requestSetPluginAppletVisible(PluginsItemInterface * const itemInter, const QString &itemKey, const bool visible) override;

    void updateDockInfo(PluginsItemInterface * const itemInter, const DockPart &part) override;

    bool eventFilter(QObject *watched, QEvent *event) override;

private:
    QList<AbstractPluginsController *> getValidController(PluginsItemInterface *itemInter) const;
    bool addPluginItems(PluginsItemInterface * const itemInter, const QString &itemKey);
    void removePluginItem(PluginsItemInterface * const itemInter);
    void startLoader();

private:
    QList<AbstractPluginsController *> m_interfaces;
    QStringList m_dirs;
    QList<PluginsItemInterface *> m_pluginsItems;
    QMap<PluginsItemInterface *, QString> m_pluginsItemKeys;
};

// 该插件用于处理插件的延迟加载，当退出安全模式后，会收到该事件并加载插件
class PluginLoadEvent : public QEvent
{
public:
    PluginLoadEvent();
    ~PluginLoadEvent() override;

    static Type eventType();
};

Q_DECLARE_METATYPE(PluginType)

#endif // PROXYPLUGINCONTROLLER_H
