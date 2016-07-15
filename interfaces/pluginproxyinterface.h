#ifndef PLUGINPROXYINTERFACE_H
#define PLUGINPROXYINTERFACE_H

#include "constants.h"

#include <QtCore>

class PluginsItemInterface;
class PluginProxyInterface
{
public:
    virtual void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey) = 0;
    virtual void itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey) = 0;
    virtual void itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey) = 0;
};

#endif // PLUGINPROXYINTERFACE_H
