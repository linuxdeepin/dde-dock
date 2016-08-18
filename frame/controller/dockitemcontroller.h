#ifndef DOCKITEMCONTROLLER_H
#define DOCKITEMCONTROLLER_H

#include "dockpluginscontroller.h"
#include "pluginsiteminterface.h"
#include "dbus/dbusdock.h"
#include "item/dockitem.h"
#include "item/stretchitem.h"
#include "item/appitem.h"
#include "item/placeholderitem.h"
#include "item/containeritem.h"

#include <QObject>

class DockItemController : public QObject
{
    Q_OBJECT

public:
    static DockItemController *instance(QObject *parent);
    ~DockItemController();

    const QList<DockItem *> itemList() const;
    bool appIsOnDock(const QString &appDesktop) const;
    bool itemIsInContainer(DockItem * const item) const;

signals:
    void itemInserted(const int index, DockItem *item) const;
    void itemRemoved(DockItem *item) const;
    void itemMoved(DockItem *item, const int index) const;
    void itemManaged(DockItem *item) const;

public slots:
    void refershItemsIcon();
    void updatePluginsItemOrderKey();
    void itemMove(DockItem * const moveItem, DockItem * const replaceItem);
    void itemDroppedIntoContainer(DockItem * const item);
    void itemDragOutFromContainer(DockItem * const item);
    void placeholderItemAdded(PlaceholderItem *item, DockItem *position);
    void placeholderItemDocked(const QString &appDesktop, DockItem *position);
    void placeholderItemRemoved(PlaceholderItem *item);

private:
    explicit DockItemController(QObject *parent = 0);
    void appItemAdded(const QDBusObjectPath &path, const int index);
    void appItemRemoved(const QString &appId);
    void appItemRemoved(AppItem *appItem);
    void pluginItemInserted(PluginsItem *item);
    void pluginItemRemoved(PluginsItem *item);
    void reloadAppItems();

private:
    QList<DockItem *> m_itemList;

    QTimer *m_updatePluginsOrderTimer;

    DBusDock *m_appInter;
    DockPluginsController *m_pluginsInter;
    StretchItem *m_placeholderItem;
    ContainerItem *m_containerItem;

    static DockItemController *INSTANCE;
};

#endif // DOCKITEMCONTROLLER_H
