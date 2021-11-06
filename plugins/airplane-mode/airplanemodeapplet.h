/*
 * Copyright (C) 2020 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     weizhixiang <weizhixiang@uniontech.com>
 *
 * Maintainer: weizhixiang <weizhixiang@uniontech.com>
 *
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

#ifndef AIRPLANEMODEAPPLET_H
#define AIRPLANEMODEAPPLET_H

#include <QWidget>

namespace Dtk {
namespace Widget {
class DSwitchButton;
}
}

class QLabel;
class AirplaneModeApplet : public QWidget
{
    Q_OBJECT

public:
    explicit AirplaneModeApplet(QWidget *parent = nullptr);
    void setEnabled(bool enable);

signals:
    void enableChanged(bool enable);

private:
    Dtk::Widget::DSwitchButton *m_switchBtn;
};

#endif // AIRPLANEMODEAPPLET_H
