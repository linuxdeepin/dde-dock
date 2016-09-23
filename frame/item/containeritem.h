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

    void setDropping(const bool dropping);
    void addItem(DockItem * const item);
    void removeItem(DockItem * const item);
    bool contains(DockItem * const item);

public slots:
    void refershIcon();

protected:
    void dragEnterEvent(QDragEnterEvent *e);
    void dragMoveEvent(QDragMoveEvent *e);
    void paintEvent(QPaintEvent *e);
    void mouseReleaseEvent(QMouseEvent *e);
    QSize sizeHint() const;

private:
    bool m_dropping;
    ContainerWidget *m_containerWidget;
    QPixmap m_icon;
};

#endif // CONTAINERITEM_H
