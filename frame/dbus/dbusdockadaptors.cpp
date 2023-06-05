/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

#include "dbusdockadaptors.h"
#include "utils.h"
#include "dockitemmanager.h"
#include "windowmanager.h"
#include "quicksettingcontroller.h"
#include "pluginsitem.h"
#include "docksettings.h"
#include "common.h"
#include "customevent.h"

#include <DGuiApplicationHelper>

#include <QScreen>
#include <QDebug>
#include <QGSettings>
#include <QDBusMetaType>

const QSize defaultIconSize = QSize(20, 20);

QDebug operator<<(QDebug argument, const DockItemInfo &info)
{
    argument << "name:" << info.name << ", displayName:" << info.displayName
             << "itemKey:" << info.itemKey << "SettingKey:" << info.settingKey
             << "icon_light:" << info.iconLight << "icon_dark:" << info.iconDark << "visible:" << info.visible;
    return argument;
}

QDBusArgument &operator<<(QDBusArgument &arg, const DockItemInfo &info)
{
    arg.beginStructure();
    arg << info.name << info.displayName << info.itemKey << info.settingKey << info.iconLight << info.iconDark << info.visible;
    arg.endStructure();
    return arg;
}

const QDBusArgument &operator>>(const QDBusArgument &arg, DockItemInfo &info)
{
    arg.beginStructure();
    arg >> info.name >> info.displayName >> info.itemKey >> info.settingKey >> info.iconLight >> info.iconDark >> info.visible;
    arg.endStructure();
    return arg;
}

void registerPluginInfoMetaType()
{
    qRegisterMetaType<DockItemInfo>("DockItemInfo");
    qDBusRegisterMetaType<DockItemInfo>();
    qRegisterMetaType<DockItemInfos>("DockItemInfos");
    qDBusRegisterMetaType<DockItemInfos>();
}

DBusDockAdaptors::DBusDockAdaptors(WindowManager* parent)
    : QDBusAbstractAdaptor(parent)
    , m_gsettings(Utils::SettingsPtr("com.deepin.dde.dock.mainwindow", QByteArray(), this))
    , m_windowManager(parent)
{
    connect(parent, &WindowManager::panelGeometryChanged, this, [ = ] {
        emit DBusDockAdaptors::geometryChanged(geometry());
    });

    if (m_gsettings) {
        connect(m_gsettings, &QGSettings::changed, this, [ = ] (const QString &key) {
            if (key == "onlyShowPrimary") {
                Q_EMIT showInPrimaryChanged(m_gsettings->get(key).toBool());
            }
        });
    }

    QList<PluginsItemInterface *> allPlugin = localPlugins();
    connect(DockItemManager::instance(), &DockItemManager::itemInserted, this, [ = ] (const int index, DockItem *item) {
        Q_UNUSED(index);
        if (item->itemType() == DockItem::Plugins
                || item->itemType() == DockItem::FixedPlugin) {
            PluginsItem *pluginItem = static_cast<PluginsItem *>(item);
            for (auto *p : allPlugin) {
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
            for (auto *p : allPlugin) {
                if (p->pluginName() == pluginItem->pluginName()) {
                    Q_EMIT pluginVisibleChanged(p->pluginDisplayName(), getPluginVisible(p->pluginDisplayName()));
                }
            }
        }
    });

    registerPluginInfoMetaType();
}

DBusDockAdaptors::~DBusDockAdaptors()
{

}

void DBusDockAdaptors::callShow()
{
    m_windowManager->callShow();
}

void DBusDockAdaptors::ReloadPlugins()
{
    if (qApp->property("PLUGINSLOADED").toBool())
        return;

    // 发送事件，通知代理来加载插件
    PluginLoadEvent event;
    QCoreApplication::sendEvent(qApp, &event);

    qApp->setProperty("PLUGINSLOADED", true);
    // 退出安全模式
    qApp->setProperty("safeMode", false);
}

QStringList DBusDockAdaptors::GetLoadedPlugins()
{
    QList<PluginsItemInterface *> allPlugin = localPlugins();
    QStringList nameList;
    QMap<QString, QString> map;
    for (auto plugin : allPlugin) {
        // 托盘本身也是一个插件，这里去除掉这个特殊的插件,还有一些没有实际名字的插件
        if (plugin->pluginName() == "tray"
                || plugin->pluginDisplayName().isEmpty()
                || !isPluginValid(plugin->pluginName()))
            continue;

        nameList << plugin->pluginName();
        map.insert(plugin->pluginName(), plugin->pluginDisplayName());
    }

    // 排序,保持和原先任务栏右键菜单中的插件列表顺序一致
    std::sort(nameList.begin(), nameList.end(), [ = ] (const QString &name1, const QString &name2) {
        return name1 > name2;
    });

    QStringList newList;
    for (auto name : nameList) {
        newList.push_back(map[name]);
    }

    return newList;
}

DockItemInfos DBusDockAdaptors::plugins()
{
#define DOCK_QUICK_PLUGINS "Dock_Quick_Plugins"
    // 获取本地加载的插件
    QList<PluginsItemInterface *> allPlugin = localPlugins();
    DockItemInfos pluginInfos;
    QStringList quickSettingKeys = DockSettings::instance()->getQuickPlugins();
    for (PluginsItemInterface *plugin : allPlugin) {
        DockItemInfo info;
        info.name = plugin->pluginName();
        info.displayName = plugin->pluginDisplayName();
        info.itemKey = plugin->pluginName();
        info.settingKey = DOCK_QUICK_PLUGINS;
        info.visible = quickSettingKeys.contains(info.itemKey);
        QSize pixmapSize;
        QIcon lightIcon = getSettingIcon(plugin, pixmapSize, DGuiApplicationHelper::ColorType::LightType);
        if (!lightIcon.isNull()) {
            QBuffer buffer(&info.iconLight);
            if (buffer.open(QIODevice::WriteOnly)) {
                QPixmap pixmap = lightIcon.pixmap(pixmapSize);
                pixmap.save(&buffer, "png");
            }
        }
        QIcon darkIcon = getSettingIcon(plugin, pixmapSize, DGuiApplicationHelper::ColorType::DarkType);
        if (!darkIcon.isNull()) {
            QBuffer buffer(&info.iconDark);
            if (buffer.open(QIODevice::WriteOnly)) {
                QPixmap pixmap = darkIcon.pixmap(pixmapSize);
                pixmap.save(&buffer, "png");
            }
        }
        pluginInfos << info;
    }

    return pluginInfos;
}

void DBusDockAdaptors::resizeDock(int offset, bool dragging)
{
    m_windowManager->resizeDock(offset, dragging);
}

// 返回每个插件的识别Key(所以此值应始终不变)，供个性化插件根据key去匹配每个插件对应的图标
QString DBusDockAdaptors::getPluginKey(const QString &pluginName)
{
    QList<PluginsItemInterface *> allPlugin = localPlugins();
    for (auto plugin : allPlugin) {
        if (plugin->pluginDisplayName() == pluginName)
            return plugin->pluginName();
    }

    return QString();
}

bool DBusDockAdaptors::getPluginVisible(const QString &pluginName)
{
    QList<PluginsItemInterface *> allPlugin = localPlugins();
    for (auto *p : allPlugin) {
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
    QList<PluginsItemInterface *> allPlugin = localPlugins();
    for (auto *p : allPlugin) {
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

void DBusDockAdaptors::setItemOnDock(const QString settingKey, const QString &itemKey, bool visible)
{
    DockSettings *settings = DockSettings::instance();
    if ( keyQuickTrayName == settingKey) {
        visible? settings->setTrayItemOnDock(itemKey) : settings->removeTrayItemOnDock(itemKey);
    } else if (keyQuickPlugins == settingKey) {
        visible? settings->setQuickPlugin(itemKey) : settings->removeQuickPlugin(itemKey);
    }
}

QRect DBusDockAdaptors::geometry() const
{
    return m_windowManager->geometry();
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

    // 未开启窗口特效时，不显示多任务视图插件
    if (name == "multitasking" && !DWindowManagerHelper::instance()->hasComposite())
        return false;

    // 录屏插件不显示,插件名如果有变化，建议发需求，避免任务栏反复适配
    if (name == "deepin-screen-recorder-plugin")
        return false;

    return true;
}

QList<PluginsItemInterface *> DBusDockAdaptors::localPlugins() const
{
    return QuickSettingController::instance()->pluginInSettings();
}

QIcon DBusDockAdaptors::getSettingIcon(PluginsItemInterface *plugin, QSize &pixmapSize, DGuiApplicationHelper::ColorType colorType) const
{
    auto iconSize = [](const QIcon &icon) {
        QList<QSize> iconSizes = icon.availableSizes();
        if (iconSizes.size() > 0)
            return iconSizes[0];

        return defaultIconSize;
    };
    // 先获取控制中心的设置图标
    QIcon icon = plugin->icon(DockPart::DCCSetting, colorType);
    if (!icon.isNull()) {
        pixmapSize = iconSize(icon);
        return icon;
    }

    // 如果插件中没有设置图标，则根据插件的类型，获取其他的图标
    QuickSettingController::PluginAttribute pluginAttr = QuickSettingController::instance()->pluginAttribute(plugin);
    switch(pluginAttr) {
    case QuickSettingController::PluginAttribute::System: {
        icon = plugin->icon(DockPart::SystemPanel, colorType);
        pixmapSize = defaultIconSize;
        QList<QSize> iconSizes = icon.availableSizes();
        if (iconSizes.size() > 0)
            pixmapSize = iconSizes[0];
        break;
    }
    case QuickSettingController::PluginAttribute::Quick: {
        icon = plugin->icon(DockPart::QuickShow, colorType);
        if (icon.isNull())
            icon = plugin->icon(DockPart::QuickPanel, colorType);
        pixmapSize = defaultIconSize;
        QList<QSize> iconSizes = icon.availableSizes();
        if (iconSizes.size() > 0)
            pixmapSize = iconSizes[0];
        break;
    }
    default:
        break;
    }

    if (icon.isNull()) {
        icon = QIcon(":/icons/resources/dcc_dock_plug_in.svg");
        pixmapSize = QSize(20, 20);
    }

    return icon;
}
