/*
 * Copyright (C) 2018 ~ 2028 Uniontech Technology Co., Ltd.
 *
 * Author:     fanpengcheng <fanpengcheng@uniontech.com>
 *
 * Maintainer: fanpengcheng <fanpengcheng@uniontech.com>
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
#ifndef DOCK_UNIT_TEST_H
#define DOCK_UNIT_TEST_H
#include <QObject>

class DockUnitTest : public QObject
{
    Q_OBJECT

public:
    DockUnitTest();
    ~DockUnitTest();

private slots:
    void dock_geometry_check();         // 显示区域
    void dock_position_check();         // 位置检查
    void dock_displayMode_check();      // 显示模式检查
    void dock_appItemCount_check();     // 应用显示数量检查
};

#endif // DOCK_UNIT_TEST_H
