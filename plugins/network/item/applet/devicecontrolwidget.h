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
#ifndef DEVICECONTROLWIDGET_H
#define DEVICECONTROLWIDGET_H

#include "horizontalseperator.h"

#include <QWidget>
#include <QLabel>
#include <dloadingindicator.h>
#include <dswitchbutton.h>

#include <com_deepin_daemon_airplanemode.h>

using AirplanInter = com::deepin::daemon::AirplaneMode;

DWIDGET_USE_NAMESPACE
class TipsWidget;
class QLabel;
class DeviceControlWidget : public QWidget
{
    Q_OBJECT

public:
    explicit DeviceControlWidget(QWidget *parent = 0);

    void setDeviceName(const QString &name);
    void setDeviceEnabled(const bool enable);

signals:
    void enableButtonToggled(const bool enable) const;
    void requestRefresh() const;

protected:
    bool eventFilter(QObject *watched, QEvent *event) Q_DECL_OVERRIDE;
    void refreshIcon();

private slots:
    void refreshNetwork();

private:
    QLabel *m_deviceName;

    Dtk::Widget::DSwitchButton *m_switchBtn;
    DLoadingIndicator *m_loadingIndicator;

    AirplanInter *m_airplaninter;           //飞行模式dbus接口(system dbus)  com.deepin.daemon.AirplaneMode

};

#endif // DEVICECONTROLWIDGET_H
