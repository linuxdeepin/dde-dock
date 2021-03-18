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

#include <QDebug>
#include <QMouseEvent>
#include <QTouchEvent>

#include <DLog>

int main(int argc, char **argv)
{
    // gerrit编译时没有显示器，需要指定环境变量
    qputenv("QT_QPA_PLATFORM", "offscreen");
    DockApplication app(argc, argv);

    qApp->setProperty("CANSHOW", true);

    QMouseEvent mouseEvent1(QMouseEvent::MouseButtonPress, QPoint(), QPoint(), QPoint(), Qt::LeftButton, Qt::LeftButton, Qt::NoModifier, Qt::MouseEventSynthesizedByQt);
    qApp->sendEvent(&app, &mouseEvent1);

    QMouseEvent mouseEvent2(QMouseEvent::MouseButtonPress, QPoint(), QPoint(), QPoint(), Qt::LeftButton, Qt::LeftButton, Qt::NoModifier, Qt::MouseEventSynthesizedByApplication);
    qApp->sendEvent(&app, &mouseEvent2);

    QList<QTouchEvent::TouchPoint> list;
    list << QTouchEvent::TouchPoint(0) << QTouchEvent::TouchPoint(1) << QTouchEvent::TouchPoint(2);
    QTouchEvent touchEven(QEvent::TouchBegin, nullptr, Qt::NoModifier, Qt::TouchPointPressed, list);
    qApp->sendEvent(&app, &touchEven);

    ::testing::InitGoogleTest(&argc, argv);

    return RUN_ALL_TESTS();
}
