#ifndef DOCKITEM_H
#define DOCKITEM_H

#include <QFrame>

class DockItem : public QWidget
{
    Q_OBJECT

public:
    explicit DockItem(QWidget *parent = 0);
};

#endif // DOCKITEM_H
