// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "dbushandler.h"
#include "dbusutil.h"
#include "taskmanager.h"
#include "common.h"
#include "entry.h"
#include "windowinfok.h"

#include <DDBusSender>

#include <QtConcurrent>

DBusHandler::DBusHandler(TaskManager *taskmanager, QObject *parent)
    : QObject(parent)
    , m_taskmanager(taskmanager)
    , m_wm(new com::deepin::wm("com.deepin.wm", "/com/deepin/wm", QDBusConnection::sessionBus(), this))
    , m_wmSwitcher(new org::deepin::dde::WMSwitcher1("org.deepin.dde.WMSwitcher1", "/org/deepin/dde/WMSwitcher1", QDBusConnection::sessionBus(), this))
    , m_kwaylandManager(nullptr)
    , m_xEventMonitor(nullptr)
    , m_launcher(new org::deepin::dde::Launcher1(launcherService, launcherPath, QDBusConnection::sessionBus(), this))
{
    connect(m_wmSwitcher, &org::deepin::dde::WMSwitcher1::WMChanged, this, [&](QString name) {m_taskmanager->setWMName(name);});
    if (!isWaylandSession()) {
        m_xEventMonitor = new org::deepin::dde::XEventMonitor1("org.deepin.dde.XEventMonitor1", "/org/deepin/dde/XEventMonitor1", QDBusConnection::sessionBus(), this);
        m_activeWindowMonitorKey = m_xEventMonitor->RegisterFullScreen();
        connect(m_xEventMonitor, &org::deepin::dde::XEventMonitor1::ButtonRelease, this, &DBusHandler::onActiveWindowButtonRelease);
    }

    connect(m_launcher,static_cast<void (org::deepin::dde::Launcher1::*)(bool)>(&org::deepin::dde::Launcher1::VisibleChanged), this, [=] (const bool visible){
        m_taskmanager->setDdeLauncherVisible(visible);
        m_taskmanager->updateHideState(true);
    });
}

void DBusHandler::listenWaylandWMSignals()
{
    m_kwaylandManager = new org::deepin::dde::kwayland1::WindowManager("org.deepin.dde.KWayland1", "/org/deepin/dde/KWayland1/WindowManager", QDBusConnection::sessionBus(), this);
    connect(m_kwaylandManager, &org::deepin::dde::kwayland1::WindowManager::ActiveWindowChanged, this, &DBusHandler::handleWlActiveWindowChange);
    connect(m_kwaylandManager, &org::deepin::dde::kwayland1::WindowManager::WindowCreated, this, [&] (const QString &ObjPath) {
        m_taskmanager->registerWindowWayland(ObjPath);
    });
    connect(m_kwaylandManager, &org::deepin::dde::kwayland1::WindowManager::WindowRemove, this, [&] (const QString &ObjPath) {
        m_taskmanager->unRegisterWindowWayland(ObjPath);
    });
}

void DBusHandler::loadClientList()
{
    if (!m_kwaylandManager)
        return;

    QDBusPendingReply<QVariantList> windowList = m_kwaylandManager->Windows();
    QVariantList windows = windowList.value();
    for (QVariant windowPath : windows)
        m_taskmanager->registerWindowWayland(windowPath.toString());
}

QString DBusHandler::getCurrentWM()
{
    return m_wmSwitcher->CurrentWM().value();
}

void DBusHandler::launchApp(QString desktopFile, uint32_t timestamp, QStringList files)
{
    if (newStartManagerAvaliable()) {
        auto objPath = desktopEscapeToObjectPath(desktopFile);
        launchAppUsingApplicationManager1(QString{DDEApplicationManager1ObjectPath} + '/' + objPath, timestamp, files);
    } else {
        launchAppUsingApplication1Manager(desktopFile, timestamp, files);
    }
}

void DBusHandler::launchAppAction(QString desktopFile, QString action, uint32_t timestamp)
{
    if (newStartManagerAvaliable()) {
        auto objPath = desktopEscapeToObjectPath(desktopFile);
        launchAppActionUsingApplicationManager1(QString{DDEApplicationManager1ObjectPath} + '/' + objPath, action, timestamp);
    } else {
        launchAppActionUsingApplication1Manager(desktopFile, action, timestamp);
    }
}

void DBusHandler::launchAppUsingApplication1Manager(QString desktopFile, uint32_t timestamp, QStringList files)
{
    QDBusInterface interface = QDBusInterface("org.deepin.dde.Application1.Manager", "/org/deepin/dde/Application1/Manager", "org.deepin.dde.Application1.Manager");
    interface.call("LaunchApp", desktopFile, timestamp, files);
}

void DBusHandler::launchAppActionUsingApplication1Manager(QString desktopFile, QString action, uint32_t timestamp)
{
    QDBusInterface interface = QDBusInterface("org.deepin.dde.Application1.Manager", "/org/deepin/dde/Application1/Manager", "org.deepin.dde.Application1.Manager");
    interface.call("LaunchAppAction", desktopFile, action, timestamp);
}

// 新AM启动接口
void DBusHandler::launchAppUsingApplicationManager1(QString dbusObjectPath, uint32_t timestamp, QStringList files)
{
    QDBusInterface interface = QDBusInterface(ApplicationManager1DBusName, dbusObjectPath, "org.desktopspec.ApplicationManager1.Application");
    interface.call("Launch", "", QStringList(), QMap<QString, QVariant>());
}

void DBusHandler::launchAppActionUsingApplicationManager1(QString dbusObjectPath, QString action, uint32_t timestamp)
{
    action = action.right(action.size() - strlen(DesktopFileActionKey));
    QDBusInterface interface = QDBusInterface(ApplicationManager1DBusName, dbusObjectPath, "org.desktopspec.ApplicationManager1.Application");
    interface.call("Launch", action, QStringList(), QMap<QString, QVariant>());
}

void DBusHandler::markAppLaunched(const QString &filePath)
{
    QDBusInterface interface = QDBusInterface("org.deepin.dde.AlRecorder1", "/org/deepin/dde/AlRecorder1", "org.deepin.dde.AlRecorder1");
    interface.call("MarkLaunched", filePath);
}

bool DBusHandler::wlShowingDesktop()
{
    bool ret = false;
    if (m_kwaylandManager)
        ret = m_kwaylandManager->IsShowingDesktop().value();

    return ret;
}

uint DBusHandler::wlActiveWindow()
{
    uint ret = 0;
    if (m_kwaylandManager)
        ret = m_kwaylandManager->ActiveWindow().value();

    return ret;
}

void DBusHandler::handleWlActiveWindowChange()
{
    uint activeWinInternalId = wlActiveWindow();
    if (activeWinInternalId == 0)
        return;

    WindowInfoK *info = m_taskmanager->handleActiveWindowChangedK(activeWinInternalId);
    if (info && info->getXid() != 0) {
        m_taskmanager->handleActiveWindowChanged(info);
    } else {
        m_taskmanager->updateHideState(false);
    }
}

void DBusHandler::onActiveWindowButtonRelease(int type, int x, int y, const QString &key)
{
    // 当鼠标松开区域事件的时候，取消注册，同时调用激活窗口的方法来触发智能隐藏的相关信号
    if (key != m_activeWindowMonitorKey)
        return;

    uint activeWinInternalId = wlActiveWindow();
    if (activeWinInternalId == 0)
        return;

    WindowInfoK *info = m_taskmanager->handleActiveWindowChangedK(activeWinInternalId);
    if (!info)
        return;

    // 如果是在当前激活窗口区域内释放的，则触发检测智能隐藏的方法
    DockRect dockRect = info->getGeometry();
    if (dockRect.x <= x && x <= int(dockRect.x + dockRect.w) && dockRect.y <= y && y <= int(dockRect.y + dockRect.h)) {
        // 取消智能隐藏
        m_taskmanager->updateHideState(false);
    }
}

void DBusHandler::listenKWindowSignals(WindowInfoK *windowInfo)
{
    PlasmaWindow *window = windowInfo->getPlasmaWindow();
    if (!window)
        return;

    connect(window, &PlasmaWindow::TitleChanged, this, [=] {
        windowInfo->updateTitle();
        auto entry = m_taskmanager->getEntryByWindowId(windowInfo->getXid());
        if (entry && entry->getCurrentWindowInfo() == windowInfo)
            entry->updateName();
    });
    connect(window, &PlasmaWindow::IconChanged, this, [=] {
        windowInfo->updateIcon();
        auto entry = m_taskmanager->getEntryByWindowId(windowInfo->getXid());
        if (!entry) return;

        entry->updateIcon();
    });

    // DemandingAttention changed
    connect(window, &PlasmaWindow::DemandsAttentionChanged, this, [=] {
        windowInfo->updateDemandingAttention();
        auto entry = m_taskmanager->getEntryByWindowId(windowInfo->getXid());
        if (!entry) return;

        entry->updateExportWindowInfos();
    });

    // Geometry changed
    connect(window, &PlasmaWindow::GeometryChanged, this, [=] {
        if (!windowInfo->updateGeometry()) return;

        m_taskmanager->handleWindowGeometryChanged();
    });
}

PlasmaWindow *DBusHandler::createPlasmaWindow(QString objPath)
{
    return new PlasmaWindow("org.deepin.dde.KWayland1", objPath, QDBusConnection::sessionBus(), this);
}

/**
 * @brief DBusHandler::removePlasmaWindowHandler 取消关联信号 TODO
 * @param window
 */
void DBusHandler::removePlasmaWindowHandler(PlasmaWindow *window)
{

}

void DBusHandler::presentWindows(QList<uint> windows)
{
    m_wm->PresentWindows(windows);
}

void DBusHandler::previewWindow(uint xid)
{
    m_wm->PreviewWindow(xid);
}

void DBusHandler::cancelPreviewWindow()
{
    m_wm->CancelPreviewWindow();
}

// TODO: 待优化点， 查看Bamf根据windowId获取对应应用desktopFile路径实现方式, 移除bamf依赖
QString DBusHandler::getDesktopFromWindowByBamf(XWindow windowId)
{
    QDBusInterface interface0 = QDBusInterface("org.ayatana.bamf", "/org/ayatana/bamf/matcher", "org.ayatana.bamf.matcher");
    QDBusReply<QString> replyApplication = interface0.call("ApplicationForXid", windowId);
    QString appObjPath = replyApplication.value();
    if (!replyApplication.isValid() || appObjPath.isEmpty())
        return "";


    QDBusInterface interface = QDBusInterface("org.ayatana.bamf", appObjPath, "org.ayatana.bamf.application");
    QDBusReply<QString> replyDesktopFile = interface.call("DesktopFile");

    if (replyDesktopFile.isValid())
        return replyDesktopFile.value();


    return "";
}

// 新的AM调用
QString DBusHandler::desktopEscapeToObjectPath(QString desktopFilePath)
{
    // to desktop id
    QString objectPath;
    decltype(auto) desktopSuffix = ".desktop";
    auto tmp = desktopFilePath.chopped(sizeof(desktopSuffix) - 1);
    auto components = tmp.split(QDir::separator());
    auto it = std::find(components.cbegin(), components.cend(), "applications");
    if (it == components.cend()) return "_";
    QString FileId;
    ++it;
    while (it != components.cend()) {
        FileId += (*(it++) + "-");
    }
    objectPath = FileId.chopped(1);

    if (objectPath.isEmpty()) {
        return "_";
    }

    // desktop id to objectPath
    QRegularExpression re{R"([^a-zA-Z0-9])"};
    auto matcher = re.globalMatch(objectPath);
    while (matcher.hasNext()) {
        auto replaceList = matcher.next().capturedTexts();
        replaceList.removeDuplicates();
        for (const auto &c : replaceList) {
            auto hexStr = QString::number(static_cast<uint>(c.front().toLatin1()), 16);
            objectPath.replace(c, QString{R"(_%1)"}.arg(hexStr));
        }
    }

    return objectPath;
}

bool DBusHandler::newStartManagerAvaliable()
{
    static bool isAvaiable = false;
    std::call_once(m_isNewStartManagerAvaliableInited, [=](){
        auto services = QDBusConnection::sessionBus().interface()->registeredServiceNames().value();
        isAvaiable = std::any_of(services.begin(), services.end(), [=](const QString &name){
            return name == ApplicationManager1DBusName;
        });
    });

    return isAvaiable;
}

void DBusHandler::sendFailedDockNotification(const QString &appName)
{
    QtConcurrent::run(QThreadPool::globalInstance(), [ = ] {
        DDBusSender()
                .service(notificationService)
                .path(notificationPath)
                .interface(notificationInterface)
                .method(QString("Notify"))
                .arg(QString("dde-control-center"))                                            // appname
                .arg(static_cast<uint>(0))                                                     // id
                .arg(QString("preferences-system"))                                            // icon
                .arg(QString(tr("failed to dock ") + appName))                                          // summary
                .arg(QString(tr("Unrecognized application, unable to dock"))) // content
                .arg(QStringList())                                                             // actions
                .arg(QVariantMap())                                                                    // hints
                .arg(5000)                                                                    // timeout
                .call();
    });
    qInfo() << "unsupported app: " + appName;
}
