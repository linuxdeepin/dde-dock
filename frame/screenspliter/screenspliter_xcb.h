// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SCREENSPLITER_XCB_H
#define SCREENSPLITER_XCB_H

#include "screenspliter.h"

#include <QRect>

class WindowInfo;
typedef QMap<quint32, WindowInfo> WindowInfoMap;

class ScreenSpliter_Xcb : public ScreenSpliter
{
public:
    explicit ScreenSpliter_Xcb(AppItem *appItem, QObject *parent = nullptr);

    void startSplit(const QRect &rect) override;
    bool split(ScreenSpliter::SplitDirection direction) override;
    bool suportSplitScreen() override;
    bool releaseSplit() override;

private:
    uint32_t direction_x11(ScreenSpliter::SplitDirection direction);
    void showSplitScreenEffect(const QRect &rect, bool visible);
    bool windowSupportSplit(quint32 winId);
};

#endif // SCREENSPLITER_XCB_H
