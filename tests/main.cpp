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
#include <gtest/gtest.h>

#include "dockapplication.h"

#include <QMouseEvent>
#include <QTouchEvent>

#include <DLog>

int main(int argc, char **argv)
{
    // gerrit编译时没有显示器，需要指定环境变量,本地Debug模式编译时不要设置这个宏，导致获取不到显示器相关信息
#ifndef QT_DEBUG
    qputenv("QT_QPA_PLATFORM", "offscreen");
#endif

    DockApplication app(argc, argv);

    qApp->setProperty("CANSHOW", true);

    ::testing::InitGoogleTest(&argc, argv);

    return RUN_ALL_TESTS();
}
