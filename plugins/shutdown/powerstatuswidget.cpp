#include "powerstatuswidget.h"

#include <QPainter>
#include <QIcon>
#include <QMouseEvent>

PowerStatusWidget::PowerStatusWidget(QWidget *parent)
    : QWidget(parent),

      m_powerInter(new DBusPower(this))
{
    QIcon::setThemeName("deepin");

    connect(m_powerInter, &DBusPower::BatteryPercentageChanged, this, static_cast<void (PowerStatusWidget::*)()>(&PowerStatusWidget::update));
    connect(m_powerInter, &DBusPower::BatteryStateChanged, this, static_cast<void (PowerStatusWidget::*)()>(&PowerStatusWidget::update));
    connect(m_powerInter, &DBusPower::OnBatteryChanged, this, static_cast<void (PowerStatusWidget::*)()>(&PowerStatusWidget::update));
}

QSize PowerStatusWidget::sizeHint() const
{
    return QSize(26, 26);
}

void PowerStatusWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    const QPixmap icon = getBatteryIcon();

    QPainter painter(this);
    painter.drawPixmap(rect().center() - icon.rect().center(), icon);
}

void PowerStatusWidget::mousePressEvent(QMouseEvent *e)
{
    if (e->button() == Qt::LeftButton)
        return QWidget::mousePressEvent(e);

    const QPoint p(e->pos() - rect().center());
    if (p.manhattanLength() < std::min(width(), height()) * 0.8 * 0.5)
    {
        emit requestContextMenu(POWER_KEY);
        return;
    }

    return QWidget::mousePressEvent(e);
}

QPixmap PowerStatusWidget::getBatteryIcon()
{
    const int percentage = m_powerInter->batteryPercentage()["Display"];
    const bool plugged = !m_powerInter->onBattery();

    QString percentageStr;
    if (percentage < 10 && percentage >= 0) {
        percentageStr = "000";
    } else if (percentage < 30) {
        percentageStr = "020";
    } else if (percentage < 50) {
        percentageStr = "040";
    } else if (percentage < 70) {
        percentageStr = "060";
    } else if (percentage < 90) {
        percentageStr = "080";
    } else if (percentage <= 100){
        percentageStr = "100";
    } else {
        percentageStr = "000";
    }

    QString iconStr;

    if (plugged) {
        iconStr = "battery-charged-symbolic";
    } else {
        iconStr = QString("battery-%1-symbolic").arg(percentageStr);
    }

    return QIcon::fromTheme(iconStr).pixmap(16, 16);
}
