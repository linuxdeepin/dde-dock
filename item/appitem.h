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

    const QString appId() const;

private:
    void paintEvent(QPaintEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void mouseMoveEvent(QMouseEvent *e);

    void startDrag();
    void initClientManager();
    void titlesChanged();

private:
    DBusDockEntry *m_itemEntry;

    bool m_draging;

    WindowDict m_titles;
    QString m_id;

    QPoint m_mousePressPos;

    static DBusClientManager *ClientInter;
//    static uint ActiveWindowId;
};

#endif // APPITEM_H
