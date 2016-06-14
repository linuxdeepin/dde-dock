#ifndef LAUNCHERITEM_H
#define LAUNCHERITEM_H

#include "dockitem.h"

class LauncherItem : public DockItem
{
    Q_OBJECT

public:
    explicit LauncherItem(QWidget *parent = 0);

private:
    void paintEvent(QPaintEvent *e);
    void resizeEvent(QResizeEvent *e);
    void mousePressEvent(QMouseEvent *e);

private:
    QPixmap m_icon;
};

#endif // LAUNCHERITEM_H
