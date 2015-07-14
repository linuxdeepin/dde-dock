#ifndef DOCKPLUGINPROXY_H
#define DOCKPLUGINPROXY_H

#include "dockplugininterface.h"
#include "dockpluginproxyinterface.h"
#include "abstractdockitem.h"

class DockPluginProxy : public QObject, public DockPluginProxyInterface
{
    Q_OBJECT
public:
    DockPluginProxy(DockPluginInterface * plugin, QObject * parent = 0);

    DockPluginInterface * plugin();

    Dock::DockMode dockMode() Q_DECL_OVERRIDE;

    void itemAddedEvent(QString uuid) Q_DECL_OVERRIDE;
    void itemRemovedEvent(QString uuid) Q_DECL_OVERRIDE;
    void itemSizeChangedEvent(QString uuid) Q_DECL_OVERRIDE;

signals:
    void itemAdded(AbstractDockItem * item);
    void itemRemoved(AbstractDockItem * item);

private:
    QList<AbstractDockItem*> m_items;
    DockPluginInterface * m_plugin;

    AbstractDockItem * getItem(QString uuid);
};

#endif // DOCKPLUGINPROXY_H
