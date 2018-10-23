#include "shutdowntrayloader.h"
#include "shutdowntraywidget.h"

#define ShutdownItemKey "system-tray-shutdown"

ShutdownTrayLoader::ShutdownTrayLoader(QObject *parent) : AbstractTrayLoader(QString(), parent)
{
}

void ShutdownTrayLoader::load()
{
    emit systemTrayAdded(ShutdownItemKey, new ShutdownTrayWidget);
}
