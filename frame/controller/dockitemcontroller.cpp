#include "dockitemcontroller.h"
#include "dbus/dbusdockentry.h"
#include "item/appitem.h"
#include "item/placeholderitem.h"
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

DockItemController::~DockItemController()
{
    qDeleteAll(m_itemList);
}

const QList<DockItem *> DockItemController::itemList() const
{
    return m_itemList;
}

void DockItemController::itemMove(DockItem * const moveItem, DockItem * const replaceItem)
{
    const DockItem::ItemType moveType = moveItem->itemType();
    const DockItem::ItemType replaceType = replaceItem->itemType();

    // app move
    if (moveType == DockItem::App)
        if (replaceType != DockItem::App && replaceType != DockItem::Placeholder)
            return;

    // plugins move
    if (moveType == DockItem::Plugins)
        if (replaceType != DockItem::Plugins)
            return;

    const int moveIndex = m_itemList.indexOf(moveItem);
    const int replaceIndex = replaceItem->itemType() == DockItem::Placeholder ?
                                // disable insert after placeholder item
                                m_itemList.indexOf(replaceItem) - 1 :
                                m_itemList.indexOf(replaceItem);

    m_itemList.removeAt(moveIndex);
    m_itemList.insert(replaceIndex, moveItem);
    emit itemMoved(moveItem, replaceIndex);

    // for app move, index 0 is launcher item, need to pass it.
    if (moveItem->itemType() == DockItem::App)
        m_appInter->MoveEntry(moveIndex - 1, replaceIndex - 1);
}

DockItemController::DockItemController(QObject *parent)
    : QObject(parent),
      m_appInter(new DBusDock(this)),
      m_pluginsInter(new DockPluginsController(this)),
      m_placeholderItem(new PlaceholderItem)
{
//    m_placeholderItem->hide();

    m_itemList.append(new LauncherItem);
    for (auto entry : m_appInter->entries())
        m_itemList.append(new AppItem(entry));
    m_itemList.append(m_placeholderItem);

    connect(m_appInter, &DBusDock::EntryAdded, this, &DockItemController::appItemAdded);
    connect(m_appInter, &DBusDock::EntryRemoved, this, static_cast<void (DockItemController::*)(const QString &)>(&DockItemController::appItemRemoved));
    connect(m_appInter, &DBusDock::ServiceRestarted, this, &DockItemController::reloadAppItems);

    connect(m_pluginsInter, &DockPluginsController::pluginItemInserted, this, &DockItemController::pluginItemInserted, Qt::QueuedConnection);
    connect(m_pluginsInter, &DockPluginsController::pluginItemRemoved, this, &DockItemController::pluginItemRemoved, Qt::QueuedConnection);
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
    appItem->deleteLater();
    m_itemList.removeOne(appItem);
}

void DockItemController::pluginItemInserted(PluginsItem *item)
{
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
    if (itemSortKey == -1)
    {
        insertIndex = m_itemList.size();
    }
    else if (itemSortKey == 0)
    {
        insertIndex = firstPluginPosition;
    }
    else
    {
        // TODO: compare other plugins to find insert position
        Q_ASSERT(false);
    }

//    qDebug() << insertIndex << item;

    m_itemList.insert(insertIndex, item);
    emit itemInserted(insertIndex, item);
}

void DockItemController::pluginItemRemoved(PluginsItem *item)
{
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
