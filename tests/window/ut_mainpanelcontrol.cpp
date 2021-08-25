/*
 * Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
 *
 * Author:     weizhixiang <weizhixiang@uniontech.com>
 *
 * Maintainer: weizhixiang <weizhixiang@uniontech.com>
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
#include <QGraphicsView>

#include <gtest/gtest.h>

#define private public
#include "mainpanelcontrol.h"
#include "appitem.h"
#include "dockitem.h"
#include "placeholderitem.h"
#include "pluginsitem.h"
#include "traypluginitem.h"
#include "launcheritem.h"
#undef private

#include "../item/testplugin.h"

using namespace ::testing;

class Test_MainPanelControl : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    MainPanelControl *mainPanel;
};

void Test_MainPanelControl::SetUp()
{
    mainPanel = new MainPanelControl();
}

void Test_MainPanelControl::TearDown()
{
    delete mainPanel;
    mainPanel = nullptr;
}

TEST_F(Test_MainPanelControl, getTrayVisableItemCount)
{
    MainPanelControl panel;
    TestPlugin plugin;

    ASSERT_EQ(panel.m_trayIconCount, 0);
    panel.getTrayVisableItemCount();
    ASSERT_EQ(panel.m_trayIconCount, 0);

    TrayPluginItem trayPluginItem(&plugin, "tray", "1.2.0");
    panel.addTrayAreaItem(0, &trayPluginItem);
    panel.getTrayVisableItemCount();
}

TEST_F(Test_MainPanelControl, paintEvent)
{
    MainPanelControl panel;
    QRect paintRect(0, 0, 10, 10);
    QPaintEvent event(paintRect);

    panel.m_isHover = true;
    panel.paintEvent(&event);

    panel.m_isHover = false;
    panel.paintEvent(&event);

    ASSERT_TRUE(true);
}

TEST_F(Test_MainPanelControl, moveAppSonWidget)
{
    MainPanelControl panel;
    panel.m_dislayMode = DisplayMode::Fashion;
    panel.m_position = Position::Top;
    panel.moveAppSonWidget();

    panel.m_position = Position::Bottom;
    panel.moveAppSonWidget();

    panel.m_position = Position::Left;
    panel.moveAppSonWidget();

    panel.m_position = Position::Right;
    panel.moveAppSonWidget();

    panel.m_dislayMode = DisplayMode::Efficient;
    panel.m_position = Position::Top;
    panel.moveAppSonWidget();

    panel.m_position = Position::Bottom;
    panel.moveAppSonWidget();

    panel.m_position = Position::Left;
    panel.moveAppSonWidget();

    panel.m_position = Position::Right;
    panel.moveAppSonWidget();

    ASSERT_TRUE(true);
}

TEST_F(Test_MainPanelControl, startDrag)
{
    MainPanelControl panel;
    TestPlugin plugin;

    AppItem appItem(nullptr, nullptr, nullptr, QDBusObjectPath());
    panel.addAppAreaItem(0, &appItem);
    panel.startDrag(&appItem);

    LauncherItem launcherItem;
    mainPanel->addFixedAreaItem(0, &launcherItem);
    panel.startDrag(&launcherItem);

    PluginsItem pluginItem(&plugin, "monitor", "1.2.1");
    mainPanel->addPluginAreaItem(0, &pluginItem);
    panel.startDrag(&pluginItem);
}

TEST_F(Test_MainPanelControl, eventFilter)
{
    MainPanelControl panel;
    QResizeEvent event((QSize()), QSize());
    panel.eventFilter(mainPanel->m_appAreaSonWidget, &event);
    panel.eventFilter(mainPanel->m_appAreaWidget, &event);

    QEvent enterEvent(QEvent::Enter);
    panel.eventFilter(mainPanel->m_desktopWidget, &enterEvent);

    QEvent leaveEvent(QEvent::Leave);
    panel.eventFilter(mainPanel->m_desktopWidget, &leaveEvent);

    QEvent moveEvent(QEvent::Move);
    panel.eventFilter(mainPanel->m_appAreaWidget, &moveEvent);

    QMouseEvent mouseMoveEvent(QEvent::MouseMove, QPointF(0, 0), Qt::LeftButton, Qt::LeftButton, Qt::NoModifier);
    panel.eventFilter(mainPanel, &mouseMoveEvent);

    //    QEvent dragMoveEvent(QEvent::DragMove);
    //    mainPanel->eventFilter(static_cast<QGraphicsView *>(mainPanel->m_appDragWidget), &dragMoveEvent);
}

TEST_F(Test_MainPanelControl, moveItem)
{
    MainPanelControl panel;
    TestPlugin plugin;

    TestPlugin fixedPlugin;
    fixedPlugin.setType(PluginsItemInterface::PluginType::Fixed);

    DockItem dockItem1;
    DockItem dockItem2;
    panel.addAppAreaItem(0, &dockItem1);
    panel.addAppAreaItem(0, &dockItem2);
    panel.moveItem(&dockItem1, &dockItem2);

    LauncherItem launcherItem1;
    LauncherItem launcherItem2;
    panel.addFixedAreaItem(0, &launcherItem1);
    panel.addFixedAreaItem(0, &launcherItem2);
    panel.moveItem(&launcherItem1, &launcherItem2);

    PluginsItem pluginItem1(&plugin, "monitor", "1.2.1");
    PluginsItem pluginItem2(&plugin, "monitor", "1.2.1");
    panel.addPluginAreaItem(0, &pluginItem1);
    panel.addPluginAreaItem(0, &pluginItem2);
    panel.moveItem(&pluginItem1, &pluginItem2);

    PluginsItem fixedPluginItem1(&fixedPlugin, "monitor", "1.2.1");
    PluginsItem fixedPluginItem2(&fixedPlugin, "monitor", "1.2.1");
    panel.addPluginAreaItem(0, &fixedPluginItem1);
    panel.addPluginAreaItem(0, &fixedPluginItem2);
    panel.moveItem(&fixedPluginItem1, &fixedPluginItem1);

    // dropTargetItem test
    panel.dropTargetItem(&dockItem1, QPoint(0, 0));
    panel.dropTargetItem(&launcherItem1, QPoint(0, 0));
    panel.dropTargetItem(&pluginItem1, QPoint(0, 0));
    panel.dropTargetItem(nullptr, QPoint(-1, -1));
}

TEST_F(Test_MainPanelControl, removeItem)
{
    MainPanelControl panel;
    TestPlugin plugin;

    DockItem dockItem;
    panel.addAppAreaItem(0, &dockItem);
    panel.removeItem(&dockItem);

    PlaceholderItem placeHolderItem;
    panel.addAppAreaItem(0, &placeHolderItem);
    panel.removeItem(&placeHolderItem);

    LauncherItem launcherItem;
    panel.addFixedAreaItem(0, &launcherItem);
    panel.removeItem(&launcherItem);

    TrayPluginItem trayPluginItem(&plugin, "tray", "1.2.0");
    panel.addTrayAreaItem(0, &trayPluginItem);
    panel.removeItem(&trayPluginItem);

    PluginsItem pluginItem(&plugin, "monitor", "1.2.1");
    panel.addPluginAreaItem(0, &pluginItem);
    panel.removeItem(&pluginItem);

    ASSERT_TRUE(true);
}

TEST_F(Test_MainPanelControl, test1)
{
    MainPanelControl panel;
    DockItem dockItem;

    panel.insertItem(0, &dockItem);

    ASSERT_TRUE(true);
}

TEST_F(Test_MainPanelControl, updateMainPanelLayout)
{
    MainPanelControl panel;

    panel.setPositonValue(Dock::Position::Top);
    panel.updateMainPanelLayout();
    QTest::qWait(10);

    panel.setPositonValue(Dock::Position::Bottom);
    panel.updateMainPanelLayout();
    QTest::qWait(10);

    panel.setPositonValue(Dock::Position::Left);
    panel.updateMainPanelLayout();
    QTest::qWait(10);

    panel.setPositonValue(Dock::Position::Right);
    panel.updateMainPanelLayout();
    QTest::qWait(10);

    ASSERT_TRUE(true);
}

TEST_F(Test_MainPanelControl, event_test)
{
    MainPanelControl panel;

    QMouseEvent event1(QEvent::MouseButtonPress, QPointF(0, 0), Qt::LeftButton, Qt::RightButton, Qt::ControlModifier);
    panel.mousePressEvent(&event1);

    QMouseEvent event2(QEvent::MouseButtonRelease, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    panel.mouseReleaseEvent(&event2);

    QMouseEvent event3(QEvent::MouseMove, QPointF(0, 0), Qt::LeftButton, Qt::LeftButton, Qt::ControlModifier);
    panel.mouseMoveEvent(&event3);

    QMouseEvent event4(QEvent::MouseMove, QPointF(0, 0), Qt::RightButton, Qt::RightButton, Qt::ControlModifier);
    panel.mouseMoveEvent(&event4);

    QResizeEvent event5((QSize()), QSize());
    panel.resizeEvent(&event5);

    QEvent event6(QEvent::Leave);
    panel.leaveEvent(&event6);

    QShowEvent event7;
    panel.showEvent(&event7);

    QMimeData *data = new QMimeData;
    data->setText("test");
    QDropEvent event8(QPointF(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::ControlModifier);
    panel.dropEvent(&event8);

    QDragEnterEvent event9(QPoint(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::NoModifier);
    panel.dragEnterEvent(&event9);

    QDragMoveEvent event10(QPoint(), Qt::DropAction::CopyAction, data, Qt::LeftButton, Qt::NoModifier);
    panel.dragMoveEvent(&event10);
}

TEST_F(Test_MainPanelControl, dragLeaveEvent)
{
    MainPanelControl panel;

    QDragLeaveEvent event11;
    panel.dragLeaveEvent(&event11);

    ASSERT_TRUE(true);
}

TEST_F(Test_MainPanelControl, coverage_test)
{
    MainPanelControl panel;
    QScopedPointer<QWidget> w(new QWidget);
    panel.removeAppAreaItem(w.get());
    panel.removeTrayAreaItem(w.get());
    panel.updateAppAreaSonWidgetSize();
    panel.checkNeedShowDesktop();
    panel.appIsOnDock("123");
}

TEST_F(Test_MainPanelControl, addItem)
{
    MainPanelControl panel;

    panel.setDisplayMode(DisplayMode::Fashion);
    ASSERT_EQ(panel.m_dislayMode, DisplayMode::Fashion);

    panel.setDisplayMode(DisplayMode::Efficient);
    ASSERT_EQ(panel.m_dislayMode, DisplayMode::Efficient);

    panel.setPositonValue(Position::Top);
    QWidget *fixedWidget = new QWidget;
    QWidget *appWidget = new QWidget;
    QWidget *pluginWidget = new QWidget;

    panel.addFixedAreaItem(0, fixedWidget);
    panel.addAppAreaItem(0, appWidget);
    panel.addPluginAreaItem(0, pluginWidget);

    panel.updateAppAreaSonWidgetSize();

    panel.removeFixedAreaItem(fixedWidget);
    panel.removeAppAreaItem(appWidget);
    panel.removePluginAreaItem(pluginWidget);

    panel.setPositonValue(Position::Left);
    panel.addFixedAreaItem(0, fixedWidget);
    panel.addAppAreaItem(0, appWidget);
    panel.addPluginAreaItem(0, pluginWidget);
    panel.updateAppAreaSonWidgetSize();

    DockItem *dockItem1 = new DockItem;
    DockItem *dockItem2 = new DockItem;
    panel.addAppAreaItem(0, dockItem1);
    panel.addAppAreaItem(0, dockItem2);
    //    panel.moveItem(dockItem1, dockItem2);

    panel.itemUpdated(dockItem2);

    delete fixedWidget;
    delete appWidget;
    delete pluginWidget;

    ASSERT_TRUE(true);
}
