#include "dockpopupwindow.h"

DWIDGET_USE_NAMESPACE

DockPopupWindow::DockPopupWindow(QWidget *parent)
    : DArrowRectangle(ArrowBottom, parent),
      m_model(false),

      m_acceptDelayTimer(new QTimer(this)),

      m_mouseInter(new DBusXMouseArea(this))
{
    m_acceptDelayTimer->setSingleShot(true);
    m_acceptDelayTimer->setInterval(100);

    connect(m_acceptDelayTimer, &QTimer::timeout, this, &DockPopupWindow::accept);
    connect(m_mouseInter, &DBusXMouseArea::ButtonRelease, this, &DockPopupWindow::globalMouseRelease);
}

bool DockPopupWindow::model() const
{
    return m_model;
}

void DockPopupWindow::show(const QPoint &pos, const bool model)
{
    m_model = model;

    DArrowRectangle::show(pos.x(), pos.y());

    if (model)
        m_mouseAreaKey = m_mouseInter->RegisterFullScreen();
}

void DockPopupWindow::mousePressEvent(QMouseEvent *e)
{
    DArrowRectangle::mousePressEvent(e);

    if (e->button() == Qt::LeftButton)
        m_acceptDelayTimer->start();
}

void DockPopupWindow::globalMouseRelease()
{
    if (!m_model)
        return;

    const QRect rect = QRect(pos(), size());
    const QPoint pos = QCursor::pos();

    if (rect.contains(pos))
        return;

    emit accept();

    m_mouseInter->UnregisterArea(m_mouseAreaKey);
}
