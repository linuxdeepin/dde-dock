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
#include <QApplication>
#include <QMouseEvent>
#include <QDebug>
#include <QSignalSpy>
#include <QTest>

#include <gtest/gtest.h>

#include "statebutton.h"

class Test_StateButton : public QObject, public ::testing::Test
{};

TEST_F(Test_StateButton, statebutton_clicked_test)
{
    StateButton button;
    QSignalSpy spy(&button, SIGNAL(click()));
    QTest::mousePress(&button, Qt::LeftButton, Qt::NoModifier);
    ASSERT_EQ(spy.count(), 1);
}

TEST_F(Test_StateButton, event_test)
{
    StateButton button;

    QEvent event(QEvent::Enter);
    button.enterEvent(&event);

    QEvent event2(QEvent::Leave);
    button.leaveEvent(&event2);

    ASSERT_TRUE(true);
}

TEST_F(Test_StateButton, paintEvent)
{
    StateButton button;

    QRect rect(0, 0, 10, 10);
    QPaintEvent e(rect);
    button.setType(StateButton::Check);
    button.paintEvent(&e);

    button.setType(StateButton::Fork);
    button.paintEvent(&e);

    QTest::qWait(10);
    button.setType(StateButton::Fork);
    ASSERT_EQ(button.m_type, StateButton::Fork);

    QTest::qWait(10);
    button.setType(StateButton::Check);
    ASSERT_EQ(button.m_type, StateButton::Check);
}
