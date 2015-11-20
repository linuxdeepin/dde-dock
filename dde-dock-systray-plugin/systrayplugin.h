#ifndef SYSTRAYPLUGIN_H
#define SYSTRAYPLUGIN_H

#include <QtPlugin>
#include <QStringList>
#include <QWindow>
#include <QWidget>

#include "interfaces/dockconstants.h"
#include "interfaces/dockplugininterface.h"
#include "interfaces/dockpluginproxyinterface.h"

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

    QString getPluginName() Q_DECL_OVERRIDE;

    QStringList ids() Q_DECL_OVERRIDE;
    QString getTitle(QString id) Q_DECL_OVERRIDE;
    QString getName(QString id) Q_DECL_OVERRIDE;
    QString getCommand(QString id) Q_DECL_OVERRIDE;
    bool canDisable(QString id) Q_DECL_OVERRIDE;
    bool isDisabled(QString id) Q_DECL_OVERRIDE;
    void setDisabled(QString id, bool disabled) Q_DECL_OVERRIDE;
    QWidget * getItem(QString id) Q_DECL_OVERRIDE;
    QWidget * getApplet(QString id) Q_DECL_OVERRIDE;
    void changeMode(Dock::DockMode newMode, Dock::DockMode oldMode) Q_DECL_OVERRIDE;

    QString getMenuContent(QString id) Q_DECL_OVERRIDE;
    void invokeMenuItem(QString id, QString itemId, bool checked) Q_DECL_OVERRIDE;

private slots:
    void onTrayIconsChanged();
    void onTrayInit();

private:
    void initTrayIcons();
    void addTrayIcon(WId winId);
    void removeTrayIcon(WId winId);

private:
    CompositeTrayItem * m_compositeItem = 0;
    DockPluginProxyInterface * m_proxy = 0;
    com::deepin::dde::TrayManager *m_dbusTrayManager = 0;
};

#endif // SYSTRAYPLUGIN_H
