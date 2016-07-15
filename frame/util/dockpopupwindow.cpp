#include "dockpopupwindow.h"

DWIDGET_USE_NAMESPACE

DockPopupWindow::DockPopupWindow(QWidget *parent)
    : DArrowRectangle(ArrowBottom, parent),
      m_model(false)
{

}

bool DockPopupWindow::model() const
{
    return m_model;
}

void DockPopupWindow::show(const QPoint &pos, const bool model)
{
    m_model = model;

    DArrowRectangle::show(pos.x(), pos.y());
}
