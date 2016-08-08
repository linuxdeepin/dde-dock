#ifndef CONTAINERITEM_H
#define CONTAINERITEM_H

#include "dockitem.h"

class ContainerItem : public DockItem
{
    Q_OBJECT

public:
    explicit ContainerItem(QWidget *parent = 0);

    inline ItemType itemType() const {return Container;}
};

#endif // CONTAINERITEM_H
