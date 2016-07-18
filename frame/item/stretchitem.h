#ifndef STRETCHITEM_H
#define STRETCHITEM_H

#include "dockitem.h"

class StretchItem : public DockItem
{
    Q_OBJECT

public:
    explicit StretchItem(QWidget *parent = 0);

private:
    void mousePressEvent(QMouseEvent *e);
};

#endif // STRETCHITEM_H
