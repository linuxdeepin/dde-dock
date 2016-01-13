#ifndef DOCKAPPMANAGER_H
#define DOCKAPPMANAGER_H

#include "dbus/dbusentrymanager.h"
#include "dbus/dbusdockentry.h"
#include "dbus/dbusdockedappmanager.h"
#include "widgets/app/dockappitem.h"

class DockAppManager : public QObject
{
    Q_OBJECT
public:
    explicit DockAppManager(QObject *parent = 0);
    void initEntries();

    void setDockingItemId(const QString &dockingItemId);

signals:
    void entryAdded(DockAppItem *item);
    void entryAppend(DockAppItem *item);
    void entryRemoved(const QString &id);

private:
    void sortItemList();    //Sort and append item to dock
    void onEntryRemoved(const QString &id);
    void onEntryAdded(const QDBusObjectPath &path);

private:
    QStringList m_ids;
    QString m_dockingItemId = "";   //drop to the dock one
    DBusEntryManager *m_entryManager = NULL;
    DBusDockedAppManager *m_dockAppManager = new DBusDockedAppManager(this);
    QMap<QString, DockAppItem *> m_initItems; //Just for initialization <id, item>
};

#endif // DOCKAPPMANAGER_H
