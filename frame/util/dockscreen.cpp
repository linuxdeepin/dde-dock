// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "dockscreen.h"
#include "displaymanager.h"

DockScreen::DockScreen()
    : m_primary(DisplayManager::instance()->primary())
    , m_currentScreen(m_primary)
    , m_lastScreen(m_primary)
{
}

DockScreen *DockScreen::instance()
{
    static DockScreen instance;
    return &instance;
}

const QString &DockScreen::current() const
{
    return m_currentScreen;
}

const QString &DockScreen::last() const
{
    return m_lastScreen;
}

const QString &DockScreen::primary() const
{
    return m_primary;
}

void DockScreen::updateDockedScreen(const QString &screenName)
{
    m_lastScreen = m_currentScreen;
    m_currentScreen = screenName;
}

void DockScreen::updatePrimary(const QString &primary)
{
    m_primary = primary;
    if (m_currentScreen.isEmpty()) {
        updateDockedScreen(primary);
    }
}
