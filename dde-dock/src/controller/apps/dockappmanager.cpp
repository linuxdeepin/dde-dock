#include "dockappmanager.h"
#include "dbus/dbuslauncher.h"

DockAppManager::DockAppManager(QObject *parent) : QObject(parent)
{
    m_entryManager = new DBusEntryManager(this);
    connect(m_entryManager, &DBusEntryManager::Added, this, &DockAppManager::onEntryAdded);
    connect(m_entryManager, &DBusEntryManager::Removed, this, &DockAppManager::onEntryRemoved);
    DBusLauncher *dbusLauncher = new DBusLauncher(this);
    connect(dbusLauncher, &DBusLauncher::ItemChanged, [=](const QString &in0, ItemInfo in1){
        if (in0 == "deleted") {
            onEntryRemoved(in1.key);
            m_dockAppManager->RequestUndock(in1.key);
        }
    });
}

void DockAppManager::initEntries()
{
//    LauncherItem * lItem = new LauncherItem();
//    emit entryAdded(lItem, false);

    QList<QDBusObjectPath> entryList = m_entryManager->entries();
//    for (int i = 0; i < entryList.count(); i ++)
    for (QDBusObjectPath objPath : entryList)
    {
        DBusDockEntry *dep = new DBusDockEntry(objPath.path());
        if (dep->isValid() && dep->type() == "App") {
            DockAppItem *item = new DockAppItem();
            item->setEntryProxyer(dep);
            m_initItems.insert(item->getItemId(), item);
        }
    }

    sortItemList();
}

void DockAppManager::onEntryAdded(const QDBusObjectPath &path)
{
    DBusDockEntry *entryProxyer = new DBusDockEntry(path.path());
    if (entryProxyer->isValid() && entryProxyer->type() == "App")
    {
        DockAppItem *item = new DockAppItem();
        item->setEntryProxyer(entryProxyer);
        QString tmpId = item->getItemId();
        if (m_ids.indexOf(tmpId) != -1) {
            item->deleteLater();
        }else{
            qDebug() << "app entry add:" << tmpId;
            bool isTheDropOne = m_dockingItemId != tmpId;
            m_ids.append(tmpId);
            emit entryAdded(item, isTheDropOne);

            if (isTheDropOne)
                setDockingItemId("");
        }
    }
}

void DockAppManager::setDockingItemId(const QString &dockingItemId)
{
    m_dockingItemId = dockingItemId;
}

void DockAppManager::onEntryRemoved(const QString &id)
{
    if (m_ids.indexOf(id) != -1) {
        qDebug() << "app entry remove:" << id;
        m_ids.removeAll(id);
        emit entryRemoved(id);
    }
}

void DockAppManager::sortItemList()
{
    QStringList dockedList = m_dockAppManager->DockedAppList().value();
    m_ids = m_initItems.keys();
    QStringList tmpIds = m_initItems.keys();
    for (QString id : dockedList) {  //For docked items
        int index = tmpIds.indexOf(id);
        if (index != -1)
            emit entryAdded(m_initItems.take(tmpIds.at(index)), false);
    }
    tmpIds = m_initItems.keys();
    for (QString id : tmpIds) { //For undocked items
        emit entryAdded(m_initItems.take(id), false);
    }
}
