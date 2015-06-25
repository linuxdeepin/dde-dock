#ifndef DOCKLAYOUT_H
#define DOCKLAYOUT_H

#include <QWidget>
#include <QList>
#include "appitem.h"

class DockLayout : public QWidget
{
    Q_OBJECT
public:
    enum Direction{
        LeftToRight,
        RightToLeft,
        TopToBottom,
        BottomToTop
    };

    enum MarginEdge{
        LeftMargin,
        RightMargin,
        TopMargin,
        BottomMargin
    };

    explicit DockLayout(QWidget *parent = 0);

    void setParent(QWidget *parent);
    void addItem(AppItem * item);
    void insertItem(AppItem *item, int index);
    void removeItem(int index);
    void moveItem(int from, int to);
    void setItemMoveable(int index, bool moveable);
    void setMargin(qreal margin);
    void setMargin(DockLayout::MarginEdge edge, qreal margin);
    void setSpacing(qreal spacing);
    void setSortDirection(DockLayout::Direction value);
    void relayout();
    int indexOf(AppItem * item);
    int indexOf(int x,int y);

signals:

public slots:

private slots:
    void slotItemDrag(int x,int y,AppItem *item);

private:
    void sortLeftToRight();
    void sortRightToLeft();
    void sortTopToBottom();
    void sortBottomToTop();

private:
    QList<AppItem *> appList;

    DockLayout::Direction sortDirection = DockLayout::LeftToRight;
    qreal itemSpacing = 10;
    qreal leftMargin = 0;
    qreal rightMargin = 0;
    qreal topMargin = 0;
    qreal bottomMargin = 0;
};

#endif // DOCKLAYOUT_H
