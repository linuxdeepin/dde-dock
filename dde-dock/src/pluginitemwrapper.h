#ifndef PLUGINITEMWRAPPER_H
#define PLUGINITEMWRAPPER_H

#include "abstractdockitem.h"
#include "dockplugininterface.h"

class PluginItemWrapper : public AbstractDockItem
{
public:
    PluginItemWrapper(DockPluginInterface *plugin, QString uuid, QWidget * parent = 0);


private:
    DockPluginInterface * m_plugin;
    QString m_uuid;
};

#endif // PLUGINITEMWRAPPER_H
