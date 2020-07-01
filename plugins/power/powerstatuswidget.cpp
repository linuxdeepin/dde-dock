/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
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
#include "powerplugin.h"
#include <DGuiApplicationHelper>

#include <QPainter>
#include <QIcon>
#include <QMouseEvent>

DGUI_USE_NAMESPACE

PowerStatusWidget::PowerStatusWidget(QWidget *parent)
    : QWidget(parent),

      m_powerInter(new DBusPower(this))
{
//    QIcon::setThemeName("deepin");

    connect(m_powerInter, &DBusPower::BatteryPercentageChanged, this, static_cast<void (PowerStatusWidget::*)()>(&PowerStatusWidget::update));
    connect(m_powerInter, &DBusPower::BatteryStateChanged, this, static_cast<void (PowerStatusWidget::*)()>(&PowerStatusWidget::update));
    connect(m_powerInter, &DBusPower::OnBatteryChanged, this, static_cast<void (PowerStatusWidget::*)()>(&PowerStatusWidget::update));
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, [ = ] {
        refreshIcon();
    });
}

void PowerStatusWidget::refreshIcon()
{
    update();
}

void PowerStatusWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    const QPixmap icon = getBatteryIcon();
    const auto ratio = devicePixelRatioF();

    QPainter painter(this);
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(icon.rect());
    painter.drawPixmap(rf.center() - rfp.center() / ratio, icon);
}

QPixmap PowerStatusWidget::getBatteryIcon()
{
    const BatteryPercentageMap data = m_powerInter->batteryPercentage();
    const uint value = uint(qMin(100.0, qMax(0.0, data.value("Display"))));
    const int percentage = int(std::round(value));
    // onBattery应该表示的是当前是否使用电池在供电，为true表示没插入电源
    const bool plugged = !m_powerInter->onBattery();
    const BatteryState batteryState = static_cast<BatteryState>(m_powerInter->batteryState()["Display"]);

    /*根据新需求，电池电量显示分别是*/
    /* 0-5%、6-10%、11%-20%、21-30%、31-40%、41-50%、51-60%、61%-70%、71-80%、81-90%、91-100% */
    QString percentageStr;
    if (percentage <= 5) {
        percentageStr = "000";
    } else if (percentage <= 10) {
        percentageStr = "010";
    } else if (percentage <= 20) {
        percentageStr = "020";
    } else if (percentage <= 30) {
        percentageStr = "030";
    } else if (percentage <= 40) {
        percentageStr = "040";
    } else if (percentage <= 50) {
        percentageStr = "050";
    } else if (percentage <= 60) {
        percentageStr = "060";
    } else if (percentage <= 70) {
        percentageStr = "070";
    } else if (percentage <= 80) {
        percentageStr = "080";
    } else if (percentage <= 90) {
        percentageStr = "090";
    } else {
        percentageStr = "100";
    }

    QString iconStr;
    if (batteryState == BatteryState::FULLY_CHARGED && plugged) {
        iconStr = QString("battery-full-charged-symbolic");
    } else {
        iconStr = QString("battery-%1-%2")
                  .arg(percentageStr)
                  .arg(plugged ? "plugged-symbolic" : "symbolic");
    }

    if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconStr.append(PLUGIN_MIN_ICON_NAME);

    const auto ratio = devicePixelRatioF();
    QPixmap pix = QIcon::fromTheme(iconStr,
                                   QIcon::fromTheme(":/batteryicons/resources/batteryicons/" + iconStr + ".svg")).pixmap(QSize(20, 20) * ratio);
    pix.setDevicePixelRatio(ratio);

    return pix;
}

void PowerStatusWidget::resizeEvent(QResizeEvent *event)
{
    QWidget::resizeEvent(event);

    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    // 保持横纵比
    if (position == Dock::Bottom || position == Dock::Top) {
        setMaximumWidth(height());
        setMaximumHeight(QWIDGETSIZE_MAX);
    } else {
        setMaximumHeight(width());
        setMaximumWidth(QWIDGETSIZE_MAX);
    }
}
