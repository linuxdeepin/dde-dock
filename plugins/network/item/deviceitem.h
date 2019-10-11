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

#ifndef DEVICEITEM_H
#define DEVICEITEM_H

#include <QWidget>
#include <QPointer>

#include <NetworkDevice>

class DeviceItem : public QWidget
{
    Q_OBJECT

public:
    explicit DeviceItem(dde::network::NetworkDevice *device);

    const QString &path() const { return m_path; }

    inline const QPointer<dde::network::NetworkDevice> device() { return m_device; }

    virtual void refreshIcon() = 0;
    virtual const QString itemCommand() const;
    virtual const QString itemContextMenu();
    virtual QWidget *itemApplet();
    virtual QWidget *itemTips();
    virtual void invokeMenuItem(const QString &menuId);

signals:
    void requestSetDeviceEnable(const QString &path, const bool enable) const;

protected:
    QSize sizeHint() const;

protected:
    QPointer<dde::network::NetworkDevice> m_device;
    QString m_pageName;

private:
    QString m_path;
};

#endif // DEVICEITEM_H
