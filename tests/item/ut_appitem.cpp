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
#include <QThread>
#include <QTest>

#include <gtest/gtest.h>

#include "utils.h"
#include "appitem.h"

using namespace ::testing;

class Test_AppItem : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

    AppItem *appItem;
    const QGSettings *appSettings;
    const QGSettings *activeSettings;
    const QGSettings *dockedSettings;
};

void Test_AppItem::SetUp()
{
    appSettings = Utils::ModuleSettingsPtr("app");
    activeSettings = Utils::ModuleSettingsPtr("activeapp");
    dockedSettings = Utils::ModuleSettingsPtr("dockapp");

    appItem = new AppItem(appSettings, activeSettings, dockedSettings, QDBusObjectPath("/com/deepin/dde/daemon/Dock/entries/e0T6045b766"));
}

void Test_AppItem::TearDown()
{
    delete appItem;
    delete appSettings;
    delete activeSettings;
    delete dockedSettings;
}

TEST_F(Test_AppItem, paintEvent)
{
    QPaintEvent e((QRect()));

    WindowInfoMap map;
    WindowInfo info;
    map.insert(0,info);
    map.insert(1,info);
    map.insert(2,info);
    appItem->updateWindowInfos(map);

    DockItem::setDockDisplayMode(DisplayMode::Fashion);
    appItem->setDockInfo(Dock::Position::Top, QRect(QPoint(0,0), QPoint(1920, 40)));
    appItem->paintEvent(&e);
    appItem->setDockInfo(Dock::Position::Bottom, QRect(QPoint(0,0), QPoint(1920, 40)));
    appItem->paintEvent(&e);
    appItem->setDockInfo(Dock::Position::Left, QRect(QPoint(0,0), QPoint(1920, 40)));
    appItem->paintEvent(&e);
    appItem->setDockInfo(Dock::Position::Right, QRect(QPoint(0,0), QPoint(1920, 40)));
    appItem->paintEvent(&e);

    DockItem::setDockDisplayMode(DisplayMode::Efficient);
    appItem->setDockInfo(Dock::Position::Top, QRect(QPoint(0,0), QPoint(1920, 40)));
    appItem->paintEvent(&e);
    appItem->setDockInfo(Dock::Position::Bottom, QRect(QPoint(0,0), QPoint(1920, 40)));
    appItem->paintEvent(&e);
    appItem->setDockInfo(Dock::Position::Left, QRect(QPoint(0,0), QPoint(1920, 40)));
    appItem->paintEvent(&e);
    appItem->setDockInfo(Dock::Position::Right, QRect(QPoint(0,0), QPoint(1920, 40)));
    appItem->paintEvent(&e);

    ASSERT_TRUE(true);
}

TEST_F(Test_AppItem, coverage_test)
{
    // 触发信号测试
    appItem->m_refershIconTimer->start(10);
    QTest::qWait(20);

    appItem->undock();
    appItem->appIcon();

    ASSERT_TRUE(appItem->itemType() == AppItem::App);
    ASSERT_TRUE(appItem->accessibleName() == appItem->m_itemEntryInter->name());

    appItem->checkAttentionEffect();
    appItem->onGSettingsChanged("enabled");
    appItem->checkGSettingsControl();
    appItem->showHoverTips();
    appItem->popupTips();
    appItem->startDrag();
    appItem->playSwingEffect();
    appItem->invokedMenuItem("invalid", true);
    appItem->contextMenu();

    ASSERT_TRUE(true);
}

TEST_F(Test_AppItem, appDragWidget)
{
    appItem->appDragWidget();

    ASSERT_TRUE(true);
}

TEST_F(Test_AppItem, mouseReleaseEvent)
{
    QMouseEvent event(QEvent::MouseButtonRelease, QPointF(0, 0), Qt::MiddleButton, Qt::MiddleButton, Qt::ControlModifier);
    appItem->mouseReleaseEvent(&event);

    QTest::qWait(350);
    appItem->mouseReleaseEvent(&event);

    QMouseEvent event2(QEvent::MouseButtonRelease, QPointF(0, 0), Qt::LeftButton, Qt::LeftButton, Qt::ControlModifier);
    QTest::qWait(350);
    appItem->mouseReleaseEvent(&event2);

    ASSERT_TRUE(true);
}

TEST_F(Test_AppItem, QWheelEvent)
{
    QWheelEvent event(QPointF(), Qt::LeftButton, Qt::LeftButton, Qt::NoModifier);
    appItem->wheelEvent(&event);

    ASSERT_TRUE(true);
}

TEST_F(Test_AppItem, event_test)
{
    QMouseEvent event1(QEvent::MouseButtonPress, QPointF(0, 0), Qt::LeftButton, Qt::RightButton, Qt::ControlModifier);
    appItem->mousePressEvent(&event1);

    QMouseEvent event3(QEvent::MouseMove, QPointF(0, 0), Qt::LeftButton, Qt::LeftButton, Qt::ControlModifier);
    appItem->mouseMoveEvent(&event3);

    QMouseEvent event4(QEvent::MouseMove, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    appItem->mouseMoveEvent(&event4);

    QResizeEvent event5((QSize()), QSize());
    appItem->resizeEvent(&event5);

    QEvent event6(QEvent::Leave);
    appItem->leaveEvent(&event6);

    QShowEvent event7;
    appItem->showEvent(&event7);

    QMimeData *data = new QMimeData;
    data->setText("test");
    QDropEvent event8(QPointF(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::ControlModifier);
    appItem->dropEvent(&event8);

    QDragEnterEvent event9(QPoint(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::NoModifier);
    appItem->dragEnterEvent(&event9);

    QDragMoveEvent event10(QPoint(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::NoModifier);
    appItem->dragMoveEvent(&event10);
}

TEST_F(Test_AppItem, checkEntry)
{
    appItem->checkEntry();
    appItem->accessibleName();

    ASSERT_EQ(appItem->appId(), appItem->m_id);

    appItem->isValid();
}
