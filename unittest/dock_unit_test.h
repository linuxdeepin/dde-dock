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

using DBusDock = com::deepin::dde::daemon::Dock;

class DockUnitTest : public QObject
{
    Q_OBJECT

public:
    DockUnitTest();
    ~DockUnitTest();

private:
    QDBusInterface *m_dockInter;
    DBusDock *m_daemonDockInter;

private:
    const DockRect dockGeometry();              // 获取任务栏实际位置
    const DockRect frontendWindowRect();        // 后端记录的任务栏前端界面位置(和实际位置不一定对应)
    void setPosition(Dock::Position pos);

private slots:
    void dock_defaultGsettings_check();                         // 默认配置项检查
    void dock_geometry_check();         // 显示区域
    void dock_position_check();         // 位置检查
    void dock_displayMode_check();      // 显示模式检查
    void dock_appItemCount_check();     // 应用显示数量检查
    void dock_defaultVolume_Check(float defaultVolume = 50.0f);  // 设备默认音量检查
    void dock_frontWindowRect_check();  // 检查FrontendWindowRect接口数据是否正确
    void dock_multi_process(); // 检查是否正常启动
    void dock_coreDump_check();     // dock是否一直崩溃
    void dock_appIconSize_check();                              // 图标大小检查
};

#endif // DOCK_UNIT_TEST_H
