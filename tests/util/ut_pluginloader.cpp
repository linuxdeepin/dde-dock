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
