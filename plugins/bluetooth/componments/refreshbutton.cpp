#include "refreshbutton.h"
#include "imageutil.h"

#include <QTimer>
#include <QPainter>
#include <QIcon>
#include <QMouseEvent>
#include <QPropertyAnimation>
#include <QDebug>

RefreshButton::RefreshButton(QWidget *parent)
    : QWidget(parent)
    , m_refreshTimer(new QTimer(this))
    , m_rotateAngle(0)
{
    setAccessibleName("RefreshButton");
    m_refreshTimer->setInterval(500 / 60);
    initConnect();
}

void RefreshButton::setRotateIcon(QString path)
{
    m_pixmap = ImageUtil::loadSvg(path, ":/", qMin(width(), height()), devicePixelRatio());
}

void RefreshButton::startRotate()
{
    m_refreshTimer->start();
    if (m_rotateAngle == 360) {
        m_rotateAngle = 0;
    }
    m_rotateAngle += 360 / 60;
    update();
}

void RefreshButton::stopRotate()
{
    m_refreshTimer->stop();
    m_rotateAngle = 0;
    update();
}

void RefreshButton::paintEvent(QPaintEvent *e)
{
    QPainter painter(this);
    painter.setPen(Qt::NoPen);
    painter.setBrush(Qt::NoBrush);
    painter.setRenderHints(QPainter::Antialiasing | QPainter::SmoothPixmapTransform);

    painter.translate(this->width() / 2, this->height() / 2);
    painter.rotate(m_rotateAngle);
    painter.translate(-(this->width() / 2), -(this->height() / 2));
    painter.drawPixmap(this->rect(), m_pixmap);

    QWidget::paintEvent(e);
}

void RefreshButton::mousePressEvent(QMouseEvent *event)
{
    m_pressPos = event->pos();
    return QWidget::mousePressEvent(event);
}

void RefreshButton::mouseReleaseEvent(QMouseEvent *event)
{
    if (rect().contains(m_pressPos) && rect().contains(event->pos()) && !m_refreshTimer->isActive())
        Q_EMIT clicked();
    return QWidget::mouseReleaseEvent(event);
}

void RefreshButton::initConnect()
{
    connect(m_refreshTimer, &QTimer::timeout, this, &RefreshButton::startRotate);
}
