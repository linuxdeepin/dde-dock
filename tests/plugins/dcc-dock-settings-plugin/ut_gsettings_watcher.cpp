#include "config_watcher.h"

#include <QWidget>

#include <gtest/gtest.h>

using namespace dcc_dock_plugin;

class Test_GSettingWatcher : public QObject, public ::testing::Test
{};

TEST_F(Test_GSettingWatcher, bind)
{
    ConfigWatcher watcher("dde.dock.plugin.dconfig");

    QWidget widget;
    watcher.bind("Control-Center_Dock_Plugins", &widget);
    watcher.bind("Control-Center_Dock_Plugins", nullptr);
    watcher.bind("invalid", &widget);
    watcher.bind("", &widget);
    watcher.bind("", nullptr);
}

TEST_F(Test_GSettingWatcher, setStatus)
{
    ConfigWatcher watcher("dde.dock.plugin.dconfig");

    QWidget widget;
    watcher.bind("Control-Center_Dock_Plugins", &widget);
    watcher.setStatus("Control-Center_Dock_Plugins", &widget);
}

TEST_F(Test_GSettingWatcher, onStatusModeChanged)
{
    ConfigWatcher watcher("dde.dock.plugin.dconfig");

    QWidget widget;
    watcher.bind("Control-Center_Dock_Plugins", &widget);
    watcher.onStatusModeChanged("Control-Center_Dock_Plugins");
    watcher.onStatusModeChanged("invalid");
    watcher.onStatusModeChanged("");
}
