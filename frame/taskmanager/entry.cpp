// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "entry.h"
#include "docksettings.h"
#include "xcbutils.h"
#include "taskmanager.h"
#include "processinfo.h"
#include "windowinfomap.h"

#include <QDebug>
#include <QDBusInterface>

#include <algorithm>
#include <signal.h>

#define XCB XCBUtils::instance()

Entry::Entry(TaskManager *_taskmanager, AppInfo *_app, QString _innerId, QObject *parent)
    : QObject(parent)
    , m_isActive(false)
    , m_isDocked(false)
    , m_winIconPreferred(false)
    , m_innerId(_innerId)
    , m_adapterEntry(nullptr)
    , m_taskmanager(_taskmanager)
    , m_current(nullptr)
    , m_currentWindow(0)
{
    setAppInfo(_app);
    m_id = m_taskmanager->allocEntryId();
    m_mode = getCurrentMode();
    m_name = getName();
    m_icon = getIcon();
}

Entry::~Entry()
{
    for (auto winInfo : m_windowInfoMap) {
        if (winInfo) winInfo->deleteLater();
    }
    m_windowInfoMap.clear();

}

bool Entry::isValid()
{
    // desktopfile 无效时且没有窗口时，该entry是无效的
    // 虽然也就是desktop是无效时，但是当前存在窗口，该entry也是有效的。
    return m_isValid || m_current;
}

QString Entry::getId() const
{
    return m_id;
}

QString Entry::getName()
{
    QString ret = m_current ? m_current->getDisplayName() : QString();
    if (m_appInfo.isNull()) return ret;
    ret = m_appInfo->getName();
    return ret;
}

void Entry::updateName()
{
    setPropName(getName());
}

QString Entry::getIcon()
{
    QString ret;
    if (hasWindow()) {
        if (!m_current) {
            return ret;
        }

        // has window && current not nullptr
        if (m_winIconPreferred) {
            // try current window icon first
            ret = m_current->getIcon();
            if (ret.size() > 0) {
                return ret;
            }
        }

        if (m_appInfo) {
            m_icon = m_appInfo->getIcon();
            if (m_icon.size() > 0) {
                return m_icon;
            }
        }

        return m_current->getIcon();
    }

    if (m_appInfo) {
        // no window
        return m_appInfo->getIcon();
    }

    return ret;
}

QString Entry::getInnerId()
{
    return m_innerId;
}

void Entry::setInnerId(QString _innerId)
{
    qDebug() << "setting innerID from: " << m_innerId << " to: " << _innerId;
    m_innerId = _innerId;
}

QString Entry::getFileName()
{
    return m_appInfo.isNull() ? QString() : m_appInfo->getFileName();
}

AppInfo *Entry::getAppInfo()
{
    return m_appInfo.data();
}

void Entry::setAppInfo(AppInfo *appinfo)
{
    if (m_appInfo.data() == appinfo || appinfo == nullptr) {
        return;
    }

    m_appInfo.reset(appinfo);
    m_isValid = appinfo->isValidApp();
    m_winIconPreferred = !appinfo;
    setPropDesktopFile(appinfo ? appinfo->getFileName(): "");
    if (!m_winIconPreferred) {
        QString id = m_appInfo->getId();
        auto perferredApps = m_taskmanager->getWinIconPreferredApps();
        if (perferredApps.contains(id)|| appinfo->getIcon().size() == 0) {
            m_winIconPreferred = true;
            return;
        }
    }
}

bool Entry::getIsDocked() const
{
    return m_isDocked;
}

void Entry::setIsDocked(bool value)
{
    if (value != m_isDocked) {
        m_isDocked = value;
        Q_EMIT isDockedChanged(value);
    }
}

void Entry::setMenu(AppMenu *_menu)
{
    _menu->setDirtyStatus(true);
    m_appMenu.reset(_menu);
    Q_EMIT menuChanged(m_appMenu->getMenuJsonStr());
}

void Entry::updateMenu()
{
    qInfo() <<"Entry: updateMenu";
    AppMenu *appMenu = new AppMenu();
    appMenu->appendItem(getMenuItemLaunch());

    for (auto &item :getMenuItemDesktopActions())
        appMenu->appendItem(item);

    if (hasWindow())
        appMenu->appendItem(getMenuItemAllWindows());

    // menu item dock or undock
    qInfo() << "entry " << m_id << " docked? " << m_isDocked;
    appMenu->appendItem(m_isDocked? getMenuItemUndock(): getMenuItemDock());

    if (hasWindow()) {
        if (m_taskmanager->getForceQuitAppStatus() != ForceQuitAppMode::Disabled) {
            appMenu->appendItem(m_appInfo && m_appInfo->getIdentifyMethod() == "Andriod" ?
                    getMenuItemForceQuitAndroid() : getMenuItemForceQuit());
        }

        if (getAllowedCloseWindows().size() > 0)
            appMenu->appendItem(getMenuItemCloseAll());
    }

    setMenu(appMenu);
}

void Entry::updateIcon()
{
    setPropIcon(getIcon());
}

int Entry::getCurrentMode()
{
    // 只要当前应用是已经驻留的应用，则让其显示为Normal
    if (getIsDocked())
        return ENTRY_NORMAL;

    // 对于未驻留的应用则做如下处理
    if (m_taskmanager->getDisplayMode() == DisplayMode::Efficient) {
        // 高效模式下，只有存在子窗口的，则让其为nornal，没有子窗口的，一般不让其显示
        return hasWindow() ? ENTRY_NORMAL : ENTRY_NONE;
    }
    // 时尚模式下对未驻留应用做如下处理
    // 如果开启了最近打开应用的功能，则显示到最近打开区域（ENTRY_RECENT）
    if (DockSettings::instance()->showRecent())
        return ENTRY_RECENT;

    // 未开启最近使用应用的功能，如果有子窗口，则显示成通用的(ENTRY_NORMAL)，如果没有子窗口，则不显示(ENTRY_NONE)
    return hasWindow() ? ENTRY_NORMAL : ENTRY_NONE;
}

void Entry::updateMode()
{
    int currentMode = getCurrentMode();
    if (m_mode != currentMode) {
        m_mode = currentMode;
        Q_EMIT modeChanged(m_mode);
    }
}

void Entry::forceUpdateIcon()
{
    m_icon = getIcon();
    Q_EMIT iconChanged(m_icon);
}

void Entry::updateIsActive()
{
    bool isActive = false;
    auto activeWin = m_taskmanager->getActiveWindow();
    if (activeWin) {
        // 判断活跃窗口是否属于当前应用
        isActive = m_windowInfoMap.find(activeWin->getXid()) != m_windowInfoMap.end();
    }

    setPropIsActive(isActive);
}

WindowInfoBase *Entry::getWindowInfoByPid(int pid)
{
    for (const auto &windowInfo : m_windowInfoMap) {
        if (windowInfo->getPid() == pid)
            return windowInfo;
    }

    return nullptr;
}

WindowInfoBase *Entry::getWindowInfoByWinId(XWindow windowId)
{
    if (m_windowInfoMap.find(windowId) != m_windowInfoMap.end())
        return m_windowInfoMap[windowId];

    return nullptr;
}

void Entry::setPropIcon(QString value)
{
    if (value != m_icon) {
        m_icon = value;
        Q_EMIT iconChanged(value);
    }
}

void Entry::setPropName(QString value)
{
    if (value != m_name) {
        m_name = value;
        Q_EMIT nameChanged(value);
    }
}

void Entry::setPropIsActive(bool active)
{
    if (m_isActive != active) {
        m_isActive = active;
        Q_EMIT isActiveChanged(active);
    }
}

void Entry::setCurrentWindowInfo(WindowInfoBase *windowInfo)
{
    m_current = windowInfo;
    setPropCurrentWindow(m_current ? m_current->getXid() : 0);
}

void Entry::setPropCurrentWindow(XWindow value)
{
    if (value != m_currentWindow) {
        m_currentWindow = value;
        Q_EMIT currentWindowChanged(value);
    }
}

WindowInfoBase *Entry::getCurrentWindowInfo()
{
    return m_current;
}

/**
 * @brief Entry::findNextLeader
 * @return
 */
WindowInfoBase *Entry::findNextLeader()
{
    auto xids = m_windowInfoMap.keys();
    std::sort(xids.begin(), xids.end());
    XWindow curWinId = m_current->getXid();
    int index = xids.indexOf(curWinId);
    if (index < 0)
        return nullptr;

    // 如果当前窗口是最大， 返回xids[0], 否则返回xids[index + 1]
    int nextIndex = 0;
    if (index < xids.size() - 1)
        nextIndex = index + 1;

    return m_windowInfoMap[xids[nextIndex]];
}

QString Entry::getExec()
{
    if (!m_current)
        return "";

    ProcessInfo *process = m_current->getProcess();
    return process->getExe();
}

QString Entry::getCmdLine()
{
    QString ret;
    if (!m_current) return ret;

    ProcessInfo *process = m_current->getProcess();
    for (auto i : process->getCmdLine()) ret += i + " ";
    return ret;

}

bool Entry::hasWindow()
{
    return m_windowInfoMap.size() > 0;
}

/**
 * @brief Entry::updateExportWindowInfos 同步更新导出窗口信息
 */
void Entry::updateExportWindowInfos()
{
    WindowInfoMap infos;
    for (auto info : m_windowInfoMap) {
        WindowInfo winInfo;
        XWindow xid = info->getXid();
        winInfo.title = info->getTitle();
        winInfo.attention = info->isDemandingAttention();
        winInfo.uuid = info->uuid();
        infos[xid] = winInfo;
    }

    bool changed = true;
    if (infos.size() == m_exportWindowInfos.size()) {
        changed = false;
        for (auto iter = infos.begin(); iter != infos.end(); iter++) {
            XWindow xid = iter.key();
            if (infos[xid].title != m_exportWindowInfos[xid].title ||
                    infos[xid].attention != m_exportWindowInfos[xid].attention ||
                    infos[xid].uuid != m_exportWindowInfos[xid].uuid) {
                changed = true;
                break;
            }
        }
    }

    if (changed) {
        Q_EMIT windowInfosChanged(infos);
    }

    // 更新导出的窗口信息
    m_exportWindowInfos = infos;
}

// 分离窗口， 返回是否需要从任务栏remove
bool Entry::detachWindow(WindowInfoBase *info)
{
    info->setEntry(nullptr);
    XWindow winId = info->getXid();
    if (m_windowInfoMap.contains(winId)) {
        m_windowInfoMap.remove(winId);
        info->deleteLater();
    }

    if (m_windowInfoMap.isEmpty()) {
        if (!m_isDocked) {
            // 既无窗口也非驻留应用，并且不是最近打开，无需在任务栏显示
            return true;
        }

        Q_EMIT windowInfosChanged(WindowInfoMap());
        setCurrentWindowInfo(nullptr);
    } else {
        for (auto window : m_windowInfoMap) {
            if (window) {   // 选择第一个窗口作为当前窗口
                setCurrentWindowInfo(window);
                break;
            }
        }
    }

    updateExportWindowInfos();
    updateIcon();
    updateMenu();

    return false;
}

bool Entry::isShowOnDock() const
{
    // 当前应用显示图标的条件是
    // 如果该图标已经固定在任务栏上，则始终显示
    if (getIsDocked())
        return true;

    // 1.时尚模式下，如果开启了显示最近使用，则不管是否有子窗口，都在任务栏上显示
    // 如果没有开启显示最近使用，则只显示有子窗口的
    if (m_taskmanager->getDisplayMode() == DisplayMode::Fashion)
        return (DockSettings::instance()->showRecent() || m_exportWindowInfos.size() > 0);

    // 2.高效模式下，只有该应用有打开窗口才显示
    return m_exportWindowInfos.size() > 0;
}

bool Entry::attachWindow(WindowInfoBase *info)
{
    XWindow winId = info->getXid();
    qInfo() << "attatchWindow: window id:" << winId;
    info->setEntry(this);

    if (m_windowInfoMap.find(winId) != m_windowInfoMap.end()) {
        qInfo() << "attachWindow: window " << winId << " is already attached";
        return false;
    }

    bool lastShowOnDock = isShowOnDock();
    m_windowInfoMap[winId] = info;
    updateExportWindowInfos();
    updateIsActive();

    if (!m_current) {
        // from no window to has window
        setCurrentWindowInfo(info);
    }

    updateIcon();
    updateMenu();

    if (!lastShowOnDock && isShowOnDock()) {
        // 新打开的窗口始终显示到最后
        Q_EMIT m_taskmanager->entryAdded(this, -1);
    }

    return true;
}

void Entry::launchApp(uint32_t timestamp)
{
    if (m_appInfo)
        m_taskmanager->launchApp(m_appInfo->getFileName(), timestamp, QStringList());
}

bool Entry::containsWindow(XWindow xid)
{
    return m_windowInfoMap.find(xid) != m_windowInfoMap.end();
}

// 处理菜单项
void Entry::handleMenuItem(uint32_t timestamp, QString itemId)
{
    m_appMenu->handleAction(timestamp, itemId);
}

// 处理拖拽事件
void Entry::handleDragDrop(uint32_t timestamp, QStringList files)
{
    m_taskmanager->launchApp(m_appInfo->getFileName(), timestamp, files);
}

// 驻留
void Entry::requestDock(bool dockToEnd)
{
    if (m_taskmanager->dockEntry(this, dockToEnd)) {
        m_taskmanager->saveDockedApps();
    }
}

// 取消驻留
void Entry::requestUndock(bool dockToEnd)
{
    m_taskmanager->undockEntry(this, dockToEnd);
}

void Entry::newInstance(uint32_t timestamp)
{
    QStringList files;
    m_taskmanager->launchApp(m_appInfo->getFileName(), timestamp, files);
}

// 检查应用窗口分离、合并状态
void Entry::check()
{
    QList<WindowInfoBase *> windows = m_windowInfoMap.values();
    for (WindowInfoBase *window : windows) {
        m_taskmanager->attachOrDetachWindow(window);
    }
}

// 强制退出
void Entry::forceQuit()
{
    QMap<int, QVector<WindowInfoBase*>> pidWinInfoMap;
    QList<WindowInfoBase *> windows = m_windowInfoMap.values();
    for (WindowInfoBase *window : windows) {
        int pid = window->getPid();
        if (pid != 0) {
            pidWinInfoMap[pid].push_back(window);
        } else {
            window->killClient();
        }
    }

    for (auto iter = pidWinInfoMap.begin(); iter != pidWinInfoMap.end(); iter++) {
        if (!killProcess(iter.key())) {         // kill pid
            for (auto &info : iter.value()) {   // kill window
                info->killClient();
            }
        }
    }
    // 所有的窗口已经退出后，清空m_windowInfoMap内容
    m_windowInfoMap.clear();
    // 退出所有的进程后，及时更新当前剩余的窗口数量
    updateExportWindowInfos();
    m_taskmanager->removeEntryFromDock(this);
}

void Entry::presentWindows()
{
    QList<uint> windows = m_windowInfoMap.keys();
    m_taskmanager->presentWindows(windows);
}

/**
 * @brief Entry::active 激活窗口
 * @param timestamp
 */
void Entry::active(uint32_t timestamp)
{
    if (m_taskmanager->getHideMode() == HideMode::SmartHide) {
        m_taskmanager->setPropHideState(HideState::Show);
        m_taskmanager->updateHideState(false);
    }

    // 无窗口则直接启动
    if (!hasWindow()) {
        launchApp(timestamp);
        return;
    }

    if (!m_current) {
        qWarning() << "active: current window is nullptr";
        return;
    }

    WindowInfoBase *winInfo = m_current;
    if (m_taskmanager->isWaylandEnv()) {
        // wayland环境
        if (!m_taskmanager->isActiveWindow(winInfo)) {
            winInfo->activate();
        } else {
            bool showing = m_taskmanager->isShowingDesktop();
            if (showing || winInfo->isMinimized()) {
                winInfo->activate();
            } else if (m_windowInfoMap.size() == 1) {
                winInfo->minimize();
            } else {
                WindowInfoBase *nextWin = findNextLeader();
                if (nextWin) {
                    nextWin->activate();
                }
            }
        }
    } else {
        // X11环境
        XWindow xid = winInfo->getXid();
        WindowInfoBase *activeWin = m_taskmanager->getActiveWindow();
        if (activeWin && xid != activeWin->getXid()) {
            m_taskmanager->doActiveWindow(xid);
        } else {
            bool found = false;
            XWindow hiddenAtom = XCB->getAtom("_NET_WM_STATE_HIDDEN");
            for (auto state : XCB->getWMState(xid)) {
                if (hiddenAtom == state) {
                    found = true;
                    break;
                }
            }

            if (found) {
                // 激活隐藏窗口
                m_taskmanager->doActiveWindow(xid);
            } else if (m_windowInfoMap.size() == 1) {
                // 窗口图标化
                XCB->minimizeWindow(xid);
            } else if (m_taskmanager->getActiveWindow() && m_taskmanager->getActiveWindow()->getXid() == xid) {
                WindowInfoBase *nextWin = findNextLeader();
                if (nextWin) {
                    nextWin->activate();
                }
            }
        }
    }
}

void Entry::activeWindow(quint32 winId)
{
    if (m_taskmanager->isWaylandEnv()) {
        if (!m_windowInfoMap.contains(winId))
            return;

        WindowInfoBase *winInfo = m_windowInfoMap[winId];
        if (m_taskmanager->isActiveWindow(winInfo)) {
            bool showing = m_taskmanager->isShowingDesktop();
            if (showing || winInfo->isMinimized()) {
                winInfo->activate();
            } else if (m_windowInfoMap.size() == 1) {
                winInfo->minimize();
            } else {
                WindowInfoBase *nextWin = findNextLeader();
                if (nextWin) {
                    nextWin->activate();
                }
            }
        } else {
            winInfo->activate();
        }
    } else {
        m_taskmanager->doActiveWindow(winId);
    }
}

int Entry::mode()
{
    return m_mode;
}

XWindow Entry::getCurrentWindow()
{
    return m_currentWindow;
}

QString Entry::getDesktopFile()
{
    return m_desktopFile;
}

bool Entry::getIsActive() const
{
    return m_isActive;
}

QString Entry::getMenu() const
{
    return m_appMenu->getMenuJsonStr();
}

QVector<XWindow> Entry::getAllowedClosedWindowIds()
{
    QVector<XWindow> ret;
    for (auto iter = m_windowInfoMap.begin(); iter != m_windowInfoMap.end(); iter++) {
        WindowInfoBase *info = iter.value();
        if (info && info->allowClose())
            ret.push_back(iter.key());
    }

    return ret;
}

WindowInfoMap Entry::getExportWindowInfos()
{
    return m_exportWindowInfos;
}

QVector<WindowInfoBase *> Entry::getAllowedCloseWindows()
{
    QVector<WindowInfoBase *> ret;
    for (auto iter = m_windowInfoMap.begin(); iter != m_windowInfoMap.end(); iter++) {
        WindowInfoBase *info = iter.value();
        if (info && info->allowClose()) {
            ret.push_back(info);
        }
    }

    return ret;
}

QVector<AppMenuItem> Entry::getMenuItemDesktopActions()
{
    QVector<AppMenuItem> ret;
    if (!m_appInfo) {
        return ret;
    }

    for (auto action : m_appInfo->getActions()) {
        AppMenuAction fn = [=](uint32_t timestamp) {
            qInfo() << "do MenuItem: " << action.name;
            m_taskmanager->launchAppAction(m_appInfo->getFileName(), action.section, timestamp);
        };

        AppMenuItem item;
        item.text = action.name;
        item.action = fn;
        item.isActive = true;
        ret.push_back(item);
    }

    return ret;
}

AppMenuItem Entry::getMenuItemLaunch()
{
    QString itemName;
    if (hasWindow()) {
        itemName = getName();
    } else {
        itemName = tr("Open");
    }

    AppMenuAction fn = [this](uint32_t timestamp) {
        qInfo() << "do MenuItem: Open";
        this->launchApp(timestamp);
    };

    AppMenuItem item;
    item.text = itemName;
    item.action = fn;
    item.isActive = true;
    return item;
}

AppMenuItem Entry::getMenuItemCloseAll()
{
    AppMenuAction fn = [this](uint32_t timestamp) {
        qInfo() << "do MenuItem: Close All";
        auto winInfos = getAllowedCloseWindows();

        // 根据创建时间从大到小排序， 方便后续关闭窗口
        for (int i = 0; i < winInfos.size() - 1; i++) {
            for (int j = i + 1; j < winInfos.size(); j++) {
                if (winInfos[i]->getCreatedTime() < winInfos[j]->getCreatedTime()) {
                    auto info = winInfos[i];
                    winInfos[i] = winInfos[j];
                    winInfos[j] = info;
                }
            }
        }

        for (auto info : winInfos) {
            qInfo() << "close WindowId " << info->getXid();
            info->close(timestamp);
        }

        // 关闭窗口后，主动刷新事件
        XCB->flush();
    };

    AppMenuItem item;
    item.text = tr("Close All");
    item.action = fn;
    item.isActive = true;
    return item;
}

AppMenuItem Entry::getMenuItemForceQuit()
{
    bool active = m_taskmanager->getForceQuitAppStatus() != ForceQuitAppMode::Deactivated;
    AppMenuAction fn = [this](uint32_t) {
        qInfo() << "do MenuItem: Force Quit";
        forceQuit();
    };

    AppMenuItem item;
    item.text = tr("Force Quit");
    item.action = fn;
    item.isActive = active;
    return item;
}

//dock栏上Android程序的Force Quit功能
AppMenuItem Entry::getMenuItemForceQuitAndroid()
{
    bool active = m_taskmanager->getForceQuitAppStatus() != ForceQuitAppMode::Deactivated;
    auto allowedCloseWindows = getAllowedCloseWindows();
    AppMenuAction fn = [](uint32_t){};
    if (allowedCloseWindows.size() > 0) {
        qInfo() << "do MenuItem: Force Quit";
        AppMenuAction fn = [&](uint32_t timestamp) {
            for (auto info : allowedCloseWindows) {
                info->close(timestamp);
            }
        };
    }

    AppMenuItem item;
    item.text = tr("Force Quit");
    item.action = fn;
    item.isActive = active;
    return item;
}

AppMenuItem Entry::getMenuItemDock()
{
    AppMenuItem item;
    item.text = tr("Dock");
    item.action = [this](uint32_t) {
        qInfo() << "do MenuItem: Dock";
        requestDock(true);
    };

    item.isActive = true;
    return item;
}

AppMenuItem Entry::getMenuItemUndock()
{
    AppMenuItem item;
    item.text = tr("Undock");
    item.action = [this](uint32_t) {
        qInfo() << "do MenuItem: Undock";
        requestUndock(true);
    };

    item.isActive = true;
    return item;
}

AppMenuItem Entry::getMenuItemAllWindows()
{
    AppMenuItem item;
    item.text = tr("All Windows");
    item.action = [this](uint32_t) {
        qInfo() << "do MenuItem: All Windows";
        presentWindows();
    };

    item.isActive = true;
    item.hint = 1;
    return item;
}

bool Entry::killProcess(int pid)
{
    return  !kill(pid, SIGTERM);
}

bool Entry::setPropDesktopFile(QString value)
{
    if (value != m_desktopFile) {
        m_desktopFile = value;
        Q_EMIT desktopFileChanged(value);
        return true;
    }

    return false;
}
