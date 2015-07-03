#ifndef DOCKPLUGINPROXYINTERFACE_H
#define DOCKPLUGINPROXYINTERFACE_H

#include <QString>

class DockPluginProxyInterface
{
public:
    virtual void itemAddedEvent(QString uuid) = 0;
    virtual void itemRemovedEvent(QString uuid) = 0;
};

#endif // DOCKPLUGINPROXYINTERFACE_H
