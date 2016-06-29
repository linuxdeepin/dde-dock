#include "fashiontrayitem.h"

#include <QPainter>
#include <QDebug>

const double pi = std::acos(-1);

FashionTrayItem::FashionTrayItem(QWidget *parent)
    : QWidget(parent),
      m_activeTray(nullptr)
{

}

TrayWidget *FashionTrayItem::activeTray()
{
    return m_activeTray;
}

void FashionTrayItem::setActiveTray(TrayWidget *tray)
{
    m_activeTray = tray;
    update();
}

void FashionTrayItem::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    const QRect r = rect();

    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing, true);

    // draw circle
    QPen circlePen(QColor(0, 164, 233));
    circlePen.setWidth(3);
    const int circleSize = r.width() * 0.6 / 2;
    painter.setPen(circlePen);
    painter.drawEllipse(r.center(), circleSize, circleSize);

    // draw red dot
    const int offset = std::sin(pi / 4) * circleSize;
    painter.setPen(Qt::transparent);
    painter.setBrush(QColor(250, 64, 151));
    painter.drawEllipse(r.center() + QPoint(offset, -offset), 5, 5);

    // draw active icon
    if (m_activeTray)
    {
        const QImage image = m_activeTray->trayImage();
        painter.drawImage(r.center().x() - image.width() / 2, r.center().y() - image.height() / 2, image);
    }
}
