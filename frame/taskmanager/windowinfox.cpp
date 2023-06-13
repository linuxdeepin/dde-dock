// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "windowinfox.h"
#include "appinfo.h"
#include "xcbutils.h"
#include "common.h"
#include "processinfo.h"

#include <QDebug>
#include <QCryptographicHash>
#include <QTimer>
#include <QImage>
#include <QIcon>
#include <QBuffer>

#include <X11/Xlib.h>
#include <algorithm>
#include <qobject.h>
#include <string>

#define XCB XCBUtils::instance()

WindowInfoX::WindowInfoX(XWindow _xid, QObject *parent)
 : WindowInfoBase (parent)
 , m_x(0)
 , m_y(0)
 , m_width(0)
 , m_height(0)
 , m_hasWMTransientFor(false)
 , m_hasXEmbedInfo(false)
 , m_updateCalled(false)
{
    xid = _xid;
    m_createdTime = std::chrono::duration_cast<std::chrono::nanoseconds>(std::chrono::system_clock::now().time_since_epoch()).count(); // 获取当前时间，精确到纳秒
}

WindowInfoX::~WindowInfoX()
{

}

bool WindowInfoX::shouldSkip()
{
    qInfo() << "window " << xid << " shouldSkip?";
    if (!m_updateCalled) {
        update();
        m_updateCalled = true;
    }

    if (hasWmStateSkipTaskBar() || isValidModal() || shouldSkipWithWMClass())
        return true;

    for (auto atom : m_wmWindowType) {
        if (atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_DIALOG") && !isActionMinimizeAllowed())
            return true;

        if (atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_UTILITY")
        || atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_COMBO")
        || atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_DESKTOP") // 桌面属性窗口
        || atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_DND")
        || atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_DOCK")    // 任务栏属性窗口
        || atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_DROPDOWN_MENU")
        || atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_MENU")
        || atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_NOTIFICATION")
        || atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_POPUP_MENU")
        || atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_SPLASH")
        || atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_TOOLBAR")
        || atom == XCB->getAtom("_NET_WM_WINDOW_TYPE_TOOLTIP"))
            return true; 
    }

    return false;
}

QString WindowInfoX::getIcon()
{
    if (icon.isEmpty())
        icon = getIconFromWindow();

    return icon;
}

void WindowInfoX::activate()
{
    XCB->changeActiveWindow(xid);
    QTimer::singleShot(50, [&] {
        XCB->restackWindow(xid);
    });
}

void WindowInfoX::minimize()
{
    XCB->minimizeWindow(xid);
}

bool WindowInfoX::isMinimized()
{
    return containAtom(m_wmState, XCB->getAtom("_NET_WM_STATE_HIDDEN"));
}

int64_t WindowInfoX::getCreatedTime()
{
    return m_createdTime;
}

QString WindowInfoX::getWindowType()
{
    return "X11";
}

bool WindowInfoX::allowClose()
{
    // 允许关闭的条件：
    // 1. 不设置 functions 字段，即MotifHintFunctions 标志位；
    // 2. 或者设置了 functions 字段并且 设置了 MotifFunctionAll 标志位；
    // 3. 或者设置了 functions 字段并且 设置了 MotifFunctionClose 标志位。
    // 相关定义在 motif-2.3.8/lib/Xm/MwmUtil.h 。
    if ((m_motifWmHints.flags & MotifHintFunctions) == 0
     || (m_motifWmHints.functions & MotifFunctionAll) != 0
     || (m_motifWmHints.functions & MotifFunctionClose) != 0)
        return true;

    for (auto action : m_wmAllowedActions) {
        if (action == XCB->getAtom("_NET_WM_ACTION_CLOSE")) {
            return true;
        }
    }

    return false;
}

QString WindowInfoX::getDisplayName()
{
    XWindow winId = xid;
    //QString role = wmRole;
    QString className(m_wmClass.className.c_str());
    QString instance;
    if (m_wmClass.instanceName.size() > 0) {
        int pos = QString(m_wmClass.instanceName.c_str()).lastIndexOf('/');
        if (pos != -1)
            instance.remove(0, pos + 1);
    }
    qInfo() << "getDisplayName class:" << className << " ,instance:" << instance;

    //if (!role.isEmpty() && !className.isEmpty())
    //    return className + " " + role;

    if (!className.isEmpty())
        return className;

    if (!instance.isEmpty())
        return instance;


    QString _wmName = m_wmName;
    if (!_wmName.isEmpty()) {
        int pos = _wmName.lastIndexOf('-');
        if (pos != -1 && !_wmName.startsWith("-")) {
            _wmName.truncate(pos);
            return _wmName;
        }
    }

    if (m_processInfo) {
        QString exe {m_processInfo->getEnv("exe")};
        if (!exe.isEmpty())
            return exe;
    }

    return QString("window:%1").arg(winId);
}

void WindowInfoX::killClient()
{
    XCB->killClientChecked(xid);
}

QString WindowInfoX::uuid()
{
    return QString();
}

QString WindowInfoX::getGtkAppId()
{
    return m_gtkAppId;
}

QString WindowInfoX::getFlatpakAppId()
{
    return m_flatpakAppId;
}

QString WindowInfoX::getWmRole()
{
    return m_wmRole;
}

WMClass WindowInfoX::getWMClass()
{
    return m_wmClass;
}

QString WindowInfoX::getWMName()
{
    return m_wmName;
}

ConfigureEvent *WindowInfoX::getLastConfigureEvent()
{
    return m_lastConfigureNotifyEvent;
}

void WindowInfoX::setLastConfigureEvent(ConfigureEvent *event)
{
    m_lastConfigureNotifyEvent = event;
}

bool WindowInfoX::isGeometryChanged(int _x, int _y, int _width, int _height)
{
    return !(_x == m_x && _y == m_y && _width == m_width && _height == m_height);
}

void WindowInfoX::setGtkAppId(QString _gtkAppId)
{
    m_gtkAppId = _gtkAppId;
}

void WindowInfoX::updateMotifWmHints()
{
    // get from XCB
    m_motifWmHints = XCB->getWindowMotifWMHints(xid);
}

// XEmbed info
// 一般 tray icon 会带有 _XEMBED_INFO 属性
void WindowInfoX::updateHasXEmbedInfo()
{
    m_hasXEmbedInfo = XCB->hasXEmbedInfo(xid);
}

/**
 * @brief WindowInfoX::genInnerId 生成innerId
 * @param winInfo
 * @return
 */
QString WindowInfoX::genInnerId(WindowInfoX *winInfo)
{
    XWindow winId = winInfo->getXid();
    QString wmClassName, wmInstance;
    WMClass wmClass = winInfo->getWMClass();
    if (wmClass.className.size() > 0)
        wmClassName = wmClass.className.c_str();

    if (wmClass.instanceName.size() > 0) {
        QString instanceName(wmClass.instanceName.c_str());
        instanceName.remove(0, instanceName.lastIndexOf('/') + 1);
        wmInstance = instanceName;
    }

    QString exe, args;
    if (winInfo->getProcess()) {
        exe = winInfo->getProcess()->getExe();
        for (auto arg : winInfo->getProcess()->getArgs()) {
            if (arg.contains("/") || arg == "." || arg == "..") {
                args += "%F ";
            } else {
                args += arg + " ";
            }
        }

        if (args.size() > 0)
            args.remove(args.size() - 2, 1);
    }

    bool hasPid = winInfo->getPid() != 0;
    QString str;
    // NOTE: 不要使用 wmRole，有些程序总会改变这个值比如 GVim
    if (wmInstance.isEmpty() && wmClassName.isEmpty() && exe.isEmpty() && winInfo->getGtkAppId().isEmpty()) {
        if (!winInfo->getWMName().isEmpty())
            str = QString("wmName:%1").arg(winInfo->getWMName());
        else
            str = QString("windowId:%1").arg(winInfo->getXid());
    } else {
        str = QString("wmInstance:%1,wmClass:%2,exe:%3,args:%4,hasPid:%5,gtkAppId:%6").arg(wmInstance).arg(wmClassName).arg(exe).arg(args).arg(hasPid).arg(winInfo->getGtkAppId());
    }

    QByteArray encryText = QCryptographicHash::hash(str.toLatin1(), QCryptographicHash::Md5);
    QString innerId = windowHashPrefix + encryText.toHex();
    qInfo() << "genInnerId window " << winId << " innerId :" << innerId;
    return innerId;
}

// 更新窗口类型
void WindowInfoX::updateWmWindowType()
{
    m_wmWindowType.clear();
    for (auto ty : XCB->getWMWindoType(xid)) {
        m_wmWindowType.push_back(ty);
    }
}

// 更新窗口许可动作
void WindowInfoX::updateWmAllowedActions()
{
    m_wmAllowedActions.clear();
    for (auto action : XCB->getWMAllowedActions(xid)) {
        m_wmAllowedActions.push_back(action);
    }
}

void WindowInfoX::updateWmState()
{
    m_wmState.clear();
    for (auto a : XCB->getWMState(xid)) {
        m_wmState.push_back(a);
    }
}

void WindowInfoX::updateWmClass()
{
    m_wmClass = XCB->getWMClass(xid);
}

void WindowInfoX::updateWmName()
{
    auto name = XCB->getWMName(xid);
    if (!name.empty())
        m_wmName = name.c_str();

    title = getTitle();
}

void WindowInfoX::updateIcon()
{
    icon = getIconFromWindow();
}

void WindowInfoX::updateHasWmTransientFor()
{
    if (XCB->getWMTransientFor(xid) == 1)
        m_hasWMTransientFor = true;
}

/**
 * @brief WindowInfoX::update 更新窗口信息（在识别窗口时执行一次）
 */
void WindowInfoX::update()
{
    updateWmClass();
    updateWmState();
    updateWmWindowType();
    updateWmAllowedActions();
    updateHasWmTransientFor();
    updateProcessInfo();
    updateWmName();
    innerId = genInnerId(this);
}

QString WindowInfoX::getIconFromWindow()
{
    WMIcon icon = XCB->getWMIcon(xid);

    // invalid icon
    if (icon.width == 0) {
        return QString();
    }

    QImage img = QImage((uchar *)icon.data.data(), icon.width, icon.width, QImage::Format_ARGB32);
    QBuffer buffer;
    buffer.open(QIODevice::WriteOnly);
    img.scaled(48, 48, Qt::KeepAspectRatio, Qt::SmoothTransformation);
    img.save(&buffer, "PNG");

    // convert to base64
    QString encode = buffer.data().toBase64();
    QString iconPath = QString("%1,%2").arg("data:image/png:base64").arg(encode);
    buffer.close();

    return iconPath;
}

bool WindowInfoX::isActionMinimizeAllowed()
{
    return containAtom(m_wmAllowedActions, XCB->getAtom("_NET_WM_ACTION_MINIMIZE"));
}

bool WindowInfoX::hasWmStateDemandsAttention()
{
    return containAtom(m_wmState, XCB->getAtom("_NET_WM_STATE_DEMANDS_ATTENTION"));
}

bool WindowInfoX::hasWmStateSkipTaskBar()
{
    return containAtom(m_wmState, XCB->getAtom("_NET_WM_STATE_SKIP_TASKBAR"));
}

bool WindowInfoX::hasWmStateModal()
{
    return containAtom(m_wmState, XCB->getAtom("_NET_WM_STATE_MODAL"));
}

bool WindowInfoX::isValidModal()
{
    return hasWmStateModal() && hasWmStateModal();
}

// 通过WMClass判断是否需要隐藏此窗口
bool WindowInfoX::shouldSkipWithWMClass()
{
    bool ret = false;
    if (m_wmClass.instanceName == "explorer.exe" && m_wmClass.className == "Wine")
        ret = true;
    else if (m_wmClass.className == "dde-launcher" ||
        m_wmClass.className == "dde-dock" ||
        m_wmClass.className == "dde-lock") {
        ret = true;
    }

    return ret;
}

void WindowInfoX::updateProcessInfo()
{
    XWindow winId = xid;
    pid = XCB->getWMPid(winId);
    qInfo() << "updateProcessInfo: pid=" << pid;
    m_processInfo.reset(new ProcessInfo(pid));
    if (!m_processInfo->isValid()) {
        // try WM_COMMAND
        auto wmComand = XCB->getWMCommand(winId);
        if (wmComand.size() > 0) {
            QStringList cmds;
            std::transform(wmComand.begin(), wmComand.end(), std::back_inserter(cmds), [=] (std::string cmd){ return QString::fromStdString(cmd);});
            m_processInfo.reset(new ProcessInfo(cmds));
        }
    }

    qInfo() << "updateProcessInfo: pid is " << pid;
}

bool WindowInfoX::getUpdateCalled()
{
    return m_updateCalled;
}

void WindowInfoX::setInnerId(QString _innerId)
{
    innerId = _innerId;
}

QString WindowInfoX::getTitle()
{
    QString name = m_wmName;
    if (name.isEmpty())
        name = getDisplayName();

    return name;
}

bool WindowInfoX::isDemandingAttention()
{
    return hasWmStateDemandsAttention();
}

void WindowInfoX::close(uint32_t timestamp)
{
    XCB->requestCloseWindow(xid, timestamp);
}
