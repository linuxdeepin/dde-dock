#ifndef APPMANAGER_H
#define APPMANAGER_H

#include <QObject>
#include <QDebug>
#include "dbus/dbusentrymanager.h"
#include "dbus/dbusdockedappmanager.h"
#include "widgets/appitem.h"
#include "widgets/launcheritem.h"

class AppManager : public QObject
{
    Q_OBJECT
public:
    explicit AppManager(QObject *parent = 0);
    void initEntries();

    void setDockingItemId(const QString &dockingItemId);

signals:
    void entryAdded(AbstractDockItem *item, bool delayShow);
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
    QMap<QString, AbstractDockItem *> m_initItemList; //Juse for initialization <id, item>
};

#endif // APPMANAGER_H
