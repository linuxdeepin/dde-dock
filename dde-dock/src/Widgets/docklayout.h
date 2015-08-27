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
        RightToLeft,
        TopToBottom
    };

    explicit DockLayout(QWidget *parent = 0);

    void addItem(AbstractDockItem *item);
    void insertItem(AbstractDockItem *item, int index);
    void moveItem(int from, int to);
    void removeItem(int index);
    void removeItem(AbstractDockItem *item);
    void setSpacing(qreal spacing);
    void setVerticalAlignment(Qt::Alignment value);
    void setSortDirection(DockLayout::Direction value);

    int indexOf(AbstractDockItem *item) const;
    int indexOf(int x,int y) const;
    int getContentsWidth();
    int getItemCount() const;
    QList<AbstractDockItem *> getItemList() const;

signals:
    void startDrag();
    void itemDropped();
    void contentsWidthChange();
    void frameUpdate();

public slots:
    void removeSpacingItem();
    void restoreTmpItem();
    void clearTmpItem();
    void relayout();

protected:
    bool eventFilter(QObject *obj, QEvent *event);
    void dragEnterEvent(QDragEnterEvent *event);
    void dropEvent(QDropEvent *event);

private slots:
    void slotItemDrag();
    void slotItemRelease();
    void slotItemEntered(QDragEnterEvent *event);
    void slotItemExited(QDragLeaveEvent *event);
    void slotAnimationFinish();

private:
    void sortLeftToRight();
    void sortTopToBottom();
    void leftToRightMove(int hoverIndex);
    void topToBottomMove(int hoverIndex);
    void addSpacingItem();
    void dragoutFromLayout(int index);

    int spacingItemWidth() const;
    int spacingItemIndex() const;
    int animatingItemCount();
    QStringList itemsIdList() const;

private:
    QList<AbstractDockItem *> m_appList;
    QMap<AbstractDockItem *,int> m_dragItemMap;//only one item inside
    DBusDockedAppManager *m_ddam = new DBusDockedAppManager(this);

    Qt::Alignment m_verticalAlignment = Qt::AlignVCenter;
    DockLayout::Direction m_sortDirection = DockLayout::LeftToRight;

    qreal m_itemSpacing = 10;
    QPoint m_lastPost = QPoint(0,0);
    int m_lastHoverIndex = -1;
    bool m_movingLeftward = true;

    const int MOVE_ANIMATION_DURATION_BASE = 300;
};

#endif // DOCKLAYOUT_H
