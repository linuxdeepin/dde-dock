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

#ifndef ACCESSPOINTWIDGET_H
#define ACCESSPOINTWIDGET_H

#include "accesspoint.h"

#include <QWidget>
#include <QLabel>
#include <QPushButton>
#include <QDBusObjectPath>

#include <dimagebutton.h>
#include <NetworkDevice>

class AccessPointWidget : public QFrame
{
    Q_OBJECT
    Q_PROPERTY(bool active READ active DESIGNABLE true)

public:
    explicit AccessPointWidget();

    const AccessPoint ap() const { return m_ap; }
    void updateAP(const AccessPoint &ap);

    bool active() const;
    void setActiveState(const dde::network::NetworkDevice::DeviceStatus state);

signals:
    void requestActiveAP(const QString &apPath, const QString &ssid) const;
    void requestDeactiveAP(const AccessPoint &ap) const;
    void clicked() const;

private:
    void enterEvent(QEvent *e);
    void leaveEvent(QEvent *e);
    void setStrengthIcon(const int strength);

private slots:
    void ssidClicked();
    void disconnectBtnClicked();

private:
    dde::network::NetworkDevice::DeviceStatus m_activeState;

    AccessPoint m_ap;
    QPushButton *m_ssidBtn;
    Dtk::Widget::DImageButton *m_disconnectBtn;
    QLabel *m_securityLabel;
    QLabel *m_strengthLabel;

    QPixmap m_securityPixmap;
    QSize m_securityIconSize;
};

#endif // ACCESSPOINTWIDGET_H
