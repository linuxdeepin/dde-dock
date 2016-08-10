#ifndef CONTAINERITEM_H
#define CONTAINERITEM_H

#include "dockitem.h"
#include "components/containerwidget.h"

#include <QPixmap>

class ContainerItem : public DockItem
{
    Q_OBJECT

public:
    explicit ContainerItem(QWidget *parent = 0);

    inline ItemType itemType() const {return Container;}

    void addItem(DockItem * const item);

protected:
    void dragEnterEvent(QDragEnterEvent *e);
    void dragMoveEvent(QDragMoveEvent *e);
    void paintEvent(QPaintEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);
    QSize sizeHint() const;

private:
    QPixmap m_icon;

    ContainerWidget *m_containerWidget;
};

#endif // CONTAINERITEM_H
