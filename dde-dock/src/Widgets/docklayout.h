#ifndef DOCKLAYOUT_H
#define DOCKLAYOUT_H

#include <QWidget>
#include <QList>
#include <QMap>
#include <QPropertyAnimation>
#include <QCursor>
#include "appitem.h"

class DockLayout : public QWidget
{
    Q_OBJECT
public:
    enum Direction{
        LeftToRight,
        RightToLeft
    };

    explicit DockLayout(QWidget *parent = 0);

    void addItem(AbstractDockItem * item);
    void insertItem(AbstractDockItem *item, int index);
    void removeItem(int index);
    void moveItem(int from, int to);
    void setSpacing(qreal spacing);
    void setSortDirection(DockLayout::Direction value);
    int indexOf(AbstractDockItem * item);
    int indexOf(int x,int y);
    int getContentsWidth();
    int getItemCount();
    QList<AbstractDockItem *> getItemList() const;

public slots:
    void relayout();

signals:
    void dragStarted();
    void itemDropped();
    void contentsWidthChange();

protected:
    void dragEnterEvent(QDragEnterEvent *event);
    void dropEvent(QDropEvent *event);

private slots:
    void slotItemDrag();
    void slotItemRelease(int x, int y);
    void slotItemEntered(QDragEnterEvent * event);
    void slotItemExited(QDragLeaveEvent *event);

private:
    void sortLeftToRight();
    void sortRightToLeft();

    void addSpacingItem();
    void dragoutFromLayout(int index);
    bool hasSpacingItemInList();
    int spacingItemIndex();

    void moveWithSpacingItem(int hoverIndex);
private:
    QList<AbstractDockItem *> appList;
    QMap<AbstractDockItem *,int> tmpAppMap;//only one item inside

    DockLayout::Direction sortDirection = DockLayout::LeftToRight;
    qreal itemSpacing = 10;

    bool movingForward = false;
    int lastHoverIndex = 0;
    int m_animationItemCount = 0;
    QPoint m_lastPost = QPoint(0,0);
};

#endif // DOCKLAYOUT_H
