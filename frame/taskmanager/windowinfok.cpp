// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "windowinfok.h"
#include "entry.h"
#include "processinfo.h"
#include "appinfo.h"

#include <chrono>
#include <qobject.h>

WindowInfoK::WindowInfoK(PlasmaWindow *window, XWindow _xid, QObject *parent)
 : WindowInfoBase (parent)
 , m_updateCalled(false)
 , m_internalId(0)
 , m_demaningAttention(false)
 , m_closeable(true)
 , m_plasmaWindow(window)
{
    xid = _xid;
    m_createdTime = std::chrono::duration_cast<std::chrono::nanoseconds>(std::chrono::system_clock::now().time_since_epoch()).count(); // 获取当前时间，精确到纳秒
}

WindowInfoK::~WindowInfoK()
{
}

bool WindowInfoK::shouldSkip()
{
    if (!m_updateCalled) {
        update();
        m_updateCalled = true;
    }

    bool skip = m_plasmaWindow->SkipTaskbar();

    // 添加窗口能否最小化判断, 如果窗口不能最小化则隐藏任务栏图标
    bool canMinimize = false;
    canMinimize = m_plasmaWindow->IsMinimizeable();
    if (!canMinimize)
        skip = true;

    if (skip) {
        // 白名单, 过滤类似“欢迎应用”， 没有最小化窗口但是需要在任务栏显示图标
        QStringList list { "dde-introduction"};
        if (list.indexOf(m_appId) != -1)
            skip = false;
    }

    return skip;
}

QString WindowInfoK::getIcon()
{
    return icon;
}

QString WindowInfoK::getTitle()
{
    return title;
}

bool WindowInfoK::isDemandingAttention()
{
    return m_demaningAttention;
}

bool WindowInfoK::allowClose()
{
    return m_closeable;
}

void WindowInfoK::close(uint32_t timestamp)
{
    m_plasmaWindow->RequestClose();
}

QString WindowInfoK::getAppId()
{
    return m_appId;
}

void WindowInfoK::setAppId(QString _appId)
{
    m_appId = _appId;
}

void WindowInfoK::activate()
{
    m_plasmaWindow->RequestActivate();
}

void WindowInfoK::minimize()
{
    m_plasmaWindow->RequestToggleMinimized();
}

bool WindowInfoK::isMinimized()
{
    return m_plasmaWindow->IsMinimized();
}

bool WindowInfoK::changeXid(XWindow _xid)
{
    xid = _xid;
    return true;
}

PlasmaWindow *WindowInfoK::getPlasmaWindow()
{
    return m_plasmaWindow.data();
}

bool WindowInfoK::updateGeometry()
{
    DockRect rect = m_plasmaWindow->Geometry();
    if (m_geometry == rect)
        return false;

    m_geometry = rect;
    return true;
}

void WindowInfoK::updateTitle()
{
    title = m_plasmaWindow->Title();
}

void WindowInfoK::updateDemandingAttention()
{
    m_demaningAttention = m_plasmaWindow->IsDemandingAttention();
}

void WindowInfoK::updateIcon()
{
    icon = m_plasmaWindow->Icon();
}

void WindowInfoK::updateAppId()
{
    m_appId = m_plasmaWindow->AppId();
}

void WindowInfoK::updateInternalId()
{
    m_internalId = m_plasmaWindow->InternalId();
}

void WindowInfoK::updateCloseable()
{
    m_closeable = m_plasmaWindow->IsCloseable();
}

void WindowInfoK::updateProcessInfo()
{
    pid = m_plasmaWindow->Pid();
    m_processInfo.reset(new ProcessInfo(pid));
}

/**
 * @brief WindowInfoK::getGeometry 获取窗口大小
 * @return
 */
DockRect WindowInfoK::getGeometry()
{
    return m_geometry;
}

int64_t WindowInfoK::getCreatedTime()
{
    return m_createdTime;
}

// 主要是为兼容X11
QString WindowInfoK::getDisplayName()
{
    return "";
}

QString WindowInfoK::getWindowType()
{
    return "Wayland";
}

void WindowInfoK::update()
{
    updateInternalId();
    updateAppId();
    updateIcon();
    updateTitle();
    updateGeometry();
    updateDemandingAttention();
    updateCloseable();
    updateProcessInfo();
}

void WindowInfoK::killClient()
{
}

QString WindowInfoK::uuid()
{
    return QString(m_plasmaWindow->uuid());
}

QString WindowInfoK::getInnerId()
{
    return QString::number(m_internalId);
}
