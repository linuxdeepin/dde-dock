#ifndef PLUGINSITEMINTERFACE_H
#define PLUGINSITEMINTERFACE_H

#include "pluginproxyinterface.h"

#include <QIcon>
#include <QtCore>

class PluginsItemInterface
{
public:
    enum ItemType
    {
        Simple,
        Complex,
    };

public:
    virtual ~PluginsItemInterface() {}

    // the unique plugin id
    virtual const QString pluginName() const = 0;
    // init plugins
    virtual void init(PluginProxyInterface *proxyInter) = 0;
    // dock display mode changed
    virtual void displayModeChanged(const Dock::DisplayMode displayMode) {Q_UNUSED(displayMode);}
    // dock position changed
    virtual void positionChanged(const Dock::Position position) {Q_UNUSED(position);}

    // plugins type, simple icon or complex widget
    virtual ItemType pluginType(const QString &itemKey) {Q_UNUSED(itemKey); return Simple;}
    // simple string tips or popup widget
    virtual ItemType tipsType(const QString &itemKey) {Q_UNUSED(itemKey); return Simple;}
    // item sort key
    virtual int itemSortKey(const QString &itemKey) {Q_UNUSED(itemKey); return 0;}
    // reset sort key when plugins order changed
    virtual void setSortKey(const QString &itemKey, const int order) {Q_UNUSED(itemKey); Q_UNUSED(order);}

    // if pluginType is complex widget mode, return a widget to plugins item
    virtual QWidget *itemWidget(const QString &itemKey) {Q_UNUSED(itemKey); return nullptr;}
    // if pluginType is simple icon mode, plugins need to implements these data source functions
    virtual const QIcon itemIcon(const QString &itemKey) {Q_UNUSED(itemKey); return QIcon();}
    virtual const QString itemCommand(const QString &itemKey) {Q_UNUSED(itemKey); return QString();}

    // return simple string tips, call this function when tips type is Simple
    virtual const QString simpleTipsString(const QString &itemKey) {Q_UNUSED(itemKey); return QString();}
    // return complex widget tips, call this function when tips type is Complex
    virtual QWidget *complexTipsWidget(const QString &itemKey) {Q_UNUSED(itemKey); return nullptr;}

protected:
    Dock::DisplayMode displayMode() const
    {
        return qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();
    }

    Dock::Position position() const
    {
        return qApp->property(PROP_POSITION).value<Dock::Position>();
    }

protected:
    PluginProxyInterface *m_proxyInter;
};

QT_BEGIN_NAMESPACE

#define ModuleInterface_iid "com.deepin.dock.PluginsItemInterface"

Q_DECLARE_INTERFACE(PluginsItemInterface, ModuleInterface_iid)
QT_END_NAMESPACE

#endif // PLUGINSITEMINTERFACE_H
