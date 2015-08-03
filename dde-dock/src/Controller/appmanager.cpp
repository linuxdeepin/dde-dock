#include "appmanager.h"

AppManager::AppManager(QObject *parent) : QObject(parent)
{
    m_entryManager = new DBusEntryManager(this);
    connect(m_entryManager, SIGNAL(Added(QDBusObjectPath)),this, SLOT(slotEntryAdded(QDBusObjectPath)));
    connect(m_entryManager, SIGNAL(Removed(QString)), this, SLOT(slotEntryRemoved(QString)));
}

void AppManager::updateEntries()
{

    LauncherItem * lItem = new LauncherItem();
    emit entryAdded(lItem);

    QList<QDBusObjectPath> entryList = m_entryManager->entries();
    for (int i = 0; i < entryList.count(); i ++)
    {
        DBusEntryProxyer *dep = new DBusEntryProxyer(entryList.at(i).path());
        if (dep->isValid() && dep->type() == "App" && dep->data().value("title") != "dde-dock")
        {
            AppItem *item = new AppItem();
            item->setEntryProxyer(dep);
            m_initItemList.insert(item->getItemId(), item);
//            emit entryAdded(item);
        }
    }

    sortItemList();
}

void AppManager::slotEntryAdded(const QDBusObjectPath &path)
{
    DBusEntryProxyer *entryProxyer = new DBusEntryProxyer(path.path());
    if (entryProxyer->isValid() && entryProxyer->type() == "App" && entryProxyer->data().value("title") != "dde-dock")
    {
        qWarning() << "entry add:" << path.path();
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

void AppManager::sortItemList()
{
    QStringList dockedList = m_ddam->DockedAppList().value();
    QStringList ids = m_initItemList.keys();
    foreach (QString id, dockedList) {  //For docked items
        int index = ids.indexOf(id);
        if (index != -1)
            emit entryAdded(m_initItemList.take(ids.at(index)));
    }
    ids = m_initItemList.keys();
    foreach (QString id, ids) { //For undocked items
        emit entryAdded(m_initItemList.take(id));
    }
}
