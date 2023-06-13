// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "waylandmanager.h"
#include "taskmanager.h"
#include "taskmanager/entry.h"
#include "xcbutils.h"

#define XCB XCBUtils::instance()

WaylandManager::WaylandManager(TaskManager *_taskmanager, QObject *parent)
 : QObject(parent)
 , m_taskmanager(_taskmanager)
 , m_mutex(QMutex(QMutex::NonRecursive))
{

}


/**
 * @brief WaylandManager::registerWindow 注册窗口
 * @param objPath
 */
void WaylandManager::registerWindow(const QString &objPath)
{
    qInfo() << "registerWindow: " << objPath;
    if (findWindowByObjPath(objPath))
        return;

    PlasmaWindow *plasmaWindow = m_taskmanager->createPlasmaWindow(objPath);
    if (!plasmaWindow) {
        qWarning() << "registerWindowWayland: createPlasmaWindow failed";
        return;
    }

    if (!plasmaWindow->IsValid() || !plasmaWindow->isValid()) {
        qWarning() << "PlasmaWindow is not valid:" << objPath;
        plasmaWindow->deleteLater();
        return;
    }

    QString appId = plasmaWindow->AppId();
    QStringList list {"dde-dock", "dde-launcher", "dde-clipboard", "dde-osd", "dde-polkit-agent", "dde-simple-egl", "dmcs", "dde-lock"};
    if (list.indexOf(appId) >= 0 || appId.startsWith("No such object path")) {
        plasmaWindow->deleteLater();
        return;
    }

    XWindow winId = XCB->allocId();     // XCB中未发现释放XID接口
    XWindow realId = plasmaWindow->WindowId();
    if (realId)
        winId = realId;

    WindowInfoK *winInfo = new WindowInfoK(plasmaWindow, winId);
    m_taskmanager->listenKWindowSignals(winInfo);
    insertWindow(objPath, winInfo);
    m_taskmanager->attachOrDetachWindow(winInfo);
    if (winId) {
        m_windowInfoMap[winId] = winInfo;
    }
}

// 取消注册窗口
void WaylandManager::unRegisterWindow(const QString &objPath)
{
    qInfo() << "unRegisterWindow: " << objPath;
    WindowInfoK *winInfo = findWindowByObjPath(objPath);
    if (!winInfo)
        return;

    m_taskmanager->removePlasmaWindowHandler(winInfo->getPlasmaWindow());
    m_taskmanager->detachWindow(winInfo);
    deleteWindow(objPath);
}

WindowInfoK *WaylandManager::findWindowById(uint activeWin)
{
    QMutexLocker locker(&m_mutex);
    for (auto iter = m_kWinInfos.begin(); iter != m_kWinInfos.end(); iter++) {
        if (iter.value()->getInnerId() == QString::number(activeWin)) {
            return iter.value();
        }
    }

    return nullptr;
}

WindowInfoK *WaylandManager::findWindowByXid(XWindow xid)
{
    WindowInfoK *winInfo = nullptr;
    for (auto iter = m_kWinInfos.begin(); iter != m_kWinInfos.end(); iter++) {
        if (iter.value()->getXid() == xid) {
            winInfo = iter.value();
            break;
        }
    }

    return winInfo;
}

WindowInfoK *WaylandManager::findWindowByObjPath(QString objPath)
{
    if (m_kWinInfos.find(objPath) == m_kWinInfos.end())
        return nullptr;

    return m_kWinInfos[objPath];
}

void WaylandManager::insertWindow(QString objPath, WindowInfoK *windowInfo)
{
    QMutexLocker locker(&m_mutex);
    m_kWinInfos[objPath] = windowInfo;
}

void WaylandManager::deleteWindow(QString objPath)
{
    m_kWinInfos.remove(objPath);
}
