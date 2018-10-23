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

#ifndef ABSTRACTNETWORKTRAYWIDGET_H
#define ABSTRACTNETWORKTRAYWIDGET_H

#include "abstracttraywidget.h"

#include <QWidget>

#include <NetworkDevice>

class AbstractNetworkTrayWidget : public AbstractTrayWidget
{
    Q_OBJECT

public:
    explicit AbstractNetworkTrayWidget(dde::network::NetworkDevice *device, QWidget *parent = nullptr);

public:
    const QString &path() const { return m_path; }
    inline const dde::network::NetworkDevice * device() { return m_device; }

    virtual const QString itemContextMenu();
    virtual void invokeMenuItem(const QString &menuId);

Q_SIGNALS:
    void requestSetDeviceEnable(const QString &path, const bool enable) const;

protected:
    QSize sizeHint() const Q_DECL_OVERRIDE;

protected:
    dde::network::NetworkDevice *m_device;

private:
    QString m_path;
};

#endif // ABSTRACTNETWORKTRAYWIDGET_H
