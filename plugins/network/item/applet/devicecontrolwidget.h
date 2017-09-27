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

#ifndef DEVICECONTROLWIDGET_H
#define DEVICECONTROLWIDGET_H

#include "horizontalseperator.h"

#include <QWidget>
#include <QLabel>

#include <dswitchbutton.h>

class RefreshButton;
class DeviceControlWidget : public QWidget
{
    Q_OBJECT

public:
    explicit DeviceControlWidget(QWidget *parent = 0);

    void setDeviceName(const QString &name);
    void setDeviceEnabled(const bool enable);
//    void setSeperatorVisible(const bool visible);

signals:
    void deviceEnableChanged(const bool enable) const;
    void requestRefresh() const;

private:
    QLabel *m_deviceName;
    Dtk::Widget::DSwitchButton *m_switchBtn;
//    HorizontalSeperator *m_seperator;
    RefreshButton *m_refreshBtn;
};

#endif // DEVICECONTROLWIDGET_H
