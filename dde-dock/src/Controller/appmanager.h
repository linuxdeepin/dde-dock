#ifndef APPMANAGER_H
#define APPMANAGER_H

#include <QObject>
#include <QDebug>
#include "DBus/dbusentrymanager.h"
#include "DBus/dbusentryproxyer.h"
#include "Widgets/appitem.h"

class AppManager : public QObject
{
    Q_OBJECT
public:
    explicit AppManager(QObject *parent = 0);
    void updateEntries();

signals:
    void entryAdded(AppItem *item);
    void entryRemoved(const QString &id);

private slots:
    void slotEntryAdded(const QDBusObjectPath &path);
    void slotEntryRemoved(const QString &id);

private:


private:
    DBusEntryManager *m_entryManager = NULL;
};

#endif // APPMANAGER_H
