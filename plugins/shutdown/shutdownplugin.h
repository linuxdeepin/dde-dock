#ifndef SHUTDOWNPLUGIN_H
#define SHUTDOWNPLUGIN_H

#include "pluginsiteminterface.h"
#include "pluginwidget.h"
#include "dbus/dbuspower.h"

#include <QLabel>

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
    void displayModeChanged(const Dock::DisplayMode displayMode);

private:
    PluginWidget *m_pluginWidget;
    QLabel *m_tipsLabel;

    DBusPower *m_powerInter;
};

#endif // SHUTDOWNPLUGIN_H
