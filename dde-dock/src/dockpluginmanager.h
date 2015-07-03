#ifndef DOCKPLUGINMANAGER_H
#define DOCKPLUGINMANAGER_H

#include <QObject>
#include <QMap>
#include <qstringlist.h>

class DockPluginProxy;
class DockPluginManager : public QObject
{
    Q_OBJECT
public:
    explicit DockPluginManager(QObject *parent = 0);

    QList<DockPluginProxy*> getAll();

private:
    QStringList m_searchPaths;
    QMap<QString, DockPluginProxy*> m_proxies;

    DockPluginProxy* loadPlugin(QString &path);
};

#endif // DOCKPLUGINMANAGER_H
