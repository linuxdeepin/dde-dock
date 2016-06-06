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
    void mouseReleaseEvent(QMouseEvent *e);

    void entryDataChanged(const QString &key, const QString &value);

private:
    DBusDockEntry *m_itemEntry;

    QMap<QString, QString> m_data;
//    QMap<uint, QString> m_windows;
};

#endif // APPITEM_H
