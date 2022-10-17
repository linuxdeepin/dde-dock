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
#include "constants.h"

#include <QSettings>

static QStringList getPathFromConf(const QString &key) {
    QSettings set("/etc/deepin/dde-dock.conf", QSettings::IniFormat);
    auto value = set.value(key).toString();
    if (!value.isEmpty()) {
        return value.split(':');
    }

    return QStringList();
}

// 该方法用来设置所有的需要加载的插件的路径信息
static QMap<PluginType, QStringList> getPluginPaths()
{
    QMap<PluginType, QStringList> plugins;
    // 添加快捷插件目录
    {
        QStringList pluginPaths;
    #ifdef QT_DEBUG
        pluginPaths << QString("%1/..%2").arg(qApp->applicationDirPath()).arg(QUICK_PATH)
                    << QString("%1/..%2").arg(qApp->applicationDirPath()).arg(PLUGIN_PATH)
                    << QString("%1/..%2").arg(qApp->applicationDirPath()).arg(TRAY_PATH);
    #else
        pluginPaths << QString("/usr/lib/dde-dock%1").arg(QUICK_PATH)
                    << QString("/usr/lib/dde-dock%1").arg(PLUGIN_PATH)
                    << QString("/usr/lib/dde-dock%1").arg(TRAY_PATH);

        const QStringList pluginsDirs = (getPathFromConf("QUICK_TRAY_PATH") << getPathFromConf("PATH") << getPathFromConf("SYSTEM_TRAY_PATH"));
        if (!pluginsDirs.isEmpty())
            pluginPaths << pluginsDirs;
    #endif
        plugins[PluginType::QuickPlugin] = pluginPaths;
    }

    // 添加系统插件目录
    {
        QStringList pluginPaths;
    #ifdef QT_DEBUG
        pluginPaths << qApp->applicationDirPath() + "/../plugins/system-trays";
    #else
        pluginPaths << "/usr/lib/dde-dock/plugins/system-trays";

        const QStringList &pluginsDirs = getPathFromConf("SYSTEM_TRAY_PATH");
        if (!pluginsDirs.isEmpty())
            pluginPaths << pluginsDirs;
    #endif
        plugins[PluginType::SystemTrays] = pluginPaths;
    }

    return plugins;
}

// 该方法根据当前加载插件的类型来生成对应的单例的类
ProxyPluginController *ProxyPluginController::instance(PluginType instanceKey)
{
    // 此处将这些单例对象存储到了qApp里面，而没有存储到本地的静态变量是因为这个对象会在dock进程和tray插件中同时调用，
    // 如果存储到内存的临时变量中，他们就是不同的内存地址，获取到的变量就是多个，这样就会导致相同的插件加载多次，
    // 而qApp是dock和插件共用的，因此将对象存储到这里是保证能获取到相同的指针对象
    QMap<PluginType, ProxyPluginController *> proxyInstances = qApp->property("proxyController").value<QMap<PluginType, ProxyPluginController *>>();
    if (proxyInstances.contains(instanceKey))
        return proxyInstances.value(instanceKey);

    // 生成单例类，获取加载插件的路径信息
    static QMap<PluginType, QStringList> pluginLoadInfos = getPluginPaths();
    ProxyPluginController *controller = new ProxyPluginController();
    controller->m_dirs = pluginLoadInfos.contains(instanceKey) ? pluginLoadInfos[instanceKey] : QStringList();
    proxyInstances[instanceKey] = controller;
    qApp->setProperty("proxyController", QVariant::fromValue(proxyInstances));
    return controller;
}

ProxyPluginController *ProxyPluginController::instance(PluginsItemInterface *itemInter)
{
    // 根据插件指针获取对应的代理对象，因为在监听者里可能存在同时加载多个不同目录的插件，用到的就是多实例，
    // 添加插件的时候，不知道当前插件是属于哪个实例，因此在此处添加获取对应插件的实例，方便监听者拿到正确的实例
    QVariant proxyProperty = qApp->property("proxyController");
    if (!proxyProperty.canConvert<QMap<PluginType, ProxyPluginController *>>())
        return nullptr;

    QMap<PluginType, ProxyPluginController *> proxyControllers = proxyProperty.value<QMap<PluginType, ProxyPluginController *>>();
    for (ProxyPluginController *proxyController : proxyControllers) {
        const QList<PluginsItemInterface *> &pluginItems = proxyController->m_pluginsItems;
        for (PluginsItemInterface *interPair : pluginItems) {
            if (interPair == itemInter)
                return proxyController;
        }
    }

    return nullptr;
}

// 新增要使用的控制器，第二个参数表示当前控制器需要加载的插件名称，为空表示加载所有插件
void ProxyPluginController::addProxyInterface(AbstractPluginsController *interface)
{
    if (!m_interfaces.contains(interface))
        m_interfaces << interface;
}

void ProxyPluginController::removeProxyInterface(AbstractPluginsController *interface)
{
    Q_ASSERT(m_interfaces.contains(interface));
    m_interfaces.removeOne(interface);
}

ProxyPluginController::ProxyPluginController(QObject *parent)
    : AbstractPluginsController(parent)
{
    // 只有在非安全模式下才加载插件，安全模式会在等退出安全模式后通过接受事件的方式来加载插件
    if (!qApp->property("safeMode").toBool())
        QMetaObject::invokeMethod(this, &ProxyPluginController::startLoader, Qt::QueuedConnection);
}

QPluginLoader *ProxyPluginController::pluginLoader(PluginsItemInterface * const itemInter)
{
    QMap<PluginsItemInterface *, QMap<QString, QObject *> > &plugin = pluginsMap();
    if (plugin.contains(itemInter))
        return qobject_cast<QPluginLoader *>(plugin[itemInter].value("pluginloader"));

    return nullptr;
}

QList<PluginsItemInterface *> ProxyPluginController::pluginsItems() const
{
    return m_pluginsItems;
}

QString ProxyPluginController::itemKey(PluginsItemInterface *itemInter) const
{
    if (m_pluginsItemKeys.contains(itemInter))
        return m_pluginsItemKeys.value(itemInter);

    return QString();
}

void ProxyPluginController::pluginItemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    if (!addPluginItems(itemInter, itemKey))
        return;

    // 获取需要加载当前插件的监听者,然后将当前插件添加到监听者
    QList<AbstractPluginsController *> validController = getValidController(itemInter);
    for (AbstractPluginsController *interface : validController)
        interface->pluginItemAdded(itemInter, itemKey);
}

void ProxyPluginController::pluginItemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    QList<AbstractPluginsController *> validController = getValidController(itemInter);
    for (AbstractPluginsController *interface : validController)
        interface->pluginItemUpdate(itemInter, itemKey);
}

void ProxyPluginController::pluginItemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    // 先获取可执行的controller，再移除，因为在判断当前插件是否加载的时候需要用到当前容器中的插件来获取当前代理
    QList<AbstractPluginsController *> validController = getValidController(itemInter);
    for (AbstractPluginsController *interface : validController)
        interface->pluginItemRemoved(itemInter, itemKey);

    removePluginItem(itemInter);
}

void ProxyPluginController::requestPluginWindowAutoHide(PluginsItemInterface * const itemInter, const QString &itemKey, const bool autoHide)
{
    QList<AbstractPluginsController *> validController = getValidController(itemInter);
    for (AbstractPluginsController *interface : validController)
        interface->requestPluginWindowAutoHide(itemInter, itemKey, autoHide);
}

void ProxyPluginController::requestRefreshPluginWindowVisible(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    QList<AbstractPluginsController *> validController = getValidController(itemInter);
    for (AbstractPluginsController *interface : validController)
        interface->requestRefreshPluginWindowVisible(itemInter, itemKey);
}

void ProxyPluginController::requestSetPluginAppletVisible(PluginsItemInterface * const itemInter, const QString &itemKey, const bool visible)
{
    QList<AbstractPluginsController *> validController = getValidController(itemInter);
    for (AbstractPluginsController *interface : validController)
        interface->requestSetPluginAppletVisible(itemInter, itemKey, visible);
}

void ProxyPluginController::updateDockInfo(PluginsItemInterface * const itemInter, const DockPart &part)
{
    QList<AbstractPluginsController *> validController = getValidController(itemInter);
    for (AbstractPluginsController *interface : validController)
        interface->updateDockInfo(itemInter, part);
}

bool ProxyPluginController::eventFilter(QObject *watched, QEvent *event)
{
    if (watched == qApp && event->type() == PluginLoadEvent::eventType()) {
        // 如果收到的是重新加载插件的消息（一般是在退出安全模式后），则直接加载插件即可
        startLoader();
    }

    return AbstractPluginsController::eventFilter(watched, event);
}

QList<AbstractPluginsController *> ProxyPluginController::getValidController(PluginsItemInterface *itemInter) const
{
    QList<AbstractPluginsController *> validController;
    for (AbstractPluginsController *interface : m_interfaces) {
        if (!interface->needLoad(itemInter))
            continue;

        validController << interface;
    }

    return validController;
}

bool ProxyPluginController::addPluginItems(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    // 如果该插件已经存在，则无需再次插入
    if (m_pluginsItems.contains(itemInter))
        return false;

    m_pluginsItems << itemInter;
    m_pluginsItemKeys[itemInter] = itemKey;
    return true;
}

void ProxyPluginController::removePluginItem(PluginsItemInterface * const itemInter)
{
    if (m_pluginsItems.contains(itemInter))
        m_pluginsItems.removeOne(itemInter);

    if (m_pluginsItemKeys.contains(itemInter))
        m_pluginsItemKeys.remove(itemInter);
}

void ProxyPluginController::startLoader()
{
    QDir dir;
    for (const QString &path : m_dirs) {
        if (!dir.exists(path))
            continue;

        AbstractPluginsController::startLoader(new PluginLoader(path, this));
    }
}

// 注册事件类型
static QEvent::Type pluginEventType = (QEvent::Type)QEvent::registerEventType(QEvent::User + 1001);

// 事件处理，当收到该事件的时候，加载插件
PluginLoadEvent::PluginLoadEvent()
    : QEvent(pluginEventType)
{
}

PluginLoadEvent::~PluginLoadEvent()
{
}

QEvent::Type PluginLoadEvent::eventType()
{
    return pluginEventType;
}