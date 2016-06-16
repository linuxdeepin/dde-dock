#ifndef DOCKPLUGINSCONTROLLER_H
#define DOCKPLUGINSCONTROLLER_H

#include <QPluginLoader>
#include <QList>

class PluginsItemInterface;
class DockPluginsController : public QObject
{
    Q_OBJECT

public:
    explicit DockPluginsController(QObject *parent = 0);
    ~DockPluginsController();

signals:
    void pluginsInserted(PluginsItemInterface *interface) const;

private slots:
    void loadPlugins();

private:
    QList<PluginsItemInterface *> m_pluginsInterfaceList;
    QList<QPluginLoader *> m_pluginLoaderList;
};

#endif // DOCKPLUGINSCONTROLLER_H
