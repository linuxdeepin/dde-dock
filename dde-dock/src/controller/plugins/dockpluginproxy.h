#ifndef DOCKPLUGINPROXY_H
#define DOCKPLUGINPROXY_H

#include "widgets/abstractdockitem.h"
#include "interfaces/dockplugininterface.h"
#include "interfaces/dockpluginproxyinterface.h"

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
    void infoChangedEvent(DockPluginInterface::InfoType type, const QString &id) Q_DECL_OVERRIDE;

signals:
    void itemAdded(AbstractDockItem * item, QString uuid);
    void itemRemoved(AbstractDockItem * item, QString uuid);
    void titleChanged(const QString &id);
    void configurableChanged(QString id);
    void enabledChanged(QString id);

private:
    QMap<QString, AbstractDockItem*> m_items;

    QPluginLoader * m_loader;
    DockPluginInterface * m_plugin;

    void itemSizeChangedEvent(QString id);
    void appletSizeChangedEvent(QString id);
};

#endif // DOCKPLUGINPROXY_H
