#ifndef DOCKITEM_H
#define DOCKITEM_H

#include <QWidget>
#include <QFrame>
#include <QMouseEvent>
#include <QRectF>
#include <QDebug>
#include "dockconstants.h"
#include "appicon.h"

class DockItem : public QFrame
{
    Q_OBJECT
public:
    explicit DockItem(QWidget *parent = 0);
    virtual ~DockItem(){}

    virtual QWidget * getContents();

    virtual void setTitle(const QString &title);
    virtual void setIcon(const QString &iconPath, int size = 42);
    virtual void setMoveable(bool value);
    virtual bool moveable();
    virtual void setActived(bool value);
    virtual bool actived();
    virtual void setIndex(int value);
    virtual int index();

protected:
    AppIcon * appIcon = NULL;

    bool itemMoveable = true;
    bool itemActived = false;

    QString itemTitle = "";
    QString itemIconPath = "";
    int itemIndex = 0;

};

#endif // DOCKITEM_H
