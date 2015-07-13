#ifndef PLUGINITEMWRAPPER_H
#define PLUGINITEMWRAPPER_H

#include "abstractdockitem.h"
#include "dockplugininterface.h"

class PluginItemWrapper : public AbstractDockItem
{
    Q_OBJECT
public:
    PluginItemWrapper(DockPluginInterface *plugin, QString uuid, QWidget * parent = 0);
    virtual ~PluginItemWrapper();

    QWidget * getContents();
    QString uuid() const;

private:
    QWidget *m_pluginItemContents = NULL;
    DockPluginInterface * m_plugin;
    QString m_uuid;
};

#endif // PLUGINITEMWRAPPER_H
