#ifndef DOCKPLUGINSCONTROLLER_H
#define DOCKPLUGINSCONTROLLER_H

#include "item/pluginsitem.h"
#include "pluginproxyinterface.h"

#include <QPluginLoader>
#include <QList>

class DockItemController;
class PluginsItemInterface;
class DockPluginsController : public QObject, PluginProxyInterface
{
    Q_OBJECT

public:
    explicit DockPluginsController(DockItemController *itemControllerInter = 0);
    ~DockPluginsController();

    void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey);

signals:
    void pluginItemInserted(PluginsItem *pluginsItem) const;

private slots:
    void loadPlugins();

private:
//    QList<PluginsItemInterface *> m_pluginsInterfaceList;
//    QList<QPluginLoader *> m_pluginLoaderList;
    DockItemController *m_itemControllerInter;
};

#endif // DOCKPLUGINSCONTROLLER_H
