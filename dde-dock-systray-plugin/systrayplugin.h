#ifndef SYSTRAYPLUGIN_H
#define SYSTRAYPLUGIN_H

#include <QtPlugin>

#include "docktrayitem.h"
#include "dockplugininterface.h"
#include "abstractdockitem.h"
#include "dbustraymanager.h"

class SystrayPlugin : public QObject, DockPluginInterface
{
    Q_OBJECT
    Q_PLUGIN_METADATA(IID "org.deepin.Dock.PluginInterface" FILE "systray.json")
    Q_INTERFACES(DockPluginInterface)

public:
    ~SystrayPlugin();

    QList<AbstractDockItem*> items() Q_DECL_OVERRIDE;

private:
    QList<AbstractDockItem*> m_items;
    com::deepin::dde::TrayManager *m_dbusTrayManager = 0;

    void clearItems();
};

#endif // SYSTRAYPLUGIN_H
