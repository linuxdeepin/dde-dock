#ifndef DOCKPLUGININTERFACE_H
#define DOCKPLUGININTERFACE_H

#include <QPixmap>
#include <QObject>
#include <QStringList>

#include "dockconstants.h"
#include "dockpluginproxyinterface.h"

class DockPluginInterface
{
public:
    virtual ~DockPluginInterface() {}

    virtual QString getPluginName() = 0;

    virtual void init(DockPluginProxyInterface *proxy) = 0;
    virtual void changeMode(Dock::DockMode newMode, Dock::DockMode oldMode) = 0;

    virtual QStringList ids() = 0;
    virtual QString getName(QString id) = 0;
    virtual QString getTitle(QString id) = 0;
    virtual QString getCommand(QString id) = 0;
    virtual QPixmap getIcon(QString id) {Q_UNUSED(id); return QPixmap("");}
    virtual bool canDisable(QString id) = 0;
    virtual bool isDisabled(QString id) = 0;
    virtual void setDisabled(QString id, bool disabled) = 0;
    virtual QWidget * getItem(QString id) = 0;
    virtual QWidget * getApplet(QString id) = 0;
    virtual QString getMenuContent(QString id) = 0;
    virtual void invokeMenuItem(QString id, QString itemId, bool checked) = 0;
};

QT_BEGIN_NAMESPACE

#define DockPluginInterface_iid "org.deepin.Dock.PluginInterface"

Q_DECLARE_INTERFACE(DockPluginInterface, DockPluginInterface_iid)

QT_END_NAMESPACE

#endif // DOCKPLUGININTERFACE_H
