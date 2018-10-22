#include "fashiontraycontrolwidget.h"

#include <QMouseEvent>
#include <QPainter>

#define ExpandedKey "fashion-tray-expanded"

FashionTrayControlWidget::FashionTrayControlWidget(Dock::Position position, QWidget *parent)
    : QWidget(parent),
      m_settings(new QSettings("deepin", "dde-dock-tray")),
      m_dockPosition(position),
      m_expanded(m_settings->value(ExpandedKey, true).toBool()),
      m_hover(false),
      m_pressed(false)
{
    setDockPostion(m_dockPosition);
    setExpanded(m_expanded);
}

void FashionTrayControlWidget::setDockPostion(Dock::Position pos)
{
    m_dockPosition = pos;
    refreshArrowPixmap();
}

bool FashionTrayControlWidget::expanded() const
{
    return m_expanded;
}

void FashionTrayControlWidget::setExpanded(const bool &expanded)
{
    if (m_expanded == expanded) {
        return;
    }

    m_expanded = expanded;
    refreshArrowPixmap();

    m_settings->setValue(ExpandedKey, m_expanded);

    Q_EMIT expandChanged(m_expanded);
}

void FashionTrayControlWidget::paintEvent(QPaintEvent *event)
{
    QPainter painter(this);
    painter.setRenderHint(QPainter::Antialiasing, true);

    painter.setOpacity(0.5);

    if (m_expanded) {
        painter.setPen(QColor::fromRgb(40, 40, 40));
        painter.setBrush(QColor::fromRgb(40, 40, 40));
        if (m_hover) {
            painter.setPen(QColor::fromRgb(60, 60, 60));
            painter.setBrush(QColor::fromRgb(60, 60, 60));
        }
        if (m_pressed) {
            painter.setPen(QColor::fromRgb(20, 20, 20));
            painter.setBrush(QColor::fromRgb(20, 20, 20));
        }
    } else {
        painter.setPen(QColor::fromRgb(255, 255, 255));
        painter.setBrush(QColor::fromRgb(255, 255, 255));
        if (m_hover) {
            painter.setOpacity(0.6);
        }
        if (m_pressed) {
            painter.setOpacity(0.3);
        }
    }

    painter.drawRoundRect(rect());

    QRect r = QRect(QPoint(0, 0), m_arrowPix.size());
    r.moveCenter(rect().center());
    painter.drawPixmap(r, m_arrowPix);
}

void FashionTrayControlWidget::mouseReleaseEvent(QMouseEvent *event)
{
    m_pressed = false;
    update();

    if (event->button() == Qt::LeftButton) {
        event->accept();
        setExpanded(!m_expanded);
        return;
    }

    QWidget::mouseReleaseEvent(event);
}

void FashionTrayControlWidget::mousePressEvent(QMouseEvent *event)
{
    m_pressed = true;
    update();

    QWidget::mousePressEvent(event);
}

void FashionTrayControlWidget::enterEvent(QEvent *event)
{
    m_hover = true;
    update();

    QWidget::enterEvent(event);
}

void FashionTrayControlWidget::leaveEvent(QEvent *event)
{
    m_hover = false;
    update();

    QWidget::leaveEvent(event);
}

void FashionTrayControlWidget::refreshArrowPixmap()
{
    switch (m_dockPosition) {
    case Dock::Top:
    case Dock::Bottom:
        m_arrowPix.load(m_expanded ? ":/icons/resources/arrow_left_light.svg" : ":/icons/resources/arrow_right_dark.svg");
        break;
    case Dock::Left:
    case Dock::Right:
        m_arrowPix.load(m_expanded ? ":/icons/resources/arrow_up_light.svg" : ":/icons/resources/arrow_down_dark.svg");
        break;
    default:
        break;
    }

    update();
}
