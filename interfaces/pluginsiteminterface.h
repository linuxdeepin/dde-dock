#ifndef PLUGINSITEMINTERFACE_H
#define PLUGINSITEMINTERFACE_H

#include "pluginproxyinterface.h"

#include <QIcon>
#include <QtCore>

class PluginsItemInterface
{
public:
    enum PluginType
    {
        Simple,
        Complex,
    };

public:
    virtual ~PluginsItemInterface() {}

    // the unique plugin id
    virtual const QString pluginName() = 0;
    // plugins type, simple icon or complex widget
    virtual PluginType pluginType(const QString &itemKey) = 0;
    // init plugins
    virtual void init(PluginProxyInterface *proxyInter) = 0;

    // if complex widget mode, only return widget to plugins item
    virtual QWidget *itemWidget(const QString &itemKey) {Q_UNUSED(itemKey); return nullptr;}
    // in simple icon mode, plugins need to implements some data source functions
    virtual const QIcon itemIcon(const QString &itemKey) {Q_UNUSED(itemKey); return QIcon();}

protected:
    PluginProxyInterface *m_proxyInter;
};

QT_BEGIN_NAMESPACE

#define ModuleInterface_iid "com.deepin.dock.PluginsItemInterface"

Q_DECLARE_INTERFACE(PluginsItemInterface, ModuleInterface_iid)
QT_END_NAMESPACE

#endif // PLUGINSITEMINTERFACE_H
