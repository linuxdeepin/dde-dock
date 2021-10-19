#include "gsetting_watcher.h"

#include <QWidget>

#include <gtest/gtest.h>

class Test_GSettingWatcher : public QObject, public ::testing::Test
{};

TEST_F(Test_GSettingWatcher, bind)
{
    GSettingWatcher watcher("com.deepin.dde.control-center", "personalization");

    QWidget widget;
    watcher.bind("displayMode", &widget);
    watcher.bind("displayMode", nullptr);
    watcher.bind("invalid", &widget);
    watcher.bind("", &widget);
    watcher.bind("", nullptr);
}

TEST_F(Test_GSettingWatcher, setStatus)
{
    GSettingWatcher watcher("com.deepin.dde.control-center", "personalization");

    QWidget widget;
    watcher.bind("displayMode", &widget);
    watcher.setStatus("displayMode", &widget);
}

TEST_F(Test_GSettingWatcher, onStatusModeChanged)
{
    GSettingWatcher watcher("com.deepin.dde.control-center", "personalization");

    QWidget widget;
    watcher.bind("displayMode", &widget);
    watcher.onStatusModeChanged("displayMode");
    watcher.onStatusModeChanged("invalid");
    watcher.onStatusModeChanged("");
}
