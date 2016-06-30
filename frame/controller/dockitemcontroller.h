#ifndef DOCKITEMCONTROLLER_H
#define DOCKITEMCONTROLLER_H

#include "dockpluginscontroller.h"
#include "pluginsiteminterface.h"
#include "dbus/dbusdock.h"
#include "item/dockitem.h"
#include "item/placeholderitem.h"
#include "item/appitem.h"

#include <QObject>

class DockItemController : public QObject
{
    Q_OBJECT

public:
    static DockItemController *instance(QObject *parent);
    ~DockItemController();

    const QList<DockItem *> itemList() const;

signals:
    void itemInserted(const int index, DockItem *item) const;
    void itemRemoved(DockItem *item) const;
    void itemMoved(DockItem *item, const int index) const;

public slots:
    void itemMove(DockItem * const moveItem, DockItem * const replaceItem);

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

    DBusDock *m_appInter;
    DockPluginsController *m_pluginsInter;
    PlaceholderItem *m_placeholderItem;

    static DockItemController *INSTANCE;
};

#endif // DOCKITEMCONTROLLER_H
