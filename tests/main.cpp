/*
 * Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
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
#include <gtest/gtest.h>
#ifdef QT_DEBUG
#include <sanitizer/asan_interface.h>
#endif

#include "dockapplication.h"

#include <QMouseEvent>
#include <QTouchEvent>
#include <QProcess>

#include <DLog>

int main(int argc, char **argv)
{
    qputenv("QT_QPA_PLATFORM", "offscreen");

    DockApplication app(argc, argv);

    qApp->setProperty("CANSHOW", true);

    ::testing::InitGoogleTest(&argc, argv);

#ifdef QT_DEBUG
    __sanitizer_set_report_path("asan.log");
#endif

    int ret = RUN_ALL_TESTS();

    // 重启一下任务栏,算是默认恢复一下吧
    QProcess::startDetached("bash -c \"pkill dde-dock; /usr/bin/dde-dock\"");

    return ret;
}
