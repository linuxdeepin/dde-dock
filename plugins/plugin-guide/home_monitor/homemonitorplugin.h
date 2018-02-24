#ifndef HOMEMONITORPLUGIN_H
#define HOMEMONITORPLUGIN_H

#include "informationwidget.h"

#include <QObject>

#include <dde-dock/pluginsiteminterface.h>

class HomeMonitorPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "home_monitor.json")

public:
    explicit HomeMonitorPlugin(QObject *parent = nullptr);

    const QString pluginName() const override;
    void init(PluginProxyInterface *proxyInter) override;

    QWidget *itemWidget(const QString &itemKey) override;

private:
    InformationWidget *m_pluginWidget;
};

#endif // HOMEMONITORPLUGIN_H
