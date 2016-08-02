#ifndef PLUGINSITEMINTERFACE_H
#define PLUGINSITEMINTERFACE_H

#include "pluginproxyinterface.h"

#include <QIcon>
#include <QtCore>

class PluginsItemInterface
{
public:
    virtual ~PluginsItemInterface() {}

    // the unique plugin id
    virtual const QString pluginName() const = 0;
    // init plugins
    virtual void init(PluginProxyInterface *proxyInter) = 0;
    // plugin item widget
    virtual QWidget *itemWidget(const QString &itemKey) = 0;

    virtual QWidget *itemTipsWidget(const QString &itemKey) {Q_UNUSED(itemKey); return nullptr;}
    virtual QWidget *itemPopupApplet(const QString &itemKey) {Q_UNUSED(itemKey); return nullptr;}
    virtual const QString itemCommand(const QString &itemKey) {Q_UNUSED(itemKey); return QString();}

    // item sort key
    virtual int itemSortKey(const QString &itemKey) {Q_UNUSED(itemKey); return 0;}
    // reset sort key when plugins order changed
    virtual void setSortKey(const QString &itemKey, const int order) {Q_UNUSED(itemKey); Q_UNUSED(order);}

    // dock display mode changed
    virtual void displayModeChanged(const Dock::DisplayMode displayMode) {Q_UNUSED(displayMode);}
    // dock position changed
    virtual void positionChanged(const Dock::Position position) {Q_UNUSED(position);}


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
