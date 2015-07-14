#ifndef DOCKPLUGINPROXYINTERFACE_H
#define DOCKPLUGINPROXYINTERFACE_H

#include <QString>

#include "dockconstants.h"

class DockPluginProxyInterface
{
public:
    virtual Dock::DockMode dockMode() = 0;

    virtual void itemAddedEvent(QString uuid) = 0;
    virtual void itemRemovedEvent(QString uuid) = 0;

    virtual void itemSizeChangedEvent(QString uuid) = 0;
};

#endif // DOCKPLUGINPROXYINTERFACE_H
