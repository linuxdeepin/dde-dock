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

#ifndef POWERTRAYWIDGET_H
#define POWERTRAYWIDGET_H

#include "abstracttraywidget.h"
#include "dbus/dbuspower.h"

#include <QWidget>

class PowerTrayWidget : public AbstractTrayWidget
{
    Q_OBJECT

public:
    explicit PowerTrayWidget(AbstractTrayWidget *parent = nullptr);

public:
    void setActive(const bool active) Q_DECL_OVERRIDE;
    void updateIcon() Q_DECL_OVERRIDE;
    void sendClick(uint8_t mouseButton, int x, int y) Q_DECL_OVERRIDE;
    const QImage trayImage() Q_DECL_OVERRIDE;

protected:
    QSize sizeHint() const Q_DECL_OVERRIDE;
    void paintEvent(QPaintEvent *e) Q_DECL_OVERRIDE;

private:
    DBusPower *m_powerInter;

    QPixmap m_pixmap;
};

#endif // POWERTRAYWIDGET_H
