#ifndef DOCKITEMCONTROLLER_H
#define DOCKITEMCONTROLLER_H

#include "dockpluginscontroller.h"
#include "pluginsiteminterface.h"
#include "dbus/dbusdock.h"
#include "item/dockitem.h"

#include <QObject>

class DockItemController : public QObject
{
    Q_OBJECT

public:
    static DockItemController *instance(QObject *parent);
    ~DockItemController();

    const QList<DockItem *> itemList() const;

signals:
    void itemInserted(const int index, DockItem *item);
    void itemRemoved(DockItem *item);

private:
    explicit DockItemController(QObject *parent = 0);
    void appItemAdded(const QDBusObjectPath &path);
    void appItemRemoved(const QString &appId);
    void pluginsItemAdded(PluginsItemInterface *interface);

private:
    QList<DockItem *> m_itemList;

    DBusDock *m_appInter;
    DockPluginsController *m_pluginsInter;

    static DockItemController *INSTANCE;
};

#endif // DOCKITEMCONTROLLER_H
