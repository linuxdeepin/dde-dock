#ifndef SYSTRAYPLUGIN_H
#define SYSTRAYPLUGIN_H

#include "docktrayitem.h"
#include "dockplugininterface.h"

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

    void clearItems();
};

#endif // SYSTRAYPLUGIN_H
