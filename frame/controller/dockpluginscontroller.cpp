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

#include "dockpluginscontroller.h"
#include "pluginsiteminterface.h"
#include "item/traypluginitem.h"

#include "../plugins/datetime/datetimeplugin.h"
#include "../plugins/keyboard-layout/keyboardplugin.h"
#include "../plugins/overlay-warning/overlay-warning-plugin.h"
#include "../plugins/trash/trashplugin.h"
#include "../plugins/tray/trayplugin.h"
#include "../plugins/shutdown/shutdownplugin.h"

#include <QDebug>
#include <QDir>

DockPluginsController::DockPluginsController(QObject *parent)
    : AbstractPluginsController(parent)
{
    setObjectName("DockPlugin");
}

void DockPluginsController::itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    QMap<PluginsItemInterface *, QMap<QString, QObject *>> &mPluginsMap = pluginsMap();

    // check if same item added
    if (mPluginsMap.contains(itemInter))
        if (mPluginsMap[itemInter].contains(itemKey))
            return;

    PluginsItem *item = nullptr;
    if (itemInter->pluginName() == "tray") {
        item = new TrayPluginItem(itemInter, itemKey);
        if (item->graphicsEffect()) {
            item->graphicsEffect()->setEnabled(false);
        }
        connect(static_cast<TrayPluginItem *>(item), &TrayPluginItem::fashionTraySizeChanged,
                this, &DockPluginsController::fashionTraySizeChanged, Qt::UniqueConnection);
    } else {
        item = new PluginsItem(itemInter, itemKey);
    }

    item->setVisible(false);

    mPluginsMap[itemInter][itemKey] = item;

    emit pluginItemInserted(item);
}

void DockPluginsController::itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = static_cast<PluginsItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    item->update();

    emit pluginItemUpdated(item);
}

void DockPluginsController::itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = static_cast<PluginsItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    item->detachPluginWidget();

    emit pluginItemRemoved(item);

    QMap<PluginsItemInterface *, QMap<QString, QObject *>> &mPluginsMap = pluginsMap();
    mPluginsMap[itemInter].remove(itemKey);

    // do not delete the itemWidget object(specified in the plugin interface)
    item->centralWidget()->setParent(nullptr);

    // just delete our wrapper object(PluginsItem)
    item->deleteLater();
}

void DockPluginsController::requestWindowAutoHide(PluginsItemInterface * const itemInter, const QString &itemKey, const bool autoHide)
{
    PluginsItem *item = static_cast<PluginsItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    Q_EMIT item->requestWindowAutoHide(autoHide);
}

void DockPluginsController::requestRefreshWindowVisible(PluginsItemInterface * const itemInter, const QString &itemKey)
{
    PluginsItem *item = static_cast<PluginsItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    Q_EMIT item->requestRefreshWindowVisible();
}

void DockPluginsController::requestSetAppletVisible(PluginsItemInterface * const itemInter, const QString &itemKey, const bool visible)
{
    PluginsItem *item = static_cast<PluginsItem *>(pluginItemAt(itemInter, itemKey));
    if (!item)
        return;

    if (visible) {
        item->showPopupApplet(itemInter->itemPopupApplet(itemKey));
    } else {
        item->hidePopup();
    }
}

void DockPluginsController::startLoader()
{
    const QList<const QObject *> list{
        new DatetimePlugin,
        new KeyboardPlugin,
        new OverlayWarningPlugin,
        new TrashPlugin,
        new TrayPlugin,
        new ShutdownPlugin,
    };

    for (const QObject *obj : list) {
        loadPlugin(qobject_cast<PluginsItemInterface *>(obj));
    }

    QString pluginsDir("../plugins");
    if (!QDir(pluginsDir).exists()) {
        pluginsDir = "/usr/lib/dde-dock/plugins";
    }
    qDebug() << "using dock plugins dir:" << pluginsDir;

    AbstractPluginsController::startLoader(new PluginLoader(pluginsDir, this));
}
