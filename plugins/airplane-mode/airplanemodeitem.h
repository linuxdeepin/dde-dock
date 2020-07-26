/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     zhaolong<zhaolong@uniontech.com>
 *
 * Maintainer:  xiehui<xiehui@uniontech.com>
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

#ifndef AIRPLANEMODEITEM_H
#define AIRPLANEMODEITEM_H

#include <com_deepin_daemon_airplanemode.h>

#include <QWidget>

using DBusAirplaneMode = com::deepin::daemon::AirplaneMode;

namespace Dock {
class TipsWidget;
}

class AirplaneModeApplet;
class AirplaneModeItem : public QWidget
{
    Q_OBJECT
public:
    explicit AirplaneModeItem(QWidget *parent = nullptr);

    QWidget *tipsWidget();
    QWidget *popupApplet();
    const QString contextMenu() const;
    void invokeMenuItem(const QString menuId, const bool checked);
    void refreshIcon();

protected:
    void resizeEvent(QResizeEvent *e);
    void paintEvent(QPaintEvent *e);

private:
    Dock::TipsWidget *m_tipsLabel;
    AirplaneModeApplet *m_applet;
    DBusAirplaneMode *m_airplaneModeInter;
    QPixmap m_iconPixmap;
};

#endif // AIRPLANEMODEITEM_H
