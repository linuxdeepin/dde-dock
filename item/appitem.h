#ifndef APPITEM_H
#define APPITEM_H

#include "dockitem.h"
#include "dbus/dbusdockentry.h"
#include "dbus/dbusclientmanager.h"

class AppItem : public DockItem
{
    Q_OBJECT

public:
    explicit AppItem(const QDBusObjectPath &entry, QWidget *parent = nullptr);

private:
    void paintEvent(QPaintEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);

    void initClientManager();
    void entryDataChanged(const QString &key, const QString &value);

private:
    DBusDockEntry *m_itemEntry;

    QMap<QString, QString> m_data;
    QMap<uint, QString> m_windows;

    static DBusClientManager *ClientInter;
    static uint ActiveWindowId;
};

#endif // APPITEM_H
