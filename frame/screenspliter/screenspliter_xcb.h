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
#ifndef SCREENSPLITER_XCB_H
#define SCREENSPLITER_XCB_H

#include "screenspliter.h"

#include <QRect>

class WindowInfo;
typedef QMap<quint32, WindowInfo> WindowInfoMap;

class ScreenSpliter_Xcb : public ScreenSpliter
{
public:
    explicit ScreenSpliter_Xcb(AppItem *appItem, DockEntryInter *entryInter, QObject *parent = nullptr);

    void startSplit(const QRect &rect) override;
    bool split(ScreenSpliter::SplitDirection direction) override;
    bool suportSplitScreen() override;
    bool releaseSplit() override;

private:
    quint32 splittingWindowWId();
    uint32_t direction_x11(ScreenSpliter::SplitDirection direction);
    void showSplitScreenEffect(const QRect &rect, bool visible);
    bool openWindow();

private Q_SLOTS:
    void onUpdateWindowInfo(const WindowInfoMap &info);

private:
    bool m_isSplitCreateWindow;
    QRect m_effectRect;
};

#endif // SCREENSPLITER_XCB_H
