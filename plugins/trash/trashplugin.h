#ifndef TRASHPLUGIN_H
#define TRASHPLUGIN_H

#include "pluginsiteminterface.h"
#include "trashwidget.h"

class TrashPlugin : public QObject, PluginsItemInterface
{
    Q_OBJECT
    Q_INTERFACES(PluginsItemInterface)
    Q_PLUGIN_METADATA(IID "com.deepin.dock.PluginsItemInterface" FILE "trash.json")

public:
    explicit TrashPlugin(QObject *parent = 0);

    const QString pluginName() const;
    void init(PluginProxyInterface *proxyInter);

    QWidget *itemWidget(const QString &itemKey);
    QWidget *itemPopupApplet(const QString &itemKey);
    const QString itemCommand(const QString &itemKey);

private:
    TrashWidget *m_trashWidget;
};

#endif // TRASHPLUGIN_H
