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
