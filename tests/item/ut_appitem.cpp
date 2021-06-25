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
#define private public
#include "appitem.h"
#undef private

using namespace ::testing;

class Test_AppItem : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;
};

void Test_AppItem::SetUp()
{
}

void Test_AppItem::TearDown()
{
}

TEST_F(Test_AppItem, coverage_test)
{
    const QGSettings *appSettings = Utils::ModuleSettingsPtr("app");
    const QGSettings *activeSettings = Utils::ModuleSettingsPtr("activeapp");
    const QGSettings *dockedSettings = Utils::ModuleSettingsPtr("dockapp");

    AppItem *appItem = new AppItem(appSettings, activeSettings, dockedSettings, QDBusObjectPath("/com/deepin/dde/daemon/Dock/entries/e0T6045b766"));

    // 触发信号测试
    //    emit appItem->m_refershIconTimer->start(10);
    QTest::qWait(20);

    // FIXME: 测试不到？
    appItem->checkEntry();
    appItem->undock();

    appItem->setDockDisplayMode(Dock::Efficient);
    appItem->update();
    QTest::qWait(10);
    appItem->setDockDisplayMode(Dock::Fashion);
    appItem->update();
    QTest::qWait(10);
    //    appItem->updateWindowIconGeometries();

    ASSERT_TRUE(appItem->itemType() == AppItem::App);

    appItem->setDockInfo(Dock::Position::Top, QRect(QPoint(0,0), QPoint(1920, 40)));

    ASSERT_TRUE(appItem->accessibleName() == appItem->m_itemEntryInter->name());

    appItem->show();

    appItem->resize(100, 100);

    QTest::qWait(10);

    ASSERT_TRUE(appItem->isVisible());

    appItem->hide();

    QTest::qWait(10);

    ASSERT_TRUE(!appItem->isVisible());

    QTest::mousePress(appItem, Qt::LeftButton, Qt::NoModifier);
    QTest::mouseRelease(appItem, Qt::LeftButton, Qt::NoModifier);
    QTest::qWait(400);
    //    QTest::mouseClick(appItem, Qt::MiddleButton, Qt::NoModifier);
    //    QTest::qWait(400);
    //    QTest::mouseClick(appItem, Qt::LeftButton, Qt::NoModifier, QPoint(-1, -1));
    //    QTest::qWait(400);
    //    QTest::mouseMove(appItem, appItem->geometry().center());

    delete appItem;
    appItem = nullptr;
    delete appSettings;
    appSettings = nullptr;
    delete activeSettings;
    activeSettings = nullptr;
    delete dockedSettings;
    dockedSettings = nullptr;
}
