#include "dockitemcontroller.h"

DockItemController *DockItemController::INSTANCE = nullptr;

DockItemController *DockItemController::instance(QObject *parent)
{
    if (!INSTANCE)
        INSTANCE = new DockItemController(parent);

    return INSTANCE;
}

DockItemController::DockItemController(QObject *parent)
    : QObject(parent)
{

}
