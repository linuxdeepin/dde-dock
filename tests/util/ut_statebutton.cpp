/*
 * Copyright (C) 2018 ~ 2028 Uniontech Technology Co., Ltd.
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

#include <gtest/gtest.h>

#include "statebutton.h"

class Test_StateButton : public QObject, public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    StateButton *stateButton = nullptr;
};

void Test_StateButton::SetUp()
{
    stateButton = new StateButton();
}

void Test_StateButton::TearDown()
{
    delete stateButton;
    stateButton = nullptr;
}

TEST_F(Test_StateButton, statebutton_test)
{
    ASSERT_NE(stateButton, nullptr);
}

TEST_F(Test_StateButton, statebutton_clicked_test)
{
    bool clicked = false;

    connect(stateButton, &StateButton::click, this, [ = ]() mutable {
        clicked = true;
    });

    Qt::MouseButton button = Qt::LeftButton;
    QMouseEvent mouseEvent(QEvent::MouseButtonPress, stateButton->rect().center(), button, Qt::NoButton, Qt::NoModifier);
    bool ret = QApplication::sendEvent(stateButton, &mouseEvent);

    ASSERT_NE(ret, clicked);
}
