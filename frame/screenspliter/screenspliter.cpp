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

ScreenSpliter::ScreenSpliter(AppItem *appItem, QObject *parent)
    : QObject(parent)
    , m_appItem(appItem)
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

ScreenSpliter *ScreenSpliterFactory::createScreenSpliter(AppItem *appItem)
{
    if (Utils::IS_WAYLAND_DISPLAY)
        return new ScreenSpliter_Wayland(appItem, appItem);

    return new ScreenSpliter_Xcb(appItem, appItem);
}
