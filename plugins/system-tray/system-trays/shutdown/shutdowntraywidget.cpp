#include "shutdowntraywidget.h"

#include <QSvgRenderer>
#include <QPainter>
#include <QMouseEvent>
#include <QApplication>

ShutdownTrayWidget::ShutdownTrayWidget(AbstractTrayWidget *parent)
    : AbstractTrayWidget(parent)
{
    updateIcon();
}

void ShutdownTrayWidget::setActive(const bool active)
{

}

void ShutdownTrayWidget::updateIcon()
{
    const auto ratio = qApp->devicePixelRatio();

    QPixmap pixmap(QSize(16, 16) * ratio);
    QSvgRenderer renderer(QString(":/icons/system-trays/shutdown/resources/icons/normal.svg"));
    pixmap.fill(Qt::transparent);

    QPainter painter;
    painter.begin(&pixmap);
    renderer.render(&painter);
    painter.end();

    pixmap.setDevicePixelRatio(ratio);

    m_pixmap = pixmap;

    update();
}

void ShutdownTrayWidget::sendClick(uint8_t mouseButton, int x, int y)
{

}

const QImage ShutdownTrayWidget::trayImage()
{
    return m_pixmap.toImage();
}

QSize ShutdownTrayWidget::sizeHint() const
{
    return QSize(26, 26);
}

void ShutdownTrayWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    QPainter painter(this);
    painter.drawPixmap(rect().center() - m_pixmap.rect().center() / qApp->devicePixelRatio(), m_pixmap);
}
