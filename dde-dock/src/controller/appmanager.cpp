#include "appmanager.h"
#include "dbus/dbuslauncher.h"

AppManager::AppManager(QObject *parent) : QObject(parent)
{
    m_entryManager = new DBusEntryManager(this);
    connect(m_entryManager, &DBusEntryManager::Added, this, &AppManager::onEntryAdded);
    connect(m_entryManager, &DBusEntryManager::Removed, this, &AppManager::onEntryRemoved);
    DBusLauncher *dbusLauncher = new DBusLauncher(this);
    connect(dbusLauncher, &DBusLauncher::ItemChanged, [=](const QString &in0, ItemInfo in1){
        if (in0 == "deleted") {
            onEntryRemoved(in1.key);
        }
    });
}

void AppManager::initEntries()
{

    LauncherItem * lItem = new LauncherItem();
    emit entryAdded(lItem, false);

    QList<QDBusObjectPath> entryList = m_entryManager->entries();
    for (int i = 0; i < entryList.count(); i ++)
    {
        DBusEntryProxyer *dep = new DBusEntryProxyer(entryList.at(i).path());
        if (dep->isValid() && dep->type() == "App")
        {
            AppItem *item = new AppItem();
            item->setEntryProxyer(dep);
            m_initItemList.insert(item->getItemId(), item);
        }
    }

    sortItemList();
}

void AppManager::onEntryAdded(const QDBusObjectPath &path)
{
    DBusEntryProxyer *entryProxyer = new DBusEntryProxyer(path.path());
    if (entryProxyer->isValid() && entryProxyer->type() == "App")
    {
        AppItem *item = new AppItem();
        item->setEntryProxyer(entryProxyer);
        QString tmpId = item->getItemId();
        if (m_ids.indexOf(tmpId) != -1){
            item->deleteLater();
        }else{
            qWarning() << "entry add:" << tmpId;
            bool isTheDropOne = m_dockingItemId != tmpId;
            m_ids.append(tmpId);
            emit entryAdded(item, isTheDropOne);

            if (isTheDropOne)
                setDockingItemId("");
        }
    }
}

void AppManager::setDockingItemId(const QString &dockingItemId)
{
    m_dockingItemId = dockingItemId;
}

void AppManager::onEntryRemoved(const QString &id)
{
    if (m_ids.indexOf(id) != -1) {
        qWarning() << "entry remove:" << id;
        m_ids.removeAll(id);
        emit entryRemoved(id);
    }
}

void AppManager::sortItemList()
{
    QStringList dockedList = m_dockAppManager->DockedAppList().value();
    m_ids = m_initItemList.keys();
    QStringList tmpIds = m_initItemList.keys();
    foreach (QString id, dockedList) {  //For docked items
        int index = tmpIds.indexOf(id);
        if (index != -1)
            emit entryAdded(m_initItemList.take(tmpIds.at(index)), false);
    }
    tmpIds = m_initItemList.keys();
    foreach (QString id, tmpIds) { //For undocked items
        emit entryAdded(m_initItemList.take(id), false);
    }
}
