#ifndef DOCKPLUGINPROXYINTERFACE_H
#define DOCKPLUGINPROXYINTERFACE_H

#include <QString>

#include "dockconstants.h"

class DockPluginProxyInterface
{
public:
    virtual Dock::DockMode dockMode() = 0;

    virtual void itemAddedEvent(QString id) = 0;
    virtual void itemRemovedEvent(QString id) = 0;

    virtual void itemSizeChangedEvent(QString id) = 0;
    virtual void appletSizeChangedEvent(QString id) = 0;
};

#endif // DOCKPLUGINPROXYINTERFACE_H
