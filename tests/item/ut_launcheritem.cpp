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

#include <gtest/gtest.h>
#include <gmock/gmock.h>

using namespace ::testing;

#define private public
#include "launcheritem.h"
#undef private

class Test_LauncherItem : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Test_LauncherItem::SetUp()
{
}

void Test_LauncherItem::TearDown()
{
}

TEST_F(Test_LauncherItem, event_test)
{
    LauncherItem *launcherItem = new LauncherItem;
    QMouseEvent event1(QEvent::MouseButtonPress, QPointF(0, 0), Qt::LeftButton, Qt::RightButton, Qt::ControlModifier);
    launcherItem->mousePressEvent(&event1);

    QMouseEvent event2(QEvent::MouseButtonRelease, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    launcherItem->mouseReleaseEvent(&event2);

    QMouseEvent event3(QEvent::MouseMove, QPointF(0, 0), Qt::LeftButton, Qt::LeftButton, Qt::ControlModifier);
    launcherItem->mouseMoveEvent(&event3);

    QMouseEvent event4(QEvent::MouseMove, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    launcherItem->mouseMoveEvent(&event4);

    QResizeEvent event5((QSize()), QSize());
    launcherItem->resizeEvent(&event5);

    QEvent event6(QEvent::Leave);
    launcherItem->leaveEvent(&event6);

    QShowEvent event7;
    launcherItem->showEvent(&event7);

    QMimeData *data = new QMimeData;
    data->setText("test");
    QDropEvent event8(QPointF(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::ControlModifier);
    launcherItem->dropEvent(&event8);

    QDragEnterEvent event9(QPoint(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::NoModifier);
    launcherItem->dragEnterEvent(&event9);

    QDragMoveEvent event10(QPoint(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::NoModifier);
    launcherItem->dragMoveEvent(&event10);

    data->deleteLater();
    delete launcherItem;
}

TEST_F(Test_LauncherItem, coverage_test)
{
    LauncherItem *launcherItem = new LauncherItem;
    ASSERT_EQ(launcherItem->itemType(), LauncherItem::Launcher);
    launcherItem->refreshIcon();
    //    launcherItem->show();
    //    QThread::msleep(10);

    //    launcherItem->hide();
    //    QThread::msleep(10);

    launcherItem->resize(100,100);
    launcherItem->popupTips();

    launcherItem->onGSettingsChanged("invalid");
    launcherItem->onGSettingsChanged("enable");

    delete launcherItem;
    launcherItem = nullptr;
}

TEST_F(Test_LauncherItem, paintEvent)
{
    LauncherItem item;
    item.setVisible(true);
    item.show();

    QRect rect;
    QPaintEvent e(rect);
    item.paintEvent(&e);
}
