#ifndef SYSTRAYPLUGIN_H
#define SYSTRAYPLUGIN_H

#include <QtPlugin>
#include <QStringList>
#include <QWindow>
#include <QWidget>

#include <dock/dockconstants.h>
#include <dock/dockplugininterface.h>
#include <dock/dockpluginproxyinterface.h>

#include "dbustraymanager.h"

class CompositeTrayItem;
class SystrayPlugin : public QObject, public DockPluginInterface
{
    Q_OBJECT
    Q_PLUGIN_METADATA(IID "org.deepin.Dock.PluginInterface" FILE "dde-dock-systray-plugin.json")
    Q_INTERFACES(DockPluginInterface)

public:
    SystrayPlugin();
    ~SystrayPlugin();

    void init(DockPluginProxyInterface * proxy) Q_DECL_OVERRIDE;

    QString name() Q_DECL_OVERRIDE;

    QStringList uuids() Q_DECL_OVERRIDE;
    QString getTitle(QString uuid) Q_DECL_OVERRIDE;
    QWidget * getItem(QString uuid) Q_DECL_OVERRIDE;
    QWidget * getApplet(QString uuid) Q_DECL_OVERRIDE;
    void changeMode(Dock::DockMode newMode, Dock::DockMode oldMode) Q_DECL_OVERRIDE;

    QString getMenuContent(QString uuid) Q_DECL_OVERRIDE;
    void invokeMenuItem(QString uuid, QString itemId, bool checked) Q_DECL_OVERRIDE;

private:
    CompositeTrayItem * m_compositeItem = 0;
    DockPluginProxyInterface * m_proxy = 0;
    com::deepin::dde::TrayManager *m_dbusTrayManager = 0;

private slots:
    void onAdded(WId winId);
    void onRemoved(WId winId);
};

#endif // SYSTRAYPLUGIN_H
