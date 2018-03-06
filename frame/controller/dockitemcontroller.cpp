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

#include "dockitemcontroller.h"
#include "item/appitem.h"
#include "item/stretchitem.h"
#include "item/launcheritem.h"
#include "item/pluginsitem.h"

#include <QDebug>

DockItemController *DockItemController::INSTANCE = nullptr;

DockItemController *DockItemController::instance(QObject *parent)
{
    if (!INSTANCE)
        INSTANCE = new DockItemController(parent);

    return INSTANCE;
}

const QList<DockItem *> DockItemController::itemList() const
{
    return m_itemList;
}

const QList<PluginsItemInterface *> DockItemController::pluginList() const
{
    return m_pluginsInter->m_pluginList.keys();
}

bool DockItemController::appIsOnDock(const QString &appDesktop) const
{
    return m_appInter->IsOnDock(appDesktop);
}

bool DockItemController::itemIsInContainer(DockItem * const item) const
{
    return m_containerItem->contains(item);
}

void DockItemController::setDropping(const bool dropping)
{
    m_containerItem->setDropping(dropping);
}

void DockItemController::refershItemsIcon()
{
    for (auto item : m_itemList)
    {
        item->refershIcon();
        item->update();
    }
}

void DockItemController::updatePluginsItemOrderKey()
{
    Q_ASSERT(sender() == m_updatePluginsOrderTimer);

    int index = 0;
    for (auto item : m_itemList)
    {
        if (item->itemType() != DockItem::Plugins)
            continue;
        static_cast<PluginsItem *>(item)->setItemSortKey(++index);
    }
}

void DockItemController::itemMove(DockItem * const moveItem, DockItem * const replaceItem)
{
    Q_ASSERT(moveItem != replaceItem);

    const DockItem::ItemType moveType = moveItem->itemType();
    const DockItem::ItemType replaceType = replaceItem->itemType();

    // app move
    if (moveType == DockItem::App)
        if (replaceType != DockItem::App && replaceType != DockItem::Stretch)
            return;

    // plugins move
    if (moveType == DockItem::Plugins)
        if (replaceType != DockItem::Plugins)
            return;

    const int moveIndex = m_itemList.indexOf(moveItem);
    const int replaceIndex = replaceType == DockItem::Stretch ?
                                // disable insert after placeholder item
                                m_itemList.indexOf(replaceItem) - 1 :
                                m_itemList.indexOf(replaceItem);

    m_itemList.removeAt(moveIndex);
    m_itemList.insert(replaceIndex, moveItem);
    emit itemMoved(moveItem, replaceIndex);

    // update plugins sort key if order changed
    if (moveType == DockItem::Plugins || replaceType == DockItem::Plugins)
        m_updatePluginsOrderTimer->start();

    // for app move, index 0 is launcher item, need to pass it.
    if (moveType == DockItem::App && replaceType == DockItem::App)
        m_appInter->MoveEntry(moveIndex - 1, replaceIndex - 1);
}

void DockItemController::itemDroppedIntoContainer(DockItem * const item)
{
    Q_ASSERT(item->itemType() == DockItem::Plugins);

    PluginsItem *pi = static_cast<PluginsItem *>(item);

    if (!pi->allowContainer())
        return;
    if (m_containerItem->contains(item))
        return;

//    qDebug() << "drag into container" << item;

    // remove from main panel
    emit itemRemoved(item);
    m_itemList.removeOne(item);

    // add to container
    pi->setInContainer(true);
    m_containerItem->addItem(item);
}

void DockItemController::itemDragOutFromContainer(DockItem * const item)
{
//    qDebug() << "drag out from container" << item;

    // remove from container
    m_containerItem->removeItem(item);

    // insert to panel
    switch (item->itemType())
    {
    case DockItem::Plugins:
        static_cast<PluginsItem *>(item)->setInContainer(false);
        pluginItemInserted(static_cast<PluginsItem *>(item));
        break;
    default:                  Q_UNREACHABLE();
    }
}

void DockItemController::placeholderItemAdded(PlaceholderItem *item, DockItem *position)
{
    const int pos = m_itemList.indexOf(position);

    m_itemList.insert(pos, item);

    emit itemInserted(pos, item);
}

void DockItemController::placeholderItemDocked(const QString &appDesktop, DockItem *position)
{
    m_appInter->RequestDock(appDesktop, m_itemList.indexOf(position) - 1).waitForFinished();
}

void DockItemController::placeholderItemRemoved(PlaceholderItem *item)
{
    emit itemRemoved(item);

    m_itemList.removeOne(item);
}

DockItemController::DockItemController(QObject *parent)
    : QObject(parent),

      m_updatePluginsOrderTimer(new QTimer(this)),

      m_appInter(new DBusDock(this)),
      m_pluginsInter(new DockPluginsController(this)),
      m_placeholderItem(new StretchItem),
      m_containerItem(new ContainerItem)
{
//    m_placeholderItem->hide();

    m_updatePluginsOrderTimer->setSingleShot(true);
    m_updatePluginsOrderTimer->setInterval(1000);

    m_itemList.append(new LauncherItem);
    for (auto entry : m_appInter->entries())
    {
        AppItem *it = new AppItem(entry);

        connect(it, &AppItem::requestActivateWindow, m_appInter, &DBusDock::ActivateWindow, Qt::QueuedConnection);
        connect(it, &AppItem::requestPreviewWindow, m_appInter, &DBusDock::PreviewWindow);
        connect(it, &AppItem::requestCancelPreview, m_appInter, &DBusDock::CancelPreviewWindow);

        m_itemList.append(it);
    }
    m_itemList.append(m_placeholderItem);
    m_itemList.append(m_containerItem);

    connect(m_updatePluginsOrderTimer, &QTimer::timeout, this, &DockItemController::updatePluginsItemOrderKey);

    connect(m_appInter, &DBusDock::EntryAdded, this, &DockItemController::appItemAdded);
    connect(m_appInter, &DBusDock::EntryRemoved, this, static_cast<void (DockItemController::*)(const QString &)>(&DockItemController::appItemRemoved), Qt::QueuedConnection);
    connect(m_appInter, &DBusDock::ServiceRestarted, this, &DockItemController::reloadAppItems);

    connect(m_pluginsInter, &DockPluginsController::pluginItemInserted, this, &DockItemController::pluginItemInserted, Qt::QueuedConnection);
    connect(m_pluginsInter, &DockPluginsController::pluginItemRemoved, this, &DockItemController::pluginItemRemoved, Qt::QueuedConnection);
    connect(m_pluginsInter, &DockPluginsController::pluginItemUpdated, this, &DockItemController::itemUpdated, Qt::QueuedConnection);

    QMetaObject::invokeMethod(this, "refershItemsIcon", Qt::QueuedConnection);
}

void DockItemController::appItemAdded(const QDBusObjectPath &path, const int index)
{
    // the first index is launcher item
    int insertIndex = 1;

    // -1 for append to app list end
    if (index != -1)
    {
        insertIndex += index;
    } else {
        for (auto item : m_itemList)
            if (item->itemType() == DockItem::App)
                ++insertIndex;
    }

    AppItem *item = new AppItem(path);

    connect(item, &AppItem::requestActivateWindow, m_appInter, &DBusDock::ActivateWindow, Qt::QueuedConnection);
    connect(item, &AppItem::requestPreviewWindow, m_appInter, &DBusDock::PreviewWindow);
    connect(item, &AppItem::requestCancelPreview, m_appInter, &DBusDock::CancelPreviewWindow);

    m_itemList.insert(insertIndex, item);
    emit itemInserted(insertIndex, item);
}

void DockItemController::appItemRemoved(const QString &appId)
{
    for (int i(0); i != m_itemList.size(); ++i)
    {
        if (m_itemList[i]->itemType() != DockItem::App)
            continue;

        AppItem *app = static_cast<AppItem *>(m_itemList[i]);
        if (app->appId() != appId)
            continue;

        appItemRemoved(app);

        break;
    }
}

void DockItemController::appItemRemoved(AppItem *appItem)
{
    emit itemRemoved(appItem);
    m_itemList.removeOne(appItem);
    appItem->deleteLater();
}

void DockItemController::pluginItemInserted(PluginsItem *item)
{
    // check item is in container
    if (item->allowContainer() && item->isInContainer())
    {
        emit itemManaged(item);
        return itemDroppedIntoContainer(item);
    }

    // find first plugins item position
    int firstPluginPosition = -1;
    for (int i(0); i != m_itemList.size(); ++i)
    {
        if (m_itemList[i]->itemType() != DockItem::Plugins)
            continue;

        firstPluginPosition = i;
        break;
    }
    if (firstPluginPosition == -1)
        firstPluginPosition = m_itemList.size();

    // find insert position
    int insertIndex = 0;
    const int itemSortKey = item->itemSortKey();
    if (itemSortKey == -1 || firstPluginPosition == -1)
    {
        insertIndex = m_itemList.size();
    }
    else if (itemSortKey == 0)
    {
        insertIndex = firstPluginPosition;
    }
    else
    {
        insertIndex = m_itemList.size();
        for (int i(firstPluginPosition + 1); i != m_itemList.size() + 1; ++i)
        {
            PluginsItem *pItem = static_cast<PluginsItem *>(m_itemList[i - 1]);
            Q_ASSERT(pItem);

            const int sortKey = pItem->itemSortKey();
            if (sortKey != -1 && itemSortKey > sortKey)
                continue;
            insertIndex = i - 1;
            break;
        }
    }

//    qDebug() << insertIndex << item;

    m_itemList.insert(insertIndex, item);
    emit itemInserted(insertIndex, item);
}

void DockItemController::pluginItemRemoved(PluginsItem *item)
{
    if (m_containerItem->contains(item))
        m_containerItem->removeItem(item);
    else
        emit itemRemoved(item);

    m_itemList.removeOne(item);
}

void DockItemController::reloadAppItems()
{
    // remove old item
    for (auto item : m_itemList)
        if (item->itemType() == DockItem::App)
            appItemRemoved(static_cast<AppItem *>(item));

    // append new item
    for (auto path : m_appInter->entries())
        appItemAdded(path, -1);
}

void DockItemController::sortPluginItems()
{
    int firstPluginIndex = -1;
    for (int i(0); i != m_itemList.size(); ++i)
    {
        if (m_itemList[i]->itemType() == DockItem::Plugins)
        {
            firstPluginIndex = i;
            break;
        }
    }

    if (firstPluginIndex == -1)
        return;

    std::sort(m_itemList.begin() + firstPluginIndex, m_itemList.end(), [](DockItem *a, DockItem *b) -> bool {
        PluginsItem *pa = static_cast<PluginsItem *>(a);
        PluginsItem *pb = static_cast<PluginsItem *>(b);

        const int aKey = pa->itemSortKey();
        const int bKey = pb->itemSortKey();

        if (bKey == -1)
            return true;
        if (aKey == -1)
            return false;

        return aKey < bKey;
    });

    // reset order
    for (int i(firstPluginIndex); i != m_itemList.size(); ++i)
        emit itemMoved(m_itemList[i], i);
}
