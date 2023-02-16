// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include <QObject>
#include <QApplication>
#include <QSignalSpy>
#include <QTest>

#include <gtest/gtest.h>

#include "pluginloader.h"

class Test_PluginLoader : public QObject, public ::testing::Test
{
public:
    virtual void SetUp() override;
    virtual void TearDown() override;

public:
    PluginLoader *loader = nullptr;
};

void Test_PluginLoader::SetUp()
{
    loader = new PluginLoader("../", nullptr);
    connect(loader, &PluginLoader::finished, loader, &PluginLoader::deleteLater, Qt::QueuedConnection);
}

void Test_PluginLoader::TearDown()
{
}

TEST_F(Test_PluginLoader, loader_test)
{
    loader->start();
}
