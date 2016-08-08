#ifndef PLACEHOLDERITEM_H
#define PLACEHOLDERITEM_H

#include "dockitem.h"

class PlaceholderItem : public DockItem
{
    Q_OBJECT

public:
    explicit PlaceholderItem(QWidget *parent = 0);

    // fake as app item
    inline ItemType itemType() const {return App;}
};

#endif // PLACEHOLDERITEM_H
