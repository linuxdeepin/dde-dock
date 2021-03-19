/*
 * Copyright (C) 2018 ~ 2028 Uniontech Technology Co., Ltd.
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
#include <QDebug>

#include <DWindowManagerHelper>

#include <gtest/gtest.h>

#define private public
#include "dockitemmanager.h"
#include "dockitem.h"
#undef private
#include "item/testplugin.h"

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

//    manager->startLoadPlugins();
}

TEST_F(Test_DockItemManager, get_method_test)
{
    manager->itemList();
    manager->pluginList();

    qDebug() << manager->m_itemList.size();
    for (auto item: manager->m_itemList)
        qDebug() << item->itemType();
}

TEST_F(Test_DockItemManager, refershItemsIcon_test)
{
    manager->refershItemsIcon();
    manager->sortPluginItems();
}

TEST_F(Test_DockItemManager, cover_test)
{
    manager->updatePluginsItemOrderKey();
    manager->itemAdded("", 0);
    manager->appItemAdded(QDBusObjectPath(), 0);
    manager->onPluginLoadFinished();
    manager->reloadAppItems();
    manager->pluginItemRemoved(new PluginsItem(new TestPlugin, "", ""));
    manager->pluginItemInserted(new PluginsItem(new TestPlugin, "", ""));
    manager->appItemRemoved("");
    manager->itemMoved(new PluginsItem(new TestPlugin, "", ""), new PluginsItem(new TestPlugin, "", ""));
}
