// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include <QObject>
#include <QDebug>
#include <QTest>

#include <DWindowManagerHelper>

#include <gtest/gtest.h>

#include "dockpluginscontroller.h"
#include "abstractpluginscontroller.h"
#include "../../plugins/bluetooth/bluetoothplugin.h"

class Test_DockPluginsController : public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    DockPluginsController *controller = nullptr;
};

void Test_DockPluginsController::SetUp()
{
    controller = new DockPluginsController();
}

void Test_DockPluginsController::TearDown()
{
    delete controller;
}

TEST_F(Test_DockPluginsController, test)
{
    controller->loadPlugin("/usr/lib/dde-dock/plugins/libtray.so");
//    BluetoothPlugin * const p = new BluetoothPlugin;
//    p->init(controller);
}
