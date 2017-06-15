#ifndef APPITEM_H
#define APPITEM_H

#include "dockitem.h"
#include "components/_previewcontainer.h"
#include "dbus/dbusdockentry.h"
#include "dbus/dbusclientmanager.h"

#include <QGraphicsView>
#include <QGraphicsItem>
#include <QFuture>

class AppItem : public DockItem
{
    Q_OBJECT

public:
    explicit AppItem(const QDBusObjectPath &entry, QWidget *parent = nullptr);
    ~AppItem();

    const QString appId() const;
    void updateWindowIconGeometries();
    static void setIconBaseSize(const int size);
    static int iconBaseSize();
    static int itemBaseHeight();
    static int itemBaseWidth();

    inline ItemType itemType() const {return App;}

signals:
    void requestActivateWindow(const WId wid) const;
    void requestPreviewWindow(const WId wid) const;
    void requestCancelPreview() const;

private:
    void moveEvent(QMoveEvent *e);
    void paintEvent(QPaintEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);
    void mousePressEvent(QMouseEvent *e);
    void mouseMoveEvent(QMouseEvent *e);
    void wheelEvent(QWheelEvent *e);
    void resizeEvent(QResizeEvent *e);
    void dragEnterEvent(QDragEnterEvent *e);
    void dragMoveEvent(QDragMoveEvent *e);
    void dropEvent(QDropEvent *e);

    void showHoverTips();
    void invokedMenuItem(const QString &itemId, const bool checked);
    const QString contextMenu() const;
    QWidget *popupTips();

    void startDrag();

private slots:
    void updateTitle();
    void refershIcon();
    void activeChanged();
    void showPreview();

    void gotSmallIcon();
    void gotLargeIcon();

private:
    QLabel *m_appNameTips;
    _PreviewContainer *m_appPreviewTips;
    DBusDockEntry *m_itemEntry;

    QGraphicsView *m_itemView;
    QGraphicsScene *m_itemScene;

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

    QTimer *m_updateIconGeometryTimer;

    static int IconBaseSize;
    static QPoint MousePressPos;
};

#endif // APPITEM_H
