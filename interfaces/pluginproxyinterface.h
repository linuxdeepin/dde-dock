#ifndef PLUGINPROXYINTERFACE_H
#define PLUGINPROXYINTERFACE_H

#include <QtCore>

class PluginsItemInterface;
class PluginProxyInterface
{
public:
    virtual void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey) = 0;

};

#endif // PLUGINPROXYINTERFACE_H
