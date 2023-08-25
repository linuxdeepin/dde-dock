// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: GPL-3.0-or-later

#include "dockdaemonadaptors.h"
#include "docksettings.h"
#include "taskmanager/taskmanager.h"

DockDaemonDBusAdaptor::DockDaemonDBusAdaptor(QObject *parent)
    : QDBusAbstractAdaptor(parent)
{
    // constructor
    setAutoRelaySignals(true);
    connect(TaskManager::instance(), &TaskManager::entryAdded, this, &DockDaemonDBusAdaptor::EntryAdded);
    connect(TaskManager::instance(), &TaskManager::entryRemoved, this, &DockDaemonDBusAdaptor::EntryRemoved);
    connect(TaskManager::instance(), &TaskManager::hideStateChanged, this, &DockDaemonDBusAdaptor::HideStateChanged);
    connect(TaskManager::instance(), &TaskManager::frontendWindowRectChanged, this, &DockDaemonDBusAdaptor::FrontendWindowRectChanged);
    connect(TaskManager::instance(), &TaskManager::showRecentChanged, this, &DockDaemonDBusAdaptor::showRecentChanged);
    connect(TaskManager::instance(), &TaskManager::showMultiWindowChanged, this, &DockDaemonDBusAdaptor::ShowMultiWindowChanged);
    connect(DockSettings::instance(), &DockSettings::positionModeChanged, this, &DockDaemonDBusAdaptor::PositionChanged);
    connect(DockSettings::instance(), &DockSettings::hideModeChanged, this, &DockDaemonDBusAdaptor::HideModeChanged);
    connect(DockSettings::instance(), &DockSettings::displayModeChanged, this, &DockDaemonDBusAdaptor::DisplayModeChanged);
    connect(DockSettings::instance(), &DockSettings::windowSizeEfficientChanged, this, &DockDaemonDBusAdaptor::WindowSizeEfficientChanged);
    connect(DockSettings::instance(), &DockSettings::windowSizeFashionChanged, this, &DockDaemonDBusAdaptor::WindowSizeFashionChanged);
}

DockDaemonDBusAdaptor::~DockDaemonDBusAdaptor()
{
    // destructor
}

int DockDaemonDBusAdaptor::displayMode() const
{
    return TaskManager::instance()->getDisplayMode();
}

void DockDaemonDBusAdaptor::setDisplayMode(int value)
{
    if (displayMode() != value) {
        TaskManager::instance()->setDisplayMode(value);
        Q_EMIT DisplayModeChanged(value);
    }
}

QStringList DockDaemonDBusAdaptor::dockedApps() const
{
    return TaskManager::instance()->getDockedApps();
}

int DockDaemonDBusAdaptor::hideMode() const
{
    return TaskManager::instance()->getHideMode();
}

void DockDaemonDBusAdaptor::setHideMode(int value)
{
    if (hideMode() != value) {
        TaskManager::instance()->setHideMode(static_cast<HideMode>(value));
        Q_EMIT HideModeChanged(value);
    }
}

int DockDaemonDBusAdaptor::hideState() const
{
    return TaskManager::instance()->getHideState();
}

uint DockDaemonDBusAdaptor::hideTimeout() const
{
    return TaskManager::instance()->getHideTimeout();
}

void DockDaemonDBusAdaptor::setHideTimeout(uint value)
{
    if (hideTimeout() != value) {
        TaskManager::instance()->setHideTimeout(value);
        Q_EMIT HideTimeoutChanged(value);
    }
}

uint DockDaemonDBusAdaptor::windowSizeEfficient() const
{
    return TaskManager::instance()->getWindowSizeEfficient();
}

void DockDaemonDBusAdaptor::setWindowSizeEfficient(uint value)
{
    if (windowSizeEfficient() != value) {
        TaskManager::instance()->setWindowSizeEfficient(value);
        Q_EMIT WindowSizeEfficientChanged(value);
    }
}

uint DockDaemonDBusAdaptor::windowSizeFashion() const
{
    return TaskManager::instance()->getWindowSizeFashion();
}

void DockDaemonDBusAdaptor::setWindowSizeFashion(uint value)
{
    if (windowSizeFashion() != value) {
        TaskManager::instance()->setWindowSizeFashion(value);
        Q_EMIT WindowSizeFashionChanged(value);
    }
}

QRect DockDaemonDBusAdaptor::frontendWindowRect() const
{
    return TaskManager::instance()->getFrontendWindowRect();
}

uint DockDaemonDBusAdaptor::iconSize() const
{
    return TaskManager::instance()->getIconSize();
}

void DockDaemonDBusAdaptor::setIconSize(uint value)
{
    if (iconSize() != value) {
        TaskManager::instance()->setIconSize(value);
        Q_EMIT IconSizeChanged(value);
    }
}

int DockDaemonDBusAdaptor::position() const
{
    return TaskManager::instance()->getPosition();
}

void DockDaemonDBusAdaptor::setPosition(int value)
{
    if (position() != value) {
        TaskManager::instance()->setPosition(value);
        Q_EMIT PositionChanged(value);
    }
}

uint DockDaemonDBusAdaptor::showTimeout() const
{
    return TaskManager::instance()->getShowTimeout();
}

void DockDaemonDBusAdaptor::setShowTimeout(uint value)
{
    if (showTimeout() != value) {
        TaskManager::instance()->setShowTimeout(value);
        Q_EMIT ShowTimeoutChanged(value);
    }
}

bool DockDaemonDBusAdaptor::showRecent() const
{
    return DockSettings::instance()->showRecent();
}

bool DockDaemonDBusAdaptor::showMultiWindow() const
{
    return TaskManager::instance()->showMultiWindow();
}

void DockDaemonDBusAdaptor::CloseWindow(uint win)
{
    TaskManager::instance()->closeWindow(win);
}

// for debug
QStringList DockDaemonDBusAdaptor::GetEntryIDs()
{
    return TaskManager::instance()->getEntryIDs();
}

bool DockDaemonDBusAdaptor::IsDocked(const QString &desktopFile)
{
    return TaskManager::instance()->isDocked(desktopFile);
}

bool DockDaemonDBusAdaptor::IsOnDock(const QString &desktopFile)
{
    return TaskManager::instance()->isOnDock(desktopFile);
}

void DockDaemonDBusAdaptor::MoveEntry(int index, int newIndex)
{
    TaskManager::instance()->moveEntry(index, newIndex);
}

QString DockDaemonDBusAdaptor::QueryWindowIdentifyMethod(uint win)
{
    return TaskManager::instance()->queryWindowIdentifyMethod(win);
}

QStringList DockDaemonDBusAdaptor::GetDockedAppsDesktopFiles()
{
    return TaskManager::instance()->getDockedAppsDesktopFiles();
}

QString DockDaemonDBusAdaptor::GetPluginSettings()
{
    return TaskManager::instance()->getPluginSettings();
}

void DockDaemonDBusAdaptor::SetPluginSettings(QString jsonStr)
{
    TaskManager::instance()->setPluginSettings(jsonStr);
}

void DockDaemonDBusAdaptor::MergePluginSettings(QString jsonStr)
{
    TaskManager::instance()->mergePluginSettings(jsonStr);
}

void DockDaemonDBusAdaptor::RemovePluginSettings(QString key1, QStringList key2List)
{
    TaskManager::instance()->removePluginSettings(key1, key2List);
}

bool DockDaemonDBusAdaptor::RequestDock(const QString &desktopFile, int index)
{
    return TaskManager::instance()->requestDock(desktopFile, index);
}

bool DockDaemonDBusAdaptor::RequestUndock(const QString &desktopFile)
{
    return TaskManager::instance()->requestUndock(desktopFile);
}

void DockDaemonDBusAdaptor::SetShowRecent(bool visible)
{
    DockSettings::instance()->setShowRecent(visible);
}

void DockDaemonDBusAdaptor::SetShowMultiWindow(bool showMultiWindow)
{
    TaskManager::instance()->setShowMultiWindow(showMultiWindow);
}

void DockDaemonDBusAdaptor::SetFrontendWindowRect(int x, int y, uint width, uint height)
{
    TaskManager::instance()->setFrontendWindowRect(x, y, width, height);
}
