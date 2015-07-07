#ifndef SYSTRAYPLUGIN_H
#define SYSTRAYPLUGIN_H

#include <QtPlugin>
#include <QStringList>
#include <QWindow>
#include <QWidget>

#include "dockconstants.h"
#include "dockplugininterface.h"
#include "dockpluginproxyinterface.h"
#include "dbustraymanager.h"

class SystrayPlugin : public QObject, public DockPluginInterface
{
    Q_OBJECT
    Q_PLUGIN_METADATA(IID "org.deepin.Dock.PluginInterface" FILE "systray.json")
    Q_INTERFACES(DockPluginInterface)

public:
    ~SystrayPlugin();

    void init(DockPluginProxyInterface * proxier) Q_DECL_OVERRIDE;

    QStringList uuids() Q_DECL_OVERRIDE;
    QWidget * getItem(QString uuid) Q_DECL_OVERRIDE;
    void changeMode(Dock::DockMode newMode, Dock::DockMode oldMode);

    QString name() Q_DECL_OVERRIDE;

private:
    QMap<QString, QWidget*> m_items;
    DockPluginProxyInterface * m_proxier = 0;
    com::deepin::dde::TrayManager *m_dbusTrayManager = 0;
    Dock::DockMode m_mode;

    void clearItems();
    void addItem(QString uuid, QWidget * item);
    void removeItem(QString uuid);

private slots:
    void onAdded(WId winId);
    void onRemoved(WId winId);
};

#endif // SYSTRAYPLUGIN_H
