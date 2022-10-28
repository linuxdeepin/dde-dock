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

#include "systemplugincontroller.h"
#include "pluginsiteminterface.h"
#include "utils.h"
#include "systempluginitem.h"

#include <QDebug>
#include <QDir>

SystemPluginController::SystemPluginController(QObject *parent)
    : AbstractPluginsController(parent)
{
    setObjectName("SystemTray");
}

void SystemPluginController::pluginItemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    QMap<PluginsItemInterface *, QMap<QString, QObject *>> &mPluginsMap = pluginsMap();

    // check if same item added
    if (mPluginsMap.contains(itemInter))
        if (mPluginsMap[itemInter].contains(itemKey))
            return;

    SystemPluginItem *item = new SystemPluginItem(itemInter, itemKey);
    connect(item, &SystemPluginItem::itemVisibleChanged, this, [ = ] (bool visible){
        emit visible ? pluginItemAdded(itemKey, item) : pluginItemRemoved(itemKey, item);
    }, Qt::QueuedConnection);

    mPluginsMap[itemInter][itemKey] = item;

    // 隐藏的插件不加入到布局中
    if (Utils::SettingValue(QString("com.deepin.dde.dock.module.") + itemInter->pluginName(), QByteArray(), "enable", true).toBool())
        emit pluginItemAdded(itemKey, item);
}

void SystemPluginController::pluginItemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    SystemPluginItem *item = static_cast<SystemPluginItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    item->update();

    emit pluginItemUpdated(itemKey, item);
}

void SystemPluginController::pluginItemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    SystemPluginItem *item = static_cast<SystemPluginItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    item->detachPluginWidget();

    emit pluginItemRemoved(itemKey, item);

    QMap<PluginsItemInterface *, QMap<QString, QObject *>> &mPluginsMap = pluginsMap();
    mPluginsMap[itemInter].remove(itemKey);

    // do not delete the itemWidget object(specified in the plugin interface)
    item->centralWidget()->setParent(nullptr);

    // just delete our wrapper object(PluginsItem)
    item->deleteLater();
}

void SystemPluginController::requestPluginWindowAutoHide(PluginsItemInterface * const itemInter, const QString &itemKey, const bool autoHide)
{
    SystemPluginItem *item = static_cast<SystemPluginItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    Q_EMIT item->requestWindowAutoHide(autoHide);
}

void SystemPluginController::requestRefreshPluginWindowVisible(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    SystemPluginItem *item = static_cast<SystemPluginItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    Q_EMIT item->requestRefershWindowVisible();
}

void SystemPluginController::requestSetPluginAppletVisible(PluginsItemInterface * const itemInter, const QString &itemKey, const bool visible)
{
    SystemPluginItem *item = static_cast<SystemPluginItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    if (visible) {
        item->showPopupApplet(itemInter->itemPopupApplet(itemKey));
    } else {
        item->hidePopup();
    }
}

int SystemPluginController::systemTrayItemSortKey(const QString &itemKey)
{
    auto inter = pluginInterAt(itemKey);

    if (!inter) {
        return -1;
    }

    return inter->itemSortKey(itemKey);
}

void SystemPluginController::setSystemTrayItemSortKey(const QString &itemKey, const int order)
{
    auto inter = pluginInterAt(itemKey);

    if (!inter) {
        return;
    }

    inter->setSortKey(itemKey, order);
}

const QVariant SystemPluginController::getValueSystemTrayItem(const QString &itemKey, const QString &key, const QVariant &fallback)
{
    auto inter = pluginInterAt(itemKey);

    if (!inter) {
        return QVariant();
    }

    return getPluginValue(inter, key, fallback);
}

void SystemPluginController::saveValueSystemTrayItem(const QString &itemKey, const QString &key, const QVariant &value)
{
    auto inter = pluginInterAt(itemKey);

    if (!inter) {
        return;
    }

    savePluginValue(inter, key, value);
}

void SystemPluginController::startLoader()
{
    QString pluginsDir("../plugins/system-trays");
    if (!QDir(pluginsDir).exists()) {
        pluginsDir = "/usr/lib/dde-dock/plugins/system-trays";
    }
    qDebug() << "using system tray plugins dir:" << pluginsDir;

    AbstractPluginsController::startLoader(new PluginLoader(pluginsDir, this));
}
