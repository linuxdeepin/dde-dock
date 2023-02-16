// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include <gtest/gtest.h>
#ifdef QT_DEBUG
#include <sanitizer/asan_interface.h>
#endif

#include "dockapplication.h"

#include <QMouseEvent>
#include <QTouchEvent>

#include <DLog>

int main(int argc, char **argv)
{
    qputenv("QT_QPA_PLATFORM", "offscreen");

    DockApplication app(argc, argv);
    // 设置应用名为dde-dock，否则dconfig相关的配置就读不到了
    app.setApplicationName("dde-dock");

    qApp->setProperty("CANSHOW", true);

    ::testing::InitGoogleTest(&argc, argv);

#ifdef QT_DEBUG
    __sanitizer_set_report_path("asan.log");
#endif

    return RUN_ALL_TESTS();
}
