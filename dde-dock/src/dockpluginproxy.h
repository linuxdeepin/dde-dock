#ifndef DOCKPLUGINPROXY_H
#define DOCKPLUGINPROXY_H

#include "dockplugininterface.h"
#include "dockpluginproxyinterface.h"
#include "abstractdockitem.h"

class QPluginLoader;
class DockPluginProxy : public QObject, public DockPluginProxyInterface
{
    Q_OBJECT
public:
    DockPluginProxy(QPluginLoader * loader, DockPluginInterface * plugin);
    ~DockPluginProxy();

    bool isSystemPlugin();
    DockPluginInterface * plugin();

    Dock::DockMode dockMode() Q_DECL_OVERRIDE;

    void itemAddedEvent(QString id) Q_DECL_OVERRIDE;
    void itemRemovedEvent(QString id) Q_DECL_OVERRIDE;
    void itemSizeChangedEvent(QString id) Q_DECL_OVERRIDE;
    void appletSizeChangedEvent(QString id) Q_DECL_OVERRIDE;

signals:
    void itemAdded(AbstractDockItem * item, QString uuid);
    void itemRemoved(AbstractDockItem * item, QString uuid);

private:
    QMap<QString, AbstractDockItem*> m_items;

    QPluginLoader * m_loader;
    DockPluginInterface * m_plugin;
};

#endif // DOCKPLUGINPROXY_H
