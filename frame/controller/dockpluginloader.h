#ifndef DOCKPLUGINLOADER_H
#define DOCKPLUGINLOADER_H

#include <QThread>

class DockPluginLoader : public QThread
{
    Q_OBJECT

public:
    explicit DockPluginLoader(QObject *parent);

signals:
    void finished() const;
    void pluginFounded(const QString &pluginFile) const;

protected:
    void run();
};

#endif // DOCKPLUGINLOADER_H
