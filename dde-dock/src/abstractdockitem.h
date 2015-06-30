#ifndef ABSTRACTDOCKITEM_H
#define ABSTRACTDOCKITEM_H

#include <QWidget>
#include <QFrame>
#include <QMouseEvent>
#include <QRectF>

class AbstractDockItem : public QFrame
{
    Q_OBJECT
public:
    explicit AbstractDockItem(QWidget *parent = 0);
    virtual ~AbstractDockItem() = 0;

    virtual QWidget * getContents() = 0;

    virtual void setTitle(const QString &title) = 0;
    virtual void setIcon(const QString &iconPath, int size = 42) = 0;
    virtual void setMoveable(bool value) = 0;
    virtual bool moveable() = 0;
    virtual void setActived(bool value) = 0;
    virtual bool actived() = 0;
    virtual void setIndex(int value) = 0;
    virtual int index() = 0;

protected:
    QPixmap * m_appIcon = NULL;

    bool m_itemMoveable = true;
    bool m_itemActived = false;

    QString m_itemTitle = "";
    QString m_itemIconPath = "";
    int m_itemIndex = 0;

};

#endif // ABSTRACTDOCKITEM_H
