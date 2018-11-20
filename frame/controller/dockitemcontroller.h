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

#ifndef DOCKITEMCONTROLLER_H
#define DOCKITEMCONTROLLER_H

#include "dockpluginscontroller.h"
#include "pluginsiteminterface.h"
#include "dbus/dbusdock.h"
#include "item/dockitem.h"
#include "item/stretchitem.h"
#include "item/appitem.h"
#include "item/placeholderitem.h"
#include "item/containeritem.h"

#include <QObject>

class DockItemController : public QObject
{
    Q_OBJECT

public:
    static DockItemController *instance(QObject *parent = nullptr);

    const QList<QPointer<DockItem> > itemList() const;
    const QList<PluginsItemInterface *> pluginList() const;
    bool appIsOnDock(const QString &appDesktop) const;
    bool itemIsInContainer(DockItem * const item) const;
    void setDropping(const bool dropping);
    void startLoadPlugins() const;

signals:
    void itemInserted(const int index, DockItem *item) const;
    void itemRemoved(DockItem *item) const;
    void itemMoved(DockItem *item, const int index) const;
    void itemManaged(DockItem *item) const;
    void itemUpdated(DockItem *item) const;
    void fashionTraySizeChanged(const QSize &traySize) const;

public slots:
    void refershItemsIcon();
    void sortPluginItems();
    void updatePluginsItemOrderKey();
    void itemMove(DockItem * const moveItem, DockItem * const replaceItem);
    void itemDroppedIntoContainer(DockItem * const item);
    void itemDragOutFromContainer(DockItem * const item);
    void placeholderItemAdded(PlaceholderItem *item, DockItem *position);
    void placeholderItemDocked(const QString &appDesktop, DockItem *position);
    void placeholderItemRemoved(PlaceholderItem *item);
    void refreshFSTItemSpliterVisible();

private:
    explicit DockItemController(QObject *parent = nullptr);
    void appItemAdded(const QDBusObjectPath &path, const int index);
    void appItemRemoved(const QString &appId);
    void appItemRemoved(AppItem *appItem);
    void pluginItemInserted(PluginsItem *item);
    void pluginItemRemoved(PluginsItem *item);
    void reloadAppItems();

private:
    QList<QPointer<DockItem>> m_itemList;

    QTimer *m_updatePluginsOrderTimer;

    DBusDock *m_appInter;
    DockPluginsController *m_pluginsInter;
    StretchItem *m_placeholderItem;
    ContainerItem *m_containerItem;

    static DockItemController *INSTANCE;
};

#endif // DOCKITEMCONTROLLER_H
