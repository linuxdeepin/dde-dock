#include "config_watcher.h"

#include <QWidget>
#include <QApplication>

#include <gtest/gtest.h>

using namespace dcc_dock_plugin;

class Test_GSettingWatcher : public QObject, public ::testing::Test
{};

TEST_F(Test_GSettingWatcher, bind)
{
    const QString &appName = qApp->applicationName();
    qApp->setApplicationName("dde-dock");
    ConfigWatcher watcher("org.deepin.dde.dock.plugin");

    QWidget widget;
    watcher.bind("dockPlugins", &widget);
    watcher.bind("dockPlugins", nullptr);
    watcher.bind("invalid", &widget);
    watcher.bind("", &widget);
    watcher.bind("", nullptr);
    qApp->setApplicationName(appName);
}

TEST_F(Test_GSettingWatcher, setStatus)
{
    const QString &appName = qApp->applicationName();
    qApp->setApplicationName("dde-control-center");
    ConfigWatcher watcher("org.deepin.dde.dock.plugin");

    QWidget widget;
    watcher.bind("dockPlugins", &widget);
    watcher.setStatus("dockPlugins", &widget);
    qApp->setApplicationName(appName);
}

TEST_F(Test_GSettingWatcher, onStatusModeChanged)
{
    ConfigWatcher watcher("org.deepin.dde.dock.plugin");

    QWidget widget;
    watcher.bind("dockPlugins", &widget);
    watcher.onStatusModeChanged("dockPlugins");
    watcher.onStatusModeChanged("invalid");
    watcher.onStatusModeChanged("");
}
