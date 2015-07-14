#ifndef DOCKPLUGININTERFACE_H
#define DOCKPLUGININTERFACE_H

#include <QObject>
#include <QStringList>

#include "dockconstants.h"
#include "dockpluginproxyinterface.h"

class DockPluginInterface
{
public:
    virtual ~DockPluginInterface() {}
    virtual void init(DockPluginProxyInterface *proxy) = 0;

    virtual QString name() = 0;

    virtual QStringList uuids() = 0;
    virtual QString getTitle(QString uuid) = 0;
    virtual QWidget * getItem(QString uuid) = 0;
    virtual QWidget * getApplet(QString uuid) = 0;
    virtual void changeMode(Dock::DockMode newMode, Dock::DockMode oldMode) = 0;

    virtual QString getMenuContent(QString uuid) = 0;
    virtual void invokeMenuItem(QString uuid, QString itemId, bool checked) = 0;
};

QT_BEGIN_NAMESPACE

#define DockPluginInterface_iid "org.deepin.Dock.PluginInterface"

Q_DECLARE_INTERFACE(DockPluginInterface, DockPluginInterface_iid)

QT_END_NAMESPACE

#endif // DOCKPLUGININTERFACE_H
