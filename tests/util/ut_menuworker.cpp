// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include <QObject>
#include <QTest>
#include <QMenu>

#include <gtest/gtest.h>
#define private public
#include "menuworker.h"
#include "dockitemmanager.h"
#undef private

class Test_MenuWorker : public ::testing::Test
{};

TEST_F(Test_MenuWorker, coverage_test)
{
    MenuWorker *worker = new MenuWorker(new DBusDock("org.deepin.dde.daemon.Dock1", "/org/deepin/dde/daemon/Dock1", QDBusConnection::sessionBus()));
    DockItemManager::instance()->m_pluginsInter->m_pluginsMap.clear();
    QMenu *menu = new QMenu();
    menu = worker->createMenu(menu);
    ASSERT_FALSE(menu->isEmpty());

//    worker->showDockSettingsMenu();

    delete menu;
    delete worker;
}

TEST_F(Test_MenuWorker, setAutoHide)
{
   MenuWorker *worker = new MenuWorker(new DBusDock("org.deepin.dde.daemon.Dock1", "/org/deepin/dde/daemon/Dock1", QDBusConnection::sessionBus()));

   ASSERT_TRUE(worker->m_autoHide);
   worker->setAutoHide(false);

   ASSERT_FALSE(worker->m_autoHide);
   worker->setAutoHide(true);
   ASSERT_TRUE(worker->m_autoHide);

   delete worker;
}
