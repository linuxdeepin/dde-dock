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
#include "dockconstants.h"
#include "appicon.h"
#include "appbackground.h"
#include "QDebug"

class QDragEnterEvent;
class QDropEvent;
class AppItem : public QFrame
{
    Q_OBJECT
    Q_PROPERTY(QPoint pos READ pos WRITE move)
public:
    AppItem(QWidget *parent = 0);
    AppItem(QString title, QWidget *parent = 0);
    AppItem(QString title, QString iconPath, QWidget *parent = 0);
    ~AppItem();

    void setTitle(const QString &title);
    void setIcon(const QString &iconPath, int size = 42);
    void resize(const QSize &size);
    void resize(int width, int height);
    void setMoveable(bool value);
    bool getMoveable();
    void setIndex(int value);
    int getIndex();
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
    AppIcon * appIcon = NULL;
    QPoint nextPos;
    int itemIndex;

    bool itemMoveable = true;
    bool itemHover = false;
    bool itemActived = false;
    bool itemDraged = false;

    QString itemTitle = "";
    QString itemIconPath = "";
};

#endif // APPITEM_H
