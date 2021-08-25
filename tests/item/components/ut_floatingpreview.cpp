/*
 * Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
 *
 * Author:     chenjun <chenjun@uniontech.com>
 *
 * Maintainer: chenjun <chenjun@uniontech.com>
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
#include <QSignalSpy>
#include <QTest>

#include <gtest/gtest.h>

#define private public
#include "appsnapshot.h"
#include "floatingpreview.h"
#undef private

class Test_FloatingPreview : public ::testing::Test
{};

TEST_F(Test_FloatingPreview, eventFilter)
{
    FloatingPreview view;

    QEvent hoverEnterEvent(QEvent::HoverEnter);
    view.eventFilter(view.m_closeBtn3D, &hoverEnterEvent);

    QEvent hoverLeaveEvent(QEvent::HoverLeave);
    view.eventFilter(view.m_closeBtn3D, &hoverLeaveEvent);

    QEvent mousePressEvent(QEvent::MouseButtonPress);
    view.eventFilter(view.m_closeBtn3D, &mousePressEvent);
}

TEST_F(Test_FloatingPreview, trackedWid)
{
    FloatingPreview view;
    AppSnapshot snap(1000000);

    view.trackWindow(&snap);
    view.onCloseBtnClicked();

    ASSERT_TRUE(view.trackedWid());
}

TEST_F(Test_FloatingPreview, paintEvent)
{
    FloatingPreview view;
    QPaintEvent event((QRect()));
    view.paintEvent(&event);

    ASSERT_TRUE(true);
}

TEST_F(Test_FloatingPreview, hideEvent)
{
    FloatingPreview view;

    AppSnapshot snap(1000000);
    view.trackWindow(&snap);

    QHideEvent event;
    view.hideEvent(&event);

    ASSERT_TRUE(true);
}

TEST_F(Test_FloatingPreview, coverage_test)
{
    QWidget parent;
    FloatingPreview view(&parent);
    AppSnapshot *shot = new AppSnapshot(1000);
    shot->fetchSnapshot();
    shot->m_snapshot = QImage(":/res/dde-calendar.svg");
    view.trackWindow(shot);

    ASSERT_TRUE(view.m_titleBtn->text() == shot->title());
    ASSERT_EQ(view.trackedWindow(), shot);

    QSignalSpy spy(shot, &AppSnapshot::clicked);
    QTest::mouseClick(&view, Qt::LeftButton, Qt::NoModifier);
    ASSERT_EQ(spy.count(), 1);

    ASSERT_TRUE(shot->contentsMargins() == QMargins(0, 0, 0, 0));

    view.trackWindow(nullptr);
    ASSERT_TRUE(view.m_titleBtn->text().isEmpty());
    ASSERT_EQ(view.trackedWindow(), shot);
}
