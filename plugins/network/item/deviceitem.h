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

#include "networkmanager.h"

#include <QWidget>

class DeviceItem : public QWidget
{
    Q_OBJECT

public:
    explicit DeviceItem(const QString &path);

    const QString path() const { return m_devicePath; }

    virtual NetworkDevice::NetworkType type() const = 0;
    virtual NetworkDevice::NetworkState state() const = 0;
    virtual void refreshIcon() = 0;
    virtual const QString itemCommand() const;
    virtual const QString itemContextMenu();
    virtual QWidget *itemApplet();
    virtual QWidget *itemPopup();
    virtual void invokeMenuItem(const QString &menuId);

signals:
    void requestContextMenu() const;

protected:
    bool enabled() const;
    void setEnabled(const bool enable);
    QSize sizeHint() const;

protected:
    const QString m_devicePath;

    NetworkManager *m_networkManager;
};

#endif // DEVICEITEM_H
