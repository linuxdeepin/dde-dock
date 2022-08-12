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
#include "screenspliter.h"
#include "appitem.h"
#include "utils.h"
#include "screenspliter_xcb.h"
#include "screenspliter_wayland.h"

bool ScreenSpliter::releaseSplit()
{
    return true;
}

ScreenSpliter::ScreenSpliter(AppItem *appItem, DockEntryInter *entryInter, QObject *parent)
    : QObject(parent)
    , m_appItem(appItem)
    , m_entryInter(entryInter)
{
}

ScreenSpliter::~ScreenSpliter()
{
    m_appItem = nullptr;
}

AppItem *ScreenSpliter::appItem() const
{
    return m_appItem;
}

DockEntryInter *ScreenSpliter::entryInter() const
{
    return m_entryInter;
}

ScreenSpliter *ScreenSpliterFactory::createScreenSpliter(AppItem *appItem, DockEntryInter *entryInter)
{
    if (Utils::IS_WAYLAND_DISPLAY)
        return new ScreenSpliter_Wayland(appItem, entryInter, appItem);

    return new ScreenSpliter_Xcb(appItem, entryInter, appItem);
}
