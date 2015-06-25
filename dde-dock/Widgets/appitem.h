#ifndef APPITEM_H
#define APPITEM_H

#include <QObject>
#include <QWidget>
#include <QPushButton>
#include <QMouseEvent>
#include <QDrag>
#include <QRectF>
#include "dockconstants.h"
#include "appicon.h"
#include "appbackground.h"
#include "QDebug"

class AppItem : public QFrame
{
    Q_OBJECT
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

protected:
    void mousePressEvent(QMouseEvent *);
    void mouseReleaseEvent(QMouseEvent *);
    void mouseMoveEvent(QMouseEvent *);
    void mouseDoubleClickEvent(QMouseEvent *);

signals:
    void mousePress(int x, int y, AppItem *item);
    void mouseRelease(int x, int y, AppItem *item);
    void mouseMove(int x, int y, AppItem *item);
    void mouseDoubleClick( AppItem *item);

private:
    void resizeResources();
    void initBackground();

private:
    AppBackground * appBackground = NULL;
    AppIcon * appIcon = NULL;
    int itemIndex;

    bool itemMoveable = true;
    bool itemHover = false;
    bool itemActived = false;
    bool itemDraged = false;

    QString itemTitle = "";
    QString itemIconPath = "";
};

#endif // APPITEM_H
