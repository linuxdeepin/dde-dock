#ifndef DOCKPLUGINSCONTROLLER_H
#define DOCKPLUGINSCONTROLLER_H

#include "item/pluginsitem.h"
#include "pluginproxyinterface.h"

#include <QPluginLoader>
#include <QList>
#include <QMap>

class DockItemController;
class PluginsItemInterface;
class DockPluginsController : public QObject, PluginProxyInterface
{
    Q_OBJECT

public:
    explicit DockPluginsController(DockItemController *itemControllerInter = 0);
    ~DockPluginsController();

    // implements PluginProxyInterface
    void itemAdded(PluginsItemInterface * const itemInter, const QString &itemKey);
    void itemUpdate(PluginsItemInterface * const itemInter, const QString &itemKey);
    void itemRemoved(PluginsItemInterface * const itemInter, const QString &itemKey);
    void requestPopupApplet(PluginsItemInterface * const itemInter, const QString &itemKey);

signals:
    void pluginItemInserted(PluginsItem *pluginItem) const;
    void pluginItemRemoved(PluginsItem *pluginItem) const;

private slots:
    void loadPlugins();
    void displayModeChanged();
    void positionChanged();

private:
    bool eventFilter(QObject *o, QEvent *e);
    PluginsItem *pluginItemAt(PluginsItemInterface * const itemInter, const QString &itemKey) const;

private:
//    QList<PluginsItemInterface *> m_pluginsInterfaceList;
//    QList<QPluginLoader *> m_pluginLoaderList;
    QMap<PluginsItemInterface *, QMap<QString, PluginsItem *>> m_pluginList;
    DockItemController *m_itemControllerInter;
};

#endif // DOCKPLUGINSCONTROLLER_H
