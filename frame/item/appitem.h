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
    void updateWindowIconGeometries();
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
    void dragEnterEvent(QDragEnterEvent *e);
    void dropEvent(QDropEvent *e);

    void invokedMenuItem(const QString &itemId, const bool checked);
    const QString contextMenu() const;
    QWidget *popupTips();

    void startDrag();

private slots:
    void updateTitle();
    void updateIcon();
    void activeChanged();

private:
    QLabel *m_appNameTips;
    DBusDockEntry *m_itemEntry;

    bool m_draging;

    bool m_active;
    WindowDict m_titles;
    QString m_id;
    QPixmap m_smallIcon;
    QPixmap m_largeIcon;
    QPixmap m_horizontalIndicator;
    QPixmap m_verticalIndicator;
    QPixmap m_activeHorizontalIndicator;
    QPixmap m_activeVerticalIndicator;

    QRect m_lastGlobalGeometry;
    QTimer * m_updateIconGeometryTimer;

    static int IconBaseSize;
    static QPoint MousePressPos;
};

#endif // APPITEM_H
