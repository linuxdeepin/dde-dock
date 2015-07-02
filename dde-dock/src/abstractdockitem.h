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

    virtual void setTitle(const QString &title) { m_itemTitle = title; }
    virtual void setIcon(const QString &iconPath, int size = 42) {
        m_appIcon = new AppIcon(iconPath, this);
        m_appIcon->resize(size, size);
        m_appIcon->move((width() - m_appIcon->width()) / 2,
                        (height() - m_appIcon->height()) / 2);
    }

    virtual void setMoveable(bool value) { m_itemMoveable = value; }
    virtual bool moveable() { return m_itemMoveable; }
    virtual void setActived(bool value) { m_isActived = value; }
    virtual bool actived() { return m_isActived; }
    virtual void setHovered(bool value) { m_isHovered = value; }
    virtual bool hovered() { return m_isHovered; }
    virtual void setIndex(int value) { m_itemIndex = value; }
    virtual int index() { return m_itemIndex; }

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
    QLabel * m_appIcon = NULL;

    bool m_itemMoveable = true;
    bool m_isActived = false;
    bool m_isHovered = false;

    QString m_itemTitle = "";
    QString m_itemIconPath = "";
    QPoint m_itemNextPos;

    int m_itemIndex = 0;

};

#endif // ABSTRACTDOCKITEM_H
