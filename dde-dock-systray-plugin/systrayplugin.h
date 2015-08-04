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
    void changeMode(Dock::DockMode newMode, Dock::DockMode oldMode);

    QString getMenuContent(QString uuid);
    void invokeMenuItem(QString uuid, QString itemId, bool checked);

private:
    QMap<QString, QWidget*> m_items;
    DockPluginProxyInterface * m_proxy = 0;
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
