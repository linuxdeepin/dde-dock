/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer: donghualin <donghualin@uniontech.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
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
