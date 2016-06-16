#ifndef PLUGINSITEMINTERFACE_H
#define PLUGINSITEMINTERFACE_H

class QString;
class PluginsItemInterface
{
public:
    virtual ~PluginsItemInterface() {}

    virtual const QString name() = 0;
};

QT_BEGIN_NAMESPACE

#define ModuleInterface_iid "com.deepin.dock.PluginsItemInterface"

Q_DECLARE_INTERFACE(PluginsItemInterface, ModuleInterface_iid)
QT_END_NAMESPACE

#endif // PLUGINSITEMINTERFACE_H
