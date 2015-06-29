#ifndef DOCKPLUGININTERFACE_H
#define DOCKPLUGININTERFACE_H

#include <QObject>
#include "abstractdockitem.h"

class DockPluginInterface : public QObject
{
    Q_OBJECT
public:
    virtual ~DockPluginInterface() {}
    virtual QList<AbstractDockItem*> items();
};

QT_BEGIN_NAMESPACE

#define DockPluginInterface_iid "org.deepin.Dock.PluginInterface"

Q_DECLARE_INTERFACE(DockPluginInterface, DockPluginInterface_iid)

QT_END_NAMESPACE

#endif // DOCKPLUGININTERFACE_H
