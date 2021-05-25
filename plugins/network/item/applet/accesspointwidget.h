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

#include <NetworkDevice>

#include <QWidget>
#include <QLabel>
#include <QDBusObjectPath>

class StateButton;
class SsidButton : public QLabel
{
    Q_OBJECT
public:
    explicit SsidButton(QWidget *parent = nullptr)
        : QLabel(parent) {}
    virtual ~SsidButton() override {}

signals:
    void clicked();

protected:
    void mouseReleaseEvent(QMouseEvent *event) override
    {
        QLabel::mouseReleaseEvent(event);

        Q_EMIT clicked();
    }
};

class AccessPointWidget : public QFrame
{
    Q_OBJECT
    Q_PROPERTY(bool active READ active DESIGNABLE true)

public:
    explicit AccessPointWidget(QWidget *parent = nullptr);

    const AccessPoint ap() const { return m_ap; }
    void updateAP(const AccessPoint &ap);

    bool active() const;
    void setActiveState(const dde::network::NetworkDevice::DeviceStatus state);

signals:
    void requestActiveAP(const QString &apPath, const QString &ssid) const;
    void requestDeactiveAP(const AccessPoint &ap) const;
    void clicked() const;

private:
    void enterEvent(QEvent *e) override;
    void leaveEvent(QEvent *e) override;
    void setStrengthIcon(const int strength);

protected:
    void paintEvent(QPaintEvent *event) override;

private slots:
    void ssidClicked();
    void disconnectBtnClicked();

private:
    dde::network::NetworkDevice::DeviceStatus m_activeState;

    AccessPoint m_ap;
    SsidButton *m_ssidBtn;
    QLabel *m_securityLabel;
    QLabel *m_strengthLabel;
    StateButton *m_stateButton;

    QPixmap m_securityPixmap;
    QSize m_securityIconSize;

    bool m_isEnter;
};

#endif // ACCESSPOINTWIDGET_H
