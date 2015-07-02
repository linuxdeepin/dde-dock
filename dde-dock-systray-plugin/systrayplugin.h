#ifndef SYSTRAYPLUGIN_H
#define SYSTRAYPLUGIN_H

#include <QtPlugin>
#include <QStringList>

#include "docktrayitem.h"
#include "dockplugininterface.h"
#include "dbustraymanager.h"

class SystrayPlugin : public QObject, DockPluginInterface
{
    Q_OBJECT
    Q_PLUGIN_METADATA(IID "org.deepin.Dock.PluginInterface" FILE "systray.json")
    Q_INTERFACES(DockPluginInterface)

public:
    ~SystrayPlugin();

    void init() Q_DECL_OVERRIDE;
    QStringList uuids() Q_DECL_OVERRIDE;
    QWidget * getItem(QString uuid) Q_DECL_OVERRIDE;

private:
    QMap<QString, QWidget*> m_items;
    com::deepin::dde::TrayManager *m_dbusTrayManager = 0;

    void clearItems();
};

#endif // SYSTRAYPLUGIN_H
