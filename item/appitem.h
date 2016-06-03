#ifndef APPITEM_H
#define APPITEM_H

#include "dockitem.h"
#include "dbus/dbusdockentry.h"

class AppItem : public DockItem
{
    Q_OBJECT

public:
    explicit AppItem(const QDBusObjectPath &entry, QWidget *parent = nullptr);

private:
    void paintEvent(QPaintEvent *e);

private:
    DBusDockEntry *m_itemEntry;
};

#endif // APPITEM_H
