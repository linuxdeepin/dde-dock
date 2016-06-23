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
    connect(m_appInter, &DBusDock::EntryRemoved, this, &DockItemController::appItemRemoved);

    connect(m_pluginsInter, &DockPluginsController::pluginsInserted, this, &DockItemController::pluginsItemInserted);
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

        emit itemRemoved(m_itemList[i]);
        m_itemList[i]->deleteLater();
        m_itemList.removeAt(i);

        break;
    }
}

void DockItemController::pluginsItemInserted(PluginsItem *item)
{
    m_itemList.append(item);
    emit itemInserted(m_itemList.size(), item);
}
