#ifndef APPITEM_H
#define APPITEM_H

#include <QObject>
#include <QWidget>
#include <QPushButton>
#include <QMouseEvent>
#include <QDrag>
#include <QRectF>
#include <QDrag>
#include <QMimeData>
#include <QPixmap>
#include <QDebug>
#include "dockitem.h"
#include "dockconstants.h"
#include "appicon.h"
#include "appbackground.h"

class AppItem : public DockItem
{
    Q_OBJECT
    Q_PROPERTY(QPoint pos READ pos WRITE move)
public:
    AppItem(QWidget *parent = 0);
    AppItem(QString title, QWidget *parent = 0);
    AppItem(QString title, QString iconPath, QWidget *parent = 0);
    ~AppItem();

    void resize(const QSize &size);
    void resize(int width, int height);
    QPoint getNextPos();
    void setNextPos(const QPoint &value);
    void setNextPos(int x, int y);

protected:
    void mousePressEvent(QMouseEvent *);
    void mouseReleaseEvent(QMouseEvent *);
    void mouseDoubleClickEvent(QMouseEvent *);
    void mouseMoveEvent(QMouseEvent *);
    void enterEvent(QEvent * event);
    void leaveEvent(QEvent * event);
    void dragEnterEvent(QDragEnterEvent * event);
    void dragLeaveEvent(QDragLeaveEvent * event);
    void dropEvent(QDropEvent * event);

signals:
    void dragStart(AppItem *item);
    void dragEntered(QDragEnterEvent * event,AppItem *item);
    void dragExited(QDragLeaveEvent * event,AppItem *item);
    void drop(QDropEvent * event,AppItem *item);
    void mouseEntered(AppItem *item);
    void mouseExited(AppItem *item);
    void mousePress(int x, int y, AppItem *item);
    void mouseRelease(int x, int y, AppItem *item);
    void mouseDoubleClick( AppItem *item);

private:
    void resizeResources();
    void initBackground();

private:
    AppBackground * appBackground = NULL;
    QPoint nextPos;
};

#endif // APPITEM_H
