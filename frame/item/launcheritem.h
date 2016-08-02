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

    QWidget *popupTips();

private:
    QPixmap m_smallIcon;
    QPixmap m_largeIcon;

    QLabel *m_tips;
};

#endif // LAUNCHERITEM_H
