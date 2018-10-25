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

#ifndef ABSTRACTNETWORKTRAYWIDGET_H
#define ABSTRACTNETWORKTRAYWIDGET_H

#include "../../abstractsystemtraywidget.h"

#include <QWidget>

#include <NetworkDevice>

class AbstractNetworkTrayWidget : public AbstractSystemTrayWidget
{
    Q_OBJECT

public:
    explicit AbstractNetworkTrayWidget(dde::network::NetworkDevice *device, QWidget *parent = nullptr);

    void setActive(const bool active) Q_DECL_OVERRIDE = 0;
    void updateIcon() Q_DECL_OVERRIDE = 0;
    const QImage trayImage() Q_DECL_OVERRIDE = 0;

    const QString contextMenu() const Q_DECL_OVERRIDE;
    void invokedMenuItem(const QString &menuId, const bool checked) Q_DECL_OVERRIDE;

    const QString &path() const { return m_path; }
    inline const dde::network::NetworkDevice * device() { return m_device; }

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
