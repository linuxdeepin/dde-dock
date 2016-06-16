#include "datetimeplugin.h"

DatetimePlugin::DatetimePlugin(QObject *parent)
    : QObject(parent)
{

}

const QString DatetimePlugin::name()
{
    return "datetime";
}
