#include "dockitemcontroller.h"
#include "dbus/dbusdockentry.h"
#include "item/appitem.h"
#include "item/placeholderitem.h"

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

DockItemController::DockItemController(QObject *parent)
    : QObject(parent),
      m_entryManager(new DBusDockEntryManager(this))
{
    for (auto entry : m_entryManager->entries())
        m_itemList.append(new AppItem(entry));
    m_itemList.append(new PlaceholderItem);
}
