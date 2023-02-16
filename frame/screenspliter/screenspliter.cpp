// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
