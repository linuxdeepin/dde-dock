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

#include <com_deepin_dde_daemon_dock.h>

#include "../interfaces/constants.h"

#include <gtest/gtest.h>

using DBusDock = com::deepin::dde::daemon::Dock;

class DockUnitTest : public QObject, public ::testing::Test
{
    Q_OBJECT

public:
    DockUnitTest();
    virtual ~DockUnitTest();
    virtual void SetUp();
    virtual void TearDown();

public:
    const DockRect dockGeometry();                               // 获取任务栏实际位置
    const DockRect frontendWindowRect();                         // 后端记录的任务栏前端界面位置(和实际位置不一定对应)
    void setPosition(Dock::Position pos);
};

#endif // DOCK_UNIT_TEST_H
