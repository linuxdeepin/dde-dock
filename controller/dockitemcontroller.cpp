#include "dockitemcontroller.h"
#include "dbus/dbusdockentry.h"

#include <QDebug>

DockItemController *DockItemController::INSTANCE = nullptr;

DockItemController *DockItemController::instance(QObject *parent)
{
    if (!INSTANCE)
        INSTANCE = new DockItemController(parent);

    return INSTANCE;
}

DockItemController::DockItemController(QObject *parent)
    : QObject(parent),
      m_entryManager(new DBusDockEntryManager(this))
{
}
