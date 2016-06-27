#ifndef DATETIMEPLUGIN_H
#define DATETIMEPLUGIN_H

#include "pluginsiteminterface.h"
#include "datetimewidget.h"

#include <QTimer>

class DatetimePlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "datetime.json")

public:
    explicit DatetimePlugin(QObject *parent = 0);
    ~DatetimePlugin();

    const QString pluginName();
    PluginType pluginType(const QString &itemKey);
    void init(PluginProxyInterface *proxyInter);

    QWidget *itemWidget(const QString &itemKey);

private:
    DatetimeWidget *m_centeralWidget;
    QTimer *m_refershTimer;
};

#endif // DATETIMEPLUGIN_H
