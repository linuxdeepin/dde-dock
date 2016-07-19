#include "powerstatuswidget.h"

#include <QPainter>
#include <QIcon>

PowerStatusWidget::PowerStatusWidget(QWidget *parent)
    : QWidget(parent),

      m_powerInter(new DBusPower(this))
{
    QIcon::setThemeName("deepin");
}

QSize PowerStatusWidget::sizeHint() const
{
    return QSize(24, 24);
}

void PowerStatusWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    const QPixmap icon = getBatteryIcon();

    QPainter painter(this);
    painter.drawPixmap(rect().center() - icon.rect().center(), icon);
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
