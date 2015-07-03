#include "appmanager.h"

AppManager::AppManager(QObject *parent) : QObject(parent)
{
    m_entryManager = new DBusEntryManager(this);
    connect(m_entryManager, SIGNAL(Added(QDBusObjectPath)),this, SLOT(slotEntryAdded(QDBusObjectPath)));
    connect(m_entryManager, SIGNAL(Removed(QString)), this, SLOT(slotEntryRemoved(QString)));
}

void AppManager::updateEntries()
{
    QList<QDBusObjectPath> entryList = m_entryManager->entries();
    for (int i = 0; i < entryList.count(); i ++)
    {
        DBusEntryProxyer *dep = new DBusEntryProxyer(entryList.at(i).path());
        if (dep->isValid() && dep->type() == "App")
        {
            AppItem *item = new AppItem();
            item->setEntryProxyer(dep);
            emit entryAdded(item);
        }
    }
}

void AppManager::slotEntryAdded(const QDBusObjectPath &path)
{
    qWarning() << "entry add:" << path.path();
    DBusEntryProxyer *entryProxyer = new DBusEntryProxyer(path.path());
    if (entryProxyer->isValid())
    {
//        qWarning() << "app-status" << entryProxyer->data().value("app-status");
//        qWarning() << "app-xids" << entryProxyer->data().value("app-xids");
//        qWarning() << "icon" << entryProxyer->data().value("icon");
//        qWarning() << "menu" << entryProxyer->data().value("menu");
//        qWarning() << "title" << entryProxyer->data().value("title");

        AppItem *item = new AppItem();
        item->setEntryProxyer(entryProxyer);
        emit entryAdded(item);
    }
}

void AppManager::slotEntryRemoved(const QString &id)
{
    qWarning() << "entry remove:" << id;
    emit entryRemoved(id);
}

