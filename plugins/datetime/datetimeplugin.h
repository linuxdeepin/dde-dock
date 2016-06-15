#ifndef DATETIMEPLUGIN_H
#define DATETIMEPLUGIN_H

#include <QObject>

#include "pluginsiteminterface.h"

class DatetimePlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "datetime.json")

public:
    explicit DatetimePlugin(QObject *parent = 0);
    PluginsItem *getPluginsItem();
};

#endif // DATETIMEPLUGIN_H
