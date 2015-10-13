#ifndef DOCKPLUGINPROXYINTERFACE_H
#define DOCKPLUGINPROXYINTERFACE_H

#include <QString>

#include "dockconstants.h"
#include "dockplugininterface.h"

class DockPluginProxyInterface
{
public:
    virtual Dock::DockMode dockMode() = 0;

    virtual void itemAddedEvent(QString id) = 0;
    virtual void itemRemovedEvent(QString id) = 0;
    virtual void infoChangedEvent(DockPluginInterface::InfoType type, const QString &id) = 0;
};

#endif // DOCKPLUGINPROXYINTERFACE_H
