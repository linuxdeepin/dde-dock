#ifndef DOCKITEMCONTROLLER_H
#define DOCKITEMCONTROLLER_H

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
    void dockItemCountChanged(const int count) const;

private:
    explicit DockItemController(QObject *parent = 0);

private:
    QList<DockItem *> m_itemList;

    DBusDock *m_dockInter;

    static DockItemController *INSTANCE;
};

#endif // DOCKITEMCONTROLLER_H
