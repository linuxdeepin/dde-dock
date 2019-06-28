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

#ifndef POWERSTATUSWIDGET_H
#define POWERSTATUSWIDGET_H

#include "../frame/dbus/dbuspower.h"

#include <QWidget>

#define POWER_KEY "power"

class PowerStatusWidget : public QWidget
{
    Q_OBJECT

public:
    explicit PowerStatusWidget(QWidget *parent = 0);

    void refreshIcon();

signals:
    void requestContextMenu(const QString &itemKey) const;

protected:
    QSize sizeHint() const;
    void paintEvent(QPaintEvent *e);

private:
    QPixmap getBatteryIcon();

private:
    DBusPower *m_powerInter;
};

#endif // POWERSTATUSWIDGET_H
