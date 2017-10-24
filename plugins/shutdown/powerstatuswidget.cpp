/*
 * Copyright (C) 2011 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include "powerstatuswidget.h"

#include <QPainter>
#include <QIcon>
#include <QMouseEvent>

PowerStatusWidget::PowerStatusWidget(QWidget *parent)
    : QWidget(parent),

      m_powerInter(new DBusPower(this))
{
//    QIcon::setThemeName("deepin");

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
    const auto ratio = devicePixelRatioF();

    QPainter painter(this);
    painter.drawPixmap(rect().center() - icon.rect().center() / ratio, icon);
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
    const BatteryPercentageMap data = m_powerInter->batteryPercentage();
    const uint value = qMin(100.0, qMax(0.0, data.value("Display")));
    const int percentage = std::round(value);
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

    const auto ratio = devicePixelRatioF();
    QPixmap pix = QIcon::fromTheme(iconStr).pixmap(QSize(16, 16) * ratio);
    pix.setDevicePixelRatio(ratio);

    return pix;
}
