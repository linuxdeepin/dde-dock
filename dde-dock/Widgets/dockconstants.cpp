#include "dockconstants.h"

DockConstants::DockConstants(QObject *parent) :
    QObject(parent)
{
}

DockConstants * DockConstants::dockConstants = NULL;
DockConstants * DockConstants::getInstants()
{
    if (dockConstants == NULL)
        dockConstants = new DockConstants();

    return dockConstants;
}

int DockConstants::getIconSize()
{
    return this->iconSize;
}

void DockConstants::setIconSize(int value)
{
    this->iconSize = value;
}
