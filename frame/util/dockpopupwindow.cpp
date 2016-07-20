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

DockPopupWindow::~DockPopupWindow()
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

    if (!model && !m_mouseAreaKey.isEmpty())
    {
        m_mouseInter->UnregisterArea(m_mouseAreaKey);
        m_mouseAreaKey.clear();
    }

    if (model && m_mouseAreaKey.isEmpty())
        m_mouseAreaKey = m_mouseInter->RegisterFullScreen();
}

void DockPopupWindow::hide()
{
    if (!m_mouseAreaKey.isEmpty())
    {
        m_mouseInter->UnregisterArea(m_mouseAreaKey);
        m_mouseAreaKey.clear();
    }

    DArrowRectangle::hide();
}

void DockPopupWindow::mousePressEvent(QMouseEvent *e)
{
    DArrowRectangle::mousePressEvent(e);

//    if (e->button() == Qt::LeftButton)
//        m_acceptDelayTimer->start();
}

void DockPopupWindow::globalMouseRelease(int button, int x, int y, const QString &id)
{
    Q_UNUSED(button);

    if (id != m_mouseAreaKey)
        return;

    Q_ASSERT(m_model);

    const QRect rect = QRect(pos(), size());
    const QPoint pos = QPoint(x, y);

    if (rect.contains(pos))
        return;

    emit accept();

    m_mouseInter->UnregisterArea(m_mouseAreaKey);
    m_mouseAreaKey.clear();
}
