// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "powerstatuswidget.h"
#include "powerplugin.h"
#include "dbus/dbuspower.h"

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

    connect(m_powerInter, &DBusPower::BatteryPercentageChanged, this, &PowerStatusWidget::refreshIcon);
    connect(m_powerInter, &DBusPower::BatteryStateChanged, this, &PowerStatusWidget::refreshIcon);
    connect(m_powerInter, &DBusPower::OnBatteryChanged, this, &PowerStatusWidget::refreshIcon);
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &PowerStatusWidget::refreshIcon);
}

void PowerStatusWidget::refreshIcon()
{
    update();
    Q_EMIT iconChanged();
}

void PowerStatusWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    int themeType = DGuiApplicationHelper::instance()->themeType();
    if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && themeType == DGuiApplicationHelper::LightType)
        themeType = DGuiApplicationHelper::DarkType;
    const QPixmap icon = getBatteryIcon(themeType);
    const auto ratio = devicePixelRatioF();

    QPainter painter(this);
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(icon.rect());
    painter.drawPixmap(rf.center() - rfp.center() / ratio, icon);
}

QPixmap PowerStatusWidget::getBatteryIcon(int themeType)
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

    if (themeType == DGuiApplicationHelper::ColorType::LightType)
        iconStr.append(PLUGIN_MIN_ICON_NAME);

    const auto ratio = devicePixelRatioF();
    QSize pixmapSize = QCoreApplication::testAttribute(Qt::AA_UseHighDpiPixmaps) ? QSize(20, 20) : (QSize(20, 20) * ratio);
    QPixmap pix = QIcon::fromTheme(iconStr, QIcon::fromTheme(":/batteryicons/resources/batteryicons/" + iconStr + ".svg")).pixmap(pixmapSize);
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
