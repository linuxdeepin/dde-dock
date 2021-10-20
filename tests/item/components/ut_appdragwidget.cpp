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
#include <QTest>

#include <gtest/gtest.h>

#include "appdragwidget.h"

class Test_AppDragWidget : public ::testing::Test
{};

TEST_F(Test_AppDragWidget, popupMarkPoint)
{
    AppDragWidget dragWidget;
    QPixmap pix(":/res/all_settings_on.png");
    dragWidget.setAppPixmap(pix);
    dragWidget.setOriginPos(QPoint(-1, -1));

    dragWidget.popupMarkPoint(Dock::Position::Top);
    dragWidget.popupMarkPoint(Dock::Position::Bottom);
    dragWidget.popupMarkPoint(Dock::Position::Left);
    dragWidget.popupMarkPoint(Dock::Position::Right);

    dragWidget.showRemoveTips();
    dragWidget.showGoBackAnimation();

    ASSERT_TRUE(true);
}

TEST_F(Test_AppDragWidget, isRemoveAble)
{
    AppDragWidget dragWidget;

    //    dragWidget.show();
    //    dragWidget.hide();

    QTest::mouseClick(&dragWidget,Qt::LeftButton, Qt::NoModifier, QPoint(dragWidget.rect().center()));

    // bottom
    const QRect &rect = QRect(QPoint(0, 1040), QPoint(1920, 1080));
    dragWidget.setDockInfo(Dock::Position::Bottom, rect);
    dragWidget.isRemoveItem();
    ASSERT_TRUE(dragWidget.isRemoveAble(QPoint(10, 10)));
    ASSERT_FALSE(dragWidget.isRemoveAble(QPoint(10, 1070)));
    ASSERT_TRUE(dragWidget.isRemoveAble(QPoint(1910, 10)));
    ASSERT_FALSE(dragWidget.isRemoveAble(QPoint(1910, 1070)));

    // top
    const QRect &rect1 = QRect(QPoint(0, 0), QPoint(1920, 40));
    dragWidget.setDockInfo(Dock::Position::Top, rect1);
    dragWidget.isRemoveItem();
    ASSERT_FALSE(dragWidget.isRemoveAble(QPoint(10, 10)));
    ASSERT_TRUE(dragWidget.isRemoveAble(QPoint(10, 1070)));
    ASSERT_FALSE(dragWidget.isRemoveAble(QPoint(1910, 10)));
    ASSERT_TRUE(dragWidget.isRemoveAble(QPoint(1910, 1070)));

    // left
    const QRect &rect2 = QRect(QPoint(0, 0), QPoint(40, 1080));
    dragWidget.setDockInfo(Dock::Position::Left, rect2);
    dragWidget.isRemoveItem();
    ASSERT_FALSE(dragWidget.isRemoveAble(QPoint(10, 10)));
    ASSERT_FALSE(dragWidget.isRemoveAble(QPoint(10, 1070)));
    ASSERT_TRUE(dragWidget.isRemoveAble(QPoint(1910, 10)));
    ASSERT_TRUE(dragWidget.isRemoveAble(QPoint(1910, 1070)));

    // right
    const QRect &rect3 = QRect(QPoint(1880, 0), QPoint(1920, 1080));
    dragWidget.setDockInfo(Dock::Position::Right, rect3);
    dragWidget.isRemoveItem();
    ASSERT_TRUE(dragWidget.isRemoveAble(QPoint(10, 10)));
    ASSERT_TRUE(dragWidget.isRemoveAble(QPoint(10, 1070)));
    ASSERT_FALSE(dragWidget.isRemoveAble(QPoint(1910, 10)));
    ASSERT_FALSE(dragWidget.isRemoveAble(QPoint(1910, 1070)));
}

TEST_F(Test_AppDragWidget, coverage_test)
{
    AppDragWidget dragWidget;

    dragWidget.showRemoveAnimation();
    dragWidget.onRemoveAnimationStateChanged(QAbstractAnimation::State::Stopped, QAbstractAnimation::State::Running);
}

TEST_F(Test_AppDragWidget, event_test)
{
    AppDragWidget dragWidget;

    QMouseEvent mouseMoveEvent_(QEvent::MouseMove, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    dragWidget.mouseMoveEvent(&mouseMoveEvent_);

    QMimeData *data = new QMimeData;
    data->setText("test");
    QDropEvent dropEvent_(QPointF(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::ControlModifier);
    dragWidget.dropEvent(&dropEvent_);

    QDragEnterEvent dragEnterEvent_(QPoint(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::NoModifier);
    dragWidget.dragEnterEvent(&dragEnterEvent_);

    QDragMoveEvent dragMoveEvent_(QPoint(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::NoModifier);
    dragWidget.dragMoveEvent(&dragMoveEvent_);

    QHideEvent hideEvent_;
    dragWidget.hideEvent(&hideEvent_);

    QEvent enterEvent_(QEvent::Enter);
    dragWidget.enterEvent(&enterEvent_);

    data->deleteLater();
    ASSERT_TRUE(true);
}
