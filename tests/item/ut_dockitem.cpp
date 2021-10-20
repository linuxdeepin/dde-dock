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
#include <gmock/gmock.h>

#define protected public
#define private public
#include "dockitem.h"
#undef private
#undef protected

class Test_DockItem : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    DockItem *dockItem = nullptr;
};

void Test_DockItem::SetUp()
{
    dockItem = new DockItem();
}

void Test_DockItem::TearDown()
{
    delete dockItem;
    dockItem = nullptr;
}

TEST_F(Test_DockItem, dockitem_test)
{
    ASSERT_NE(dockItem, nullptr);
}

TEST_F(Test_DockItem, show_test)
{
    dockItem->show();

    QThread::msleep(450);

    ASSERT_EQ(dockItem->isVisible(), true);
}

TEST_F(Test_DockItem, hide_test)
{
    dockItem->hide();

    QThread::msleep(450);

    ASSERT_EQ(dockItem->isVisible(), false);
}

TEST_F(Test_DockItem, cover_test)
{
    DockItem::setDockPosition(Dock::Top);
    DockItem::setDockDisplayMode(Dock::Fashion);

    ASSERT_EQ(dockItem->itemType(), DockItem::App);
    dockItem->sizeHint();

    dockItem->refreshIcon();
    dockItem->contextMenu();
    dockItem->popupTips();
    dockItem->popupWindowAccept();
    dockItem->showPopupApplet(new QWidget);
    dockItem->invokedMenuItem("", true);
    dockItem->checkAndResetTapHoldGestureState();

    ASSERT_EQ(dockItem->accessibleName(), "");
}

TEST_F(Test_DockItem, event_test)
{
    DockItem *item = new DockItem;
    item->m_popupShown = true;
    item->update();

    QMouseEvent event(QEvent::MouseButtonPress, QPointF(0.0, 0.0), Qt::NoButton, Qt::NoButton, Qt::NoModifier);
    qApp->sendEvent(item, &event);

    QEnterEvent event1(QPointF(0.0, 0.0), QPointF(0.0, 0.0), QPointF(0.0, 0.0));
    qApp->sendEvent(item, &event1);

    QEvent event2(QEvent::Leave);
    qApp->sendEvent(item, &event2);

    QTest::qWait(10);

    QEvent e(QEvent::Enter);
    item->enterEvent(&e);

    QAction *action = new QAction();
    item->menuActionClicked(action);

    item->onContextMenuAccepted();

    item->showHoverTips();

    QEvent *deleteEvent = new QEvent(QEvent::DeferredDelete);
    qApp->postEvent(item, deleteEvent);
    deleteEvent = new QEvent(QEvent::DeferredDelete);
    qApp->postEvent(action, deleteEvent);

    item->showContextMenu();
}

TEST_F(Test_DockItem, topleftPoint_test)
{
    DockItem::setDockPosition(Dock::Top);
    dockItem->popupMarkPoint();
    dockItem->topleftPoint();
    DockItem::setDockPosition(Dock::Right);
    dockItem->popupMarkPoint();
    dockItem->topleftPoint();
    DockItem::setDockPosition(Dock::Bottom);
    dockItem->popupMarkPoint();
    dockItem->topleftPoint();
    DockItem::setDockPosition(Dock::Left);
    dockItem->popupMarkPoint();
    dockItem->topleftPoint();

    ASSERT_TRUE(true);
}
