/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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

#include "powertraywidget.h"

#include <QPainter>
#include <QIcon>
#include <QMouseEvent>

#define BATTERY_DISCHARED   2
#define BATTERY_FULL        4

PowerTrayWidget::PowerTrayWidget(QWidget *parent)
    : AbstractSystemTrayWidget(parent),
      m_powerInter(new DBusPower(this)),
      m_tipsLabel(new TipsWidget)
{
    m_tipsLabel->setVisible(false);

    connect(m_powerInter, &DBusPower::BatteryPercentageChanged, this, &PowerTrayWidget::updateIcon);
    connect(m_powerInter, &DBusPower::BatteryStateChanged, this, &PowerTrayWidget::updateIcon);
    connect(m_powerInter, &DBusPower::OnBatteryChanged, this, &PowerTrayWidget::updateIcon);

    updateIcon();
}

void PowerTrayWidget::setActive(const bool active)
{

}

void PowerTrayWidget::updateIcon()
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

    const QString iconStr = QString("battery-%1-%2")
                                .arg(percentageStr)
                                .arg(plugged ? "plugged-symbolic" : "symbolic");
    const auto ratio = devicePixelRatioF();
    m_pixmap = QIcon::fromTheme(iconStr).pixmap(QSize(16, 16) * ratio);
    m_pixmap.setDevicePixelRatio(ratio);

    update();
}

const QImage PowerTrayWidget::trayImage()
{
    return m_pixmap.toImage();
}

QWidget *PowerTrayWidget::trayTipsWidget()
{
    const BatteryPercentageMap data = m_powerInter->batteryPercentage();

    if (data.isEmpty())
    {
        m_tipsLabel->setText(tr("Shut down"));
        return m_tipsLabel;
    }

    const uint percentage = qMin(100.0, qMax(0.0, data.value("Display")));
    const QString value = QString("%1%").arg(std::round(percentage));
    const bool charging = !m_powerInter->onBattery();
    if (!charging)
        m_tipsLabel->setText(tr("Remaining Capacity %1").arg(value));
    else
    {
        const int batteryState = m_powerInter->batteryState()["Display"];

        if (batteryState == BATTERY_FULL || percentage == 100.)
            m_tipsLabel->setText(tr("Charged %1").arg(value));
        else
            m_tipsLabel->setText(tr("Charging %1").arg(value));
    }

    return m_tipsLabel;
}

const QString PowerTrayWidget::trayClickCommand()
{
    return QString("dbus-send --print-reply --dest=com.deepin.dde.ControlCenter /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.ShowModule \"string:power\"");
}

const QString PowerTrayWidget::contextMenu() const
{
    QList<QVariant> items;
    items.reserve(6);

    QMap<QString, QVariant> power;
    power["itemId"] = "power";
    power["itemText"] = tr("Power settings");
    power["isActive"] = true;
    items.push_back(power);

    QMap<QString, QVariant> menu;
    menu["items"] = items;
    menu["checkableMenu"] = false;
    menu["singleCheck"] = false;

    return QJsonDocument::fromVariant(menu).toJson();
}

void PowerTrayWidget::invokedMenuItem(const QString &menuId, const bool checked)
{
    Q_UNUSED(checked)

    if (menuId == "power")
        QProcess::startDetached("dbus-send --print-reply --dest=com.deepin.dde.ControlCenter /com/deepin/dde/ControlCenter com.deepin.dde.ControlCenter.ShowModule \"string:power\"");
}

QSize PowerTrayWidget::sizeHint() const
{
    return QSize(26, 26);
}

void PowerTrayWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    QPainter painter(this);

    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(m_pixmap.rect());
    const QPointF &p = rf.center() - rfp.center() / m_pixmap.devicePixelRatioF();
    painter.drawPixmap(p, m_pixmap);
}
