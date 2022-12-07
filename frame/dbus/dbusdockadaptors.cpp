// SPDX-FileCopyrightText: 2016 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "dbusdockadaptors.h"
#include "utils.h"
#include "dockitemmanager.h"

#include <QScreen>
#include <QDebug>
#include <QGSettings>

DBusDockAdaptors::DBusDockAdaptors(MainWindow* parent)
    : QDBusAbstractAdaptor(parent)
    , m_gsettings(Utils::SettingsPtr("com.deepin.dde.dock.mainwindow", QByteArray(), this))
{
    connect(parent, &MainWindow::panelGeometryChanged, this, [=] {
        emit DBusDockAdaptors::geometryChanged(geometry());

        // 手动触发PropertiesChanged信号，可能是qt不支持属性自动触发
        QDBusMessage msg = QDBusMessage::createSignal("/com/deepin/dde/Dock", "org.freedesktop.DBus.Properties", "PropertiesChanged");
        QList<QVariant> arguments;
        arguments.push_back("com.deepin.dde.Dock");
        QVariantMap changedProps;
        changedProps.insert("geometry", geometry());
        arguments.push_back(changedProps);
        msg.setArguments(arguments);
        QDBusConnection::connectToBus(QDBusConnection::SessionBus, "com.deepin.dde.Dock").send(msg);
    });

    if (m_gsettings) {
        connect(m_gsettings, &QGSettings::changed, this, [ = ] (const QString &key) {
            if (key == "onlyShowPrimary") {
                Q_EMIT showInPrimaryChanged(m_gsettings->get(key).toBool());
            }
        });
    }

    connect(DockItemManager::instance(), &DockItemManager::itemInserted, this, [ = ] (const int index, DockItem *item) {
        Q_UNUSED(index);
        if (item->itemType() == DockItem::Plugins
                || item->itemType() == DockItem::FixedPlugin) {
            PluginsItem *pluginItem = static_cast<PluginsItem *>(item);
            for (auto *p : DockItemManager::instance()->pluginList()) {
                if (p->pluginName() == pluginItem->pluginName()) {
                    Q_EMIT pluginVisibleChanged(p->pluginDisplayName(), getPluginVisible(p->pluginDisplayName()));
                }
            }
        }
    });

    connect(DockItemManager::instance(), &DockItemManager::itemRemoved, this, [ = ] (DockItem *item) {
        if (item->itemType() == DockItem::Plugins
                || item->itemType() == DockItem::FixedPlugin) {
            PluginsItem *pluginItem = static_cast<PluginsItem *>(item);
            for (auto *p : DockItemManager::instance()->pluginList()) {
                if (p->pluginName() == pluginItem->pluginName()) {
                    Q_EMIT pluginVisibleChanged(p->pluginDisplayName(), getPluginVisible(p->pluginDisplayName()));
                }
            }
        }
    });
}

DBusDockAdaptors::~DBusDockAdaptors()
{

}

MainWindow *DBusDockAdaptors::parent() const
{
    return static_cast<MainWindow *>(QObject::parent());
}

void DBusDockAdaptors::callShow()
{
    return parent()->callShow();
}

void DBusDockAdaptors::ReloadPlugins()
{
    return parent()->relaodPlugins();
}

QStringList DBusDockAdaptors::GetLoadedPlugins()
{
    auto pluginList = DockItemManager::instance()->pluginList();
    QStringList nameList;
    QMap<QString, QString> map;
    for (auto plugin : pluginList) {
        // 托盘本身也是一个插件，这里去除掉这个特殊的插件,还有一些没有实际名字的插件
        if (plugin->pluginName() == "tray"
                || plugin->pluginDisplayName().isEmpty()
                || !isPluginValid(plugin->pluginName()))
            continue;

        nameList << plugin->pluginName();
        map.insert(plugin->pluginName(), plugin->pluginDisplayName());
    }

    // 排序,保持和原先任务栏右键菜单中的插件列表顺序一致
    qSort(nameList.begin(), nameList.end(), [ = ] (const QString &name1, const QString &name2) {
        return name1 > name2;
    });

    QStringList newList;
    for (auto name : nameList) {
        newList.push_back(map[name]);
    }

    return newList;
}

void DBusDockAdaptors::resizeDock(int offset, bool dragging)
{
    parent()->resizeDock(offset, dragging);
}

// 返回每个插件的识别Key(所以此值应始终不变)，供个性化插件根据key去匹配每个插件对应的图标
QString DBusDockAdaptors::getPluginKey(const QString &pluginName)
{
    for (auto plugin : DockItemManager::instance()->pluginList()) {
        if (plugin->pluginDisplayName() == pluginName)
            return plugin->pluginName();
    }

    return QString();
}

bool DBusDockAdaptors::getPluginVisible(const QString &pluginName)
{
    for (auto *p : DockItemManager::instance()->pluginList()) {
        if (!p->pluginIsAllowDisable())
            continue;

        const QString &display = p->pluginDisplayName();
        if (display != pluginName)
            continue;

        const QString &name = p->pluginName();
        if (!isPluginValid(name))
            continue;

        return !p->pluginIsDisable();
    }

    qInfo() << "Unable to get information about this plugin";
    return false;
}

void DBusDockAdaptors::setPluginVisible(const QString &pluginName, bool visible)
{
    for (auto *p : DockItemManager::instance()->pluginList()) {
        if (!p->pluginIsAllowDisable())
            continue;

        const QString &display = p->pluginDisplayName();
        if (display != pluginName)
            continue;

        const QString &name = p->pluginName();
        if (!isPluginValid(name))
            continue;

        if (p->pluginIsDisable() == visible) {
            p->pluginStateSwitched();
            Q_EMIT pluginVisibleChanged(pluginName, visible);
        }
        return;
    }

    qInfo() << "Unable to set information for this plugin";
}

QRect DBusDockAdaptors::geometry() const
{
    return parent()->geometry();
}

bool DBusDockAdaptors::showInPrimary() const
{
    return Utils::SettingValue("com.deepin.dde.dock.mainwindow", QByteArray(), "onlyShowPrimary", false).toBool();
}

void DBusDockAdaptors::setShowInPrimary(bool showInPrimary)
{
    if (this->showInPrimary() == showInPrimary)
        return;

    if (Utils::SettingSaveValue("com.deepin.dde.dock.mainwindow", QByteArray(), "onlyShowPrimary", showInPrimary)) {
        Q_EMIT showInPrimaryChanged(showInPrimary);
    }
}

bool DBusDockAdaptors::isPluginValid(const QString &name)
{
    // 插件被全局禁用时，理应获取不到此插件的任何信息
    if (!Utils::SettingValue("com.deepin.dde.dock.module." + name, QByteArray(), "enable", true).toBool())
        return false;

    // 在非wayland环境下，未开启窗口特效时，不显示多任务视图插件
    if (name == "multitasking" && !DWindowManagerHelper::instance()->hasComposite() && !Utils::IS_WAYLAND_DISPLAY)
        return false;

    // 录屏插件不显示,插件名如果有变化，建议发需求，避免任务栏反复适配
    if (name == "deepin-screen-recorder-plugin")
        return false;

    // 远程助手只有在启动时才会加载在任务栏，不需要在控制中心加载
    if (name == "uos-remote-assistance")
        return false;

    return true;
}
