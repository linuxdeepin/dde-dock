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
    static void setIconBaseSize(const int size);
    static int iconBaseSize();
    static int itemBaseHeight();
    static int itemBaseWidth();

private:
    void paintEvent(QPaintEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void mouseMoveEvent(QMouseEvent *e);
    void resizeEvent(QResizeEvent *e);

    void invokedMenuItem(const QString &itemId, const bool checked);
    const QString contextMenu() const;

    void startDrag();
    void initClientManager();

private slots:
    void updateTitle();
    void updateIcon();

private:
    DBusDockEntry *m_itemEntry;

    bool m_draging;

    WindowDict m_titles;
    QString m_id;
    QPixmap m_icon;

    static int IconBaseSize;
    static QPoint MousePressPos;
    static DBusClientManager *ClientInter;
//    static uint ActiveWindowId;
};

#endif // APPITEM_H
