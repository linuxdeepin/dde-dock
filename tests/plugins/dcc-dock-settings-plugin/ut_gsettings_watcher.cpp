#include "config_watcher.h"

#include <QWidget>

#include <gtest/gtest.h>

using namespace dcc_dock_plugin;

class Test_GSettingWatcher : public QObject, public ::testing::Test
{};

TEST_F(Test_GSettingWatcher, bind)
{
    ConfigWatcher watcher("org.deepin.dock.plugin");

    QWidget widget;
    watcher.bind("dockPlugins", &widget);
    watcher.bind("dockPlugins", nullptr);
    watcher.bind("invalid", &widget);
    watcher.bind("", &widget);
    watcher.bind("", nullptr);
}

TEST_F(Test_GSettingWatcher, setStatus)
{
    ConfigWatcher watcher("org.deepin.dock.plugin");

    QWidget widget;
    watcher.bind("dockPlugins", &widget);
    watcher.setStatus("dockPlugins", &widget);
}

TEST_F(Test_GSettingWatcher, onStatusModeChanged)
{
    ConfigWatcher watcher("org.deepin.dock.plugin");

    QWidget widget;
    watcher.bind("dockPlugins", &widget);
    watcher.onStatusModeChanged("dockPlugins");
    watcher.onStatusModeChanged("invalid");
    watcher.onStatusModeChanged("");
}
