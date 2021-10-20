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

#include <QObject>
#include <QThread>
#include <QTest>
#include <QSignalSpy>
#include <QThread>

#include <gtest/gtest.h>

#define private public
#include "appsnapshot.h"
#undef private

class Test_AppSnapshot : public ::testing::Test
{};

TEST_F(Test_AppSnapshot, eventFilter)
{
    AppSnapshot snapShot(1000000);

    QEvent hoverEnterEvent(QEvent::HoverEnter);
    snapShot.eventFilter(snapShot.m_closeBtn2D, &hoverEnterEvent);

    QEvent hoverMoveEvent(QEvent::HoverMove);
    snapShot.eventFilter(snapShot.m_closeBtn2D, &hoverMoveEvent);

    QEvent hoverLeaveEvent(QEvent::HoverLeave);
    snapShot.eventFilter(snapShot.m_closeBtn2D, &hoverLeaveEvent);

    QEvent mousePressEvent(QEvent::MouseButtonPress);
    snapShot.eventFilter(snapShot.m_closeBtn2D, &mousePressEvent);
}

TEST_F(Test_AppSnapshot, paintEvent)
{
    AppSnapshot snapShot(1000000);
    QRect rect(0, 0, 10, 10);
    QPaintEvent paintEvent(rect);
    snapShot.paintEvent(&paintEvent);
}

TEST_F(Test_AppSnapshot, enterEvent)
{
    AppSnapshot snapShot(1000000);
    QEvent enterEvent(QEvent::Enter);
    snapShot.enterEvent(&enterEvent);

    ASSERT_TRUE(true);
}

TEST_F(Test_AppSnapshot, event_test)
{
    AppSnapshot snapShot(1000000);

    QMouseEvent event1(QEvent::MouseButtonPress, QPointF(0, 0), Qt::LeftButton, Qt::RightButton, Qt::ControlModifier);
    snapShot.mousePressEvent(&event1);

    QMouseEvent event2(QEvent::MouseButtonRelease, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    snapShot.mouseReleaseEvent(&event2);

    QMouseEvent event3(QEvent::MouseMove, QPointF(0, 0), Qt::LeftButton, Qt::LeftButton, Qt::ControlModifier);
    snapShot.mouseMoveEvent(&event3);

    QMouseEvent event4(QEvent::MouseMove, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    snapShot.mouseMoveEvent(&event4);

    QResizeEvent event5((QSize()), QSize());
    snapShot.resizeEvent(&event5);

    QEvent event6(QEvent::Leave);
    snapShot.leaveEvent(&event6);

    QShowEvent event7;
    snapShot.showEvent(&event7);

    QMimeData *data = new QMimeData;
    data->setText("test");
    QDropEvent event8(QPointF(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::ControlModifier);
    snapShot.dropEvent(&event8);

    QDragEnterEvent event9(QPoint(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::NoModifier);
    snapShot.dragEnterEvent(&event9);

    QDragMoveEvent event10(QPoint(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::NoModifier);
    snapShot.dragMoveEvent(&event10);

    data->deleteLater();
}

TEST_F(Test_AppSnapshot, setWindowState)
{
    AppSnapshot snapShot(1000000);

    snapShot.m_isWidowHidden = true;
    snapShot.setWindowState();

    snapShot.m_isWidowHidden = false;
    snapShot.setWindowState();

    ASSERT_TRUE(true);
}

TEST_F(Test_AppSnapshot, coverage_test)
{
    AppSnapshot snapShot(1000000);
    snapShot.closeWindow();

    QImage img;
    snapShot.rectRemovedShadow(img, nullptr);
}
