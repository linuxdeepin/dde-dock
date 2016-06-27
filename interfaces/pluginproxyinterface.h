#ifndef PLUGINPROXYINTERFACE_H
#define PLUGINPROXYINTERFACE_H

#include "constants.h"

#include <QtCore>

class PluginsItemInterface;
class PluginProxyInterface
{
public:
    virtual void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey) = 0;

    virtual Dock::DisplayMode displayMode() const = 0;
};

#endif // PLUGINPROXYINTERFACE_H
