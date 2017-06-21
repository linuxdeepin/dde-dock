#ifndef SHUTDOWNPLUGIN_H
#define SHUTDOWNPLUGIN_H

#include "pluginsiteminterface.h"
#include "pluginwidget.h"
#include "powerstatuswidget.h"
#include "dbus/dbuspower.h"

#include <QLabel>

#define BATTERY_DISCHARED   2
#define BATTERY_FULL        4

class ShutdownPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "shutdown.json")

public:
    explicit ShutdownPlugin(QObject *parent = 0);

    const QString pluginName() const;
    void init(PluginProxyInterface *proxyInter);

    QWidget *itemWidget(const QString &itemKey);
    QWidget *itemTipsWidget(const QString &itemKey);
    const QString itemCommand(const QString &itemKey);
    const QString itemContextMenu(const QString &itemKey);
    void invokedMenuItem(const QString &itemKey, const QString &menuId, const bool checked);
    void displayModeChanged(const Dock::DisplayMode displayMode);

private:
    void updateBatteryVisible();
    void requestContextMenu(const QString &itemKey);
    void delayLoader();

private:
    PluginWidget *m_shutdownWidget;
    PowerStatusWidget *m_powerStatusWidget;
    QLabel *m_tipsLabel;

    DBusPower *m_powerInter;
};

#endif // SHUTDOWNPLUGIN_H
