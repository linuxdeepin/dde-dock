// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include <QObject>
#include <QDebug>
#include <QTest>

#include <DWindowManagerHelper>

#include <gtest/gtest.h>

#define private public
#include "dockitemmanager.h"
#include "dockitem.h"
#undef private
#include "item/testplugin.h"
#include "traypluginitem.h"

class Test_DockItemManager : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    DockItemManager *manager = nullptr;
};

void Test_DockItemManager::SetUp()
{
    manager = DockItemManager::instance();
}

void Test_DockItemManager::TearDown()
{
}

TEST_F(Test_DockItemManager, appIsOnDock_test)
{
    manager->appIsOnDock("test");

    //TODO 问题从这里开始产生
//    manager->startLoadPlugins();
//    QTest::qWait(10);
}

TEST_F(Test_DockItemManager, get_method_test)
{
    manager->itemList();
    manager->pluginList();

    for (auto item: manager->m_itemList)
        qDebug() << item->itemType();
}

TEST_F(Test_DockItemManager, refreshItemsIcon_test)
{
    manager->refreshItemsIcon();
}

TEST_F(Test_DockItemManager, coverage_test)
{
    manager->updatePluginsItemOrderKey();
    manager->itemAdded("", 0);
    manager->appItemAdded(QDBusObjectPath(), 0);
    manager->onPluginLoadFinished();
    manager->reloadAppItems();

    QScopedPointer<TestPlugin> testPlugin(new TestPlugin);
    TrayPluginItem item(testPlugin.get(), "", "");
    manager->pluginItemInserted(&item);

    manager->itemMoved(manager->itemList().first(), &item);

    manager->pluginItemRemoved(&item);

    manager->appItemRemoved("");
}
