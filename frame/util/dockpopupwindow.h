#ifndef DOCKPOPUPWINDOW_H
#define DOCKPOPUPWINDOW_H

#include <darrowrectangle.h>

class DockPopupWindow : public Dtk::Widget::DArrowRectangle
{
    Q_OBJECT

public:
    explicit DockPopupWindow(QWidget *parent = 0);

    bool model() const;

    void show(const QPoint &pos, const bool model = false);

private:
    bool m_model;
};

#endif // DOCKPOPUPWINDOW_H
