#ifndef DOCKLAYOUT_H
#define DOCKLAYOUT_H

#include <QWidget>
#include <QList>
#include <QMap>
#include <QPropertyAnimation>
#include <QCursor>
#include <QJsonDocument>
#include <QJsonObject>
#include "Controller/dockmodedata.h"
#include "DBus/dbusdockedappmanager.h"
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
    void setSpacing(qreal spacing);
    void setVerticalAlignment(Qt::Alignment value);
    void setSortDirection(DockLayout::Direction value);
    int indexOf(AbstractDockItem * item);
    int indexOf(int x,int y);
    int getContentsWidth();
    int getItemCount();
    QList<AbstractDockItem *> getItemList() const;

public slots:
    void restoreTmpItem();
    void relayout();
    void clearTmpItem();

signals:
    void dragStarted();
    void itemDropped();
    void contentsWidthChange();
    void frameUpdate();

protected:
    void dragEnterEvent(QDragEnterEvent *event);
    void dropEvent(QDropEvent *event);

private slots:
    void slotItemDrag();
    void slotItemRelease();
    void slotItemEntered(QDragEnterEvent * event);
    void slotItemExited(QDragLeaveEvent *event);
    void slotAnimationFinish();

private:
    void sortLeftToRight();
    void sortRightToLeft();
    void leftToRightMove(int hoverIndex);
    void rightToLeftMove(int hoverIndex);
    void addSpacingItem();
    void dragoutFromLayout(int index);

    int spacingItemWidth();
    int spacingItemIndex();
    QStringList itemsIdList();

private:
    QList<AbstractDockItem *> m_appList;
    QMap<AbstractDockItem *,int> m_dragItemMap;//only one item inside
    DBusDockedAppManager *m_ddam = new DBusDockedAppManager(this);

    Qt::Alignment m_verticalAlignment = Qt::AlignVCenter;
    DockLayout::Direction m_sortDirection = DockLayout::LeftToRight;

    qreal m_itemSpacing = 10;
    QPoint m_lastPost = QPoint(0,0);
    int m_lastHoverIndex = -1;
    int m_animationItemCount = 0;
    bool m_movingForward = false;

    const int MOVE_ANIMATION_DURATION_BASE = 300;
};

#endif // DOCKLAYOUT_H
