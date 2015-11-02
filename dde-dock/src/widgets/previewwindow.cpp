#include "previewwindow.h"

PreviewWindow::PreviewWindow(ArrowDirection direction, QWidget *parent) : DArrowRectangle(direction, parent)
{
    setWindowFlags(Qt::X11BypassWindowManagerHint | Qt::Tool);
    setArrowWidth(ARROW_WIDTH);
    setArrowHeight(ARROW_HEIGHT);

    m_showTimer = new QTimer(this);
    m_showTimer->setSingleShot(true);
    connect(m_showTimer, &QTimer::timeout, this, &PreviewWindow::onShowTimerTriggered);

    m_hideTimer = new QTimer(this);
    m_hideTimer->setSingleShot(true);
    connect(m_hideTimer, &QTimer::timeout, this, &PreviewWindow::hide);

    m_animation = new QPropertyAnimation(this, "arrowPos");
    m_animation->setDuration(MOVE_ANIMATION_DURATION);
    m_animation->setEasingCurve(MOVE_ANIMATION_CURVE);
}

PreviewWindow::~PreviewWindow()
{

}

void PreviewWindow::showPreview(int x, int y, int interval)
{
    m_hideTimer->stop();

    if (m_showTimer->isActive())
        return;

    m_lastPos = QPoint(m_x, m_y);
    m_x = x;
    m_y = y;

    m_showTimer->start(interval);
}

void PreviewWindow::hidePreview(int interval)
{
    m_showTimer->stop();

    if (interval <= 0) {
        m_animation->stop();
        hide();
    }
    else
        m_hideTimer->start(interval);
}

void PreviewWindow::setContent(QWidget *content)
{
    m_tmpContent = content;
}

void PreviewWindow::setArrowPos(const QPoint &pos)
{
    show(pos.x(), pos.y());
}

void PreviewWindow::hide()
{

    emit hideFinish(m_lastContent);

    DArrowRectangle::hide();
}

void PreviewWindow::enterEvent(QEvent *)
{
    m_hideTimer->stop();
}

void PreviewWindow::leaveEvent(QEvent *)
{
    m_hideTimer->start();
}

void PreviewWindow::onShowTimerTriggered()
{
    if (m_lastContent != m_tmpContent)
        emit showFinish(m_lastContent);

    DArrowRectangle::setContent(m_tmpContent);
    m_lastContent = m_tmpContent;

    if (isHidden())
        show(m_x, m_y);
    else{
        m_animation->setStartValue(m_lastPos);
        m_animation->setEndValue(QPoint(m_x, m_y));
        m_animation->start();
    }
}

