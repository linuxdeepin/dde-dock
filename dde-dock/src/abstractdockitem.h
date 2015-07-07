#ifndef ABSTRACTDOCKITEM_H
#define ABSTRACTDOCKITEM_H

#include <QWidget>
#include <QFrame>
#include <QLabel>
#include "Widgets/appicon.h"

class AbstractDockItem : public QFrame
{
    Q_OBJECT
public:
    explicit AbstractDockItem(QWidget *parent = 0) :
        QFrame(parent) {}
    virtual ~AbstractDockItem() {}

    virtual QWidget * getContents() { return NULL; }

    virtual bool moveable() { return m_moveable; }
    virtual bool actived() { return m_isActived; }

    void resize(int width,int height){
        QFrame::resize(width,height);
        emit widthChanged();
    }
    void resize(const QSize &size){
        QFrame::resize(size);
        emit widthChanged();
    }

    QPoint getNextPos() { return m_itemNextPos; }
    void setNextPos(const QPoint &value) { m_itemNextPos = value; }
    void setNextPos(int x, int y) { m_itemNextPos.setX(x); m_itemNextPos.setY(y); }

    int globalX(){return mapToGlobal(QPoint(0,0)).x();}
    int globalY(){return mapToGlobal(QPoint(0,0)).y();}
    QPoint globalPos(){return mapToGlobal(QPoint(0,0));}
signals:
    void dragStart();
    void dragEntered(QDragEnterEvent * event);
    void dragExited(QDragLeaveEvent * event);
    void drop(QDropEvent * event);
    void mouseEntered();
    void mouseExited();
    void mousePress(int x, int y);
    void mouseRelease(int x, int y);
    void mouseDoubleClick();
    void widthChanged();

protected:

    bool m_moveable = true;
    bool m_isActived = false;

    QPoint m_itemNextPos;
};

#endif // ABSTRACTDOCKITEM_H
