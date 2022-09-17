// SPDX-FileCopyrightText: 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "module_widget.h"

#include <QWidget>

#include <gtest/gtest.h>

class Test_ModuleWidget : public QObject, public ::testing::Test
{};

TEST_F(Test_ModuleWidget, updateSliderValue)
{
    ModuleWidget widget;

    widget.updateSliderValue();
}
