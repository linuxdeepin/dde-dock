#include "module_widget.h"

#include <QWidget>

#include <gtest/gtest.h>

class Test_ModuleWidget : public QObject, public ::testing::Test
{};

TEST(Test_ModuleWidget_DeathTest, updateSliderValue)
{
    ModuleWidget widget;
#ifdef QT_DEBUG
    EXPECT_DEBUG_DEATH({widget.updateSliderValue(-1);}, "");
#endif
}

TEST_F(Test_ModuleWidget, updateSliderValue)
{
    ModuleWidget widget;

    widget.updateSliderValue(0);
    widget.updateSliderValue(1);
}
