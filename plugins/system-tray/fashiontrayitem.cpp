#include "fashiontrayitem.h"

#include <QPainter>
#include <QDebug>
#include <QMouseEvent>

#include <xcb/xproto.h>

#define DRAG_THRESHOLD  10

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
    const int circleSize = std::min(r.width(), r.height()) * 0.8 / 2 - 3;
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

void FashionTrayItem::mousePressEvent(QMouseEvent *e)
{
    QWidget::mousePressEvent(e);

    emit requestPopupApplet();

    m_pressPoint = e->pos();
}

void FashionTrayItem::mouseReleaseEvent(QMouseEvent *e)
{
    const QPoint point = e->pos() - m_pressPoint;

    if (point.manhattanLength() > DRAG_THRESHOLD)
        return;

    if (!m_activeTray)
        return;

    QPoint globalPos = mapToGlobal(QPoint(0, 0));
    uint8_t buttonIndex = XCB_BUTTON_INDEX_1;

    switch (e->button()) {
    case Qt:: MiddleButton:
        buttonIndex = XCB_BUTTON_INDEX_2;
        break;
    case Qt::RightButton:
        buttonIndex = XCB_BUTTON_INDEX_3;
        break;
    default:
        break;
    }

    m_activeTray->sendClick(buttonIndex, globalPos.x(), globalPos.y());
}
