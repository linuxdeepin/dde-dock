#ifndef APPMANAGER_H
#define APPMANAGER_H

#include <QObject>
#include <QDebug>
#include "DBus/dbusentrymanager.h"
#include "DBus/dbusentryproxyer.h"
#include "DBus/dbusdockedappmanager.h"
#include "Widgets/appitem.h"
#include "Widgets/launcheritem.h"

class AppManager : public QObject
{
    Q_OBJECT
public:
    explicit AppManager(QObject *parent = 0);
    void updateEntries();

signals:
    void entryAdded(AbstractDockItem *item);
    void entryRemoved(const QString &id);

private slots:
    void slotEntryAdded(const QDBusObjectPath &path);
    void slotEntryRemoved(const QString &id);

private:
    void sortItemList();    //Sort and append item to dock

private:
    QMap<QString, AbstractDockItem *> m_initItemList; //Juse for initialization <id, item>
    DBusEntryManager *m_entryManager = NULL;
    DBusDockedAppManager *m_ddam = new DBusDockedAppManager(this);
};

#endif // APPMANAGER_H
