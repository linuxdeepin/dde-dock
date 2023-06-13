// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "entries.h"
#include "taskmanager.h"
#include "docksettings.h"
#include "taskmanager/windowinfobase.h"

#include <QList>
#include <algorithm>
#include <iterator>

Entries::Entries(TaskManager *_taskmanager)
 : m_taskmanager(_taskmanager)
{

}

QVector<Entry *> Entries::filterDockedEntries()
{
    QVector<Entry *> ret;
    for (auto entry : m_items) {
        if (entry->isValid() && entry->getIsDocked()) ret.push_back(entry);
    }
    return ret;
}

Entry *Entries::getByInnerId(QString innerId)
{
    Entry *ret = nullptr;
    for (auto &entry : m_items) {
        if (entry->getInnerId() == innerId)
            ret = entry;
    }

    return ret;
}

void Entries::append(Entry *entry)
{
    insert(entry, -1);
}

void Entries::insert(Entry *entry, int index)
{
    // 如果当前应用在列表中存在(通常是该应用为最近打开应用但是关闭了最近打开应用的接口或者当前为高效模式)
    if (m_items.contains(entry))
        m_items.removeOne(entry);

    if (index < 0 || index >= m_items.size()) {
        // append
        index = m_items.size();
        m_items.push_back(entry);
    } else {
        // insert
        m_items.insert(index, entry);
    }

    insertCb(entry, index);
}

void Entries::remove(Entry *entry)
{
    for (auto iter = m_items.begin(); iter != m_items.end();) {
        if ((*iter)->getId() == entry->getId()) {
            iter = m_items.erase(iter);
        } else {
            iter++;
        }
    }

    removeCb(entry);
    entry->deleteLater();
}

void Entries::move(int oldIndex, int newIndex)
{
    if (oldIndex == newIndex || oldIndex < 0 || newIndex < 0 || oldIndex >= m_items.size() || newIndex >= m_items.size())
        return;

    m_items.swapItemsAt(oldIndex, newIndex);
}

Entry *Entries::getByWindowPid(int pid)
{
    Entry *ret = nullptr;
    for (auto &entry : m_items) {
        if (entry->getWindowInfoByPid(pid)) {
            ret = entry;
            break;
         }
    }

    return ret;
}

QStringList Entries::getEntryIDs()
{
    QStringList list;
    if (m_taskmanager->getDisplayMode() == DisplayMode::Fashion
            && DockSettings::instance()->showRecent()) {
        for (Entry *item : m_items) list << item->getId();
    } else {
        // 如果是高效模式或者没有开启显示最近应用的功能，那么未驻留并且没有子窗口的就不显示
        // 换句话说，只显示已经驻留或者有子窗口的应用
        for (Entry *item : m_items) {
            if (item->getIsDocked() || item->hasWindow()) 
                list << item->getId();
        }
    }

    return list;
}

Entry *Entries::getByWindowId(XWindow windowId)
{
    Entry *ret = nullptr;
    for (auto &entry : m_items) {
        if (entry->getWindowInfoByWinId(windowId)) {
            ret = entry;
            break;
         }
    }

    return ret;
}

Entry *Entries::getByDesktopFilePath(const QString &filePath)
{
    Entry *ret = nullptr;
    for (auto &entry : m_items) {
        qDebug() << entry->getName();
        if (entry->getFileName() == filePath) {
            ret = entry;
            break;
        }
    }

    return ret;
}

QList<Entry*> Entries::getEntries()
{
    QList<Entry*> list;
    if (static_cast<DisplayMode>(m_taskmanager->getDisplayMode()) == DisplayMode::Fashion
            && DockSettings::instance()->showRecent()) {
        for (Entry *item : m_items)
            list << item;
    } else {
        // 如果是高效模式或者没有开启显示最近应用的功能，那么未驻留并且没有子窗口的就不显示
        // 换句话说，只显示已经驻留或者有子窗口的应用
        for (Entry *item : m_items) {
            if (!item->getIsDocked() && !item->hasWindow())
                continue;
            list << item;
        }
    }

    return list;
}

Entry *Entries::getDockedEntryByDesktopFile(const QString &desktopFile)
{
    Entry *ret = nullptr;
    for (auto entry : filterDockedEntries()) {
        if ((entry->isValid()) && desktopFile == entry->getFileName()) {
            ret = entry;
            break;
        }
    }

    return ret;
}

QString Entries::queryWindowIdentifyMethod(XWindow windowId)
{
    QString ret;
    for (auto entry : m_items) {
        auto window = entry->getWindowInfoByWinId(windowId);
        if (window) {
            auto app = window->getAppInfo();
            ret = app ? app->getIdentifyMethod() : "Failed";
            break;
        }
    }

    return ret;
}

void Entries::handleActiveWindowChanged(XWindow activeWindId)
{
    for (auto entry : m_items) {
        auto windowInfo = entry->getWindowInfoByWinId(activeWindId);
        if (windowInfo) {
            entry->setPropIsActive(true);
            entry->setCurrentWindowInfo(windowInfo);
            entry->updateName();
            entry->updateIcon();
        } else {
            entry->setPropIsActive(false);
        }
    }
}

void Entries::updateEntriesMenu()
{
    for (auto entry : m_items) {
        entry->updateMenu();
    }
}

const QList<Entry *> Entries::unDockedEntries() const
{
    QList<Entry *> entrys;
    for (Entry *entry : m_items) {
        if (!entry->isValid() || entry->getIsDocked())
            continue;

        entrys << entry;
    }

    return entrys;
}

void Entries::moveEntryToLast(Entry *entry)
{
    if (m_items.contains(entry)) {
        m_items.removeOne(entry);
        m_items << entry;
    }
}

void Entries::insertCb(Entry *entry, int index)
{
    if (entry->getIsDocked() || entry->hasWindow() ||
            ((m_taskmanager->getDisplayMode() == DisplayMode::Fashion) && DockSettings::instance()->showRecent())){
                Q_EMIT m_taskmanager->entryAdded(entry, index);
            }
}

void Entries::removeCb(Entry *entry)
{
    Q_EMIT m_taskmanager->entryRemoved(entry->getId());
}

bool Entries::shouldInRecent()
{
    // 如果当前移除的应用是未驻留应用，则判断未驻留应用的数量是否小于等于3，则让其始终显示
    QList<Entry *> unDocktrys;
    for (Entry *entry : m_items) {
        if (entry->isValid() && !entry->getIsDocked())
            unDocktrys << entry;
    }

    // 如果当前未驻留应用的数量小于3个，则认为后续的应用应该显示到最近打开应用
    return (unDocktrys.size() <= MAX_UNOPEN_RECENT_COUNT);
}

void Entries::removeLastRecent()
{
    // 先查找最近使用的应用，删除没有使用的
    int unDockCount = 0;
    QList<Entry *> unDockEntrys;
    QList<Entry *> removeEntrys;

    for (Entry *entry : m_items) {
        if (entry->getIsDocked())
            continue;

        // 此处只移除没有子窗口的图标
        if (!entry->hasWindow()) {
            if (!entry->isValid())
                removeEntrys << entry; // 如果应用已经被卸载，那么需要删除
            else
                unDockEntrys << entry;
        }

        unDockCount++;
    }
    if (unDockCount >= MAX_UNOPEN_RECENT_COUNT && unDockEntrys.size() > 0) {
        // 只有当最近使用区域的图标大于等于某个数值（3）的时候，并且存在没有子窗口的Entry，那么就移除该Entry
        Entry *entry = unDockEntrys[0];
        removeEntrys << entry;
    }
    for (Entry *entry : removeEntrys) {
        m_items.removeOne(entry);
        removeCb(entry);
        entry->deleteLater();
    }
}

void Entries::setDisplayMode(DisplayMode displayMode)
{
    if (!DockSettings::instance()->showRecent())
        return;

    // 如果从时尚模式变成高效模式，对列表中所有的没有打开窗口的应用发送移除信号
    if (displayMode == DisplayMode::Efficient) {
        for (Entry *entry : m_items) {
            entry->updateMode();
            if (!entry->getIsDocked() && !entry->hasWindow())
                Q_EMIT m_taskmanager->entryRemoved(entry->getId());
        }
    } else {
        // 如果从高效模式变成时尚模式，列表中所有的未驻留且不存在打开窗口的应用认为是最近打开应用，发送新增信号
        for (Entry *entry : m_items) {
            entry->updateMode();
            if (!entry->getIsDocked() && !entry->hasWindow()) {
                // QString objPath = entry->path();
                int index = m_items.indexOf(entry);
                Q_EMIT m_taskmanager->entryAdded(entry, index);
                qDebug() << entry->getName();
            }
        }
    }
}

void Entries::updateShowRecent()
{
    // 高效模式无需做任何操作
    if (static_cast<DisplayMode>(m_taskmanager->getDisplayMode()) != DisplayMode::Fashion)
        return;

    bool showRecent = DockSettings::instance()->showRecent();
    if (showRecent) {
        // 如果显示最近打开应用，则发送新增信号
        for (Entry *entry : m_items) {
            // 已经驻留的或者有子窗口的本来就在任务栏上面，无需发送信号
            entry->updateMode();
            if (entry->getIsDocked() || entry->hasWindow())
                continue;

            // QString objPath = entry->path();
            int index = m_items.indexOf(entry);
            Q_EMIT m_taskmanager->entryAdded(entry, index);
        }
    } else {
        // 如果是隐藏最近打开的应用，则发送移除的信号
        for (Entry *entry : m_items) {
            // 已经驻留的或者有子窗口的本来就在任务栏上面，无需发送信号
            entry->updateMode();
            if (entry->getIsDocked() || entry->hasWindow())
                continue;

            Q_EMIT m_taskmanager->entryRemoved(entry->getId());
        }
    }
}
