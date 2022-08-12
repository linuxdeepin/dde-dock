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
#ifndef SCREENSPLITER_H
#define SCREENSPLITER_H

#include "dbusutil.h"

#include <QObject>

class AppItem;

class ScreenSpliter : public QObject
{
    Q_OBJECT

public:
    enum SplitDirection {
        None,           // 无操作
        Left,           // 左分屏
        Right,          // 右分屏
        Top,            // 上分屏
        Bottom,         // 下分屏
        LeftTop,        // 左上
        RightTop,       // 右上
        LeftBottom,     // 左下
        RightBottom,    // 右下
        Full            // 全屏
    };

public:
    virtual void startSplit(const QRect &) = 0;                      // 触发分屏提示效果
    virtual bool split(SplitDirection) = 0;                          // 开始触发分屏
    virtual bool suportSplitScreen() = 0;                            // 是否支持分屏
    virtual bool releaseSplit();                                     // 释放分屏

protected:
    explicit ScreenSpliter(AppItem *appItem, DockEntryInter *entryInter, QObject *parent = nullptr);
    virtual ~ScreenSpliter();
    AppItem *appItem() const;
    DockEntryInter *entryInter() const;

private:
    AppItem *m_appItem;
    DockEntryInter *m_entryInter;
};

class ScreenSpliterFactory
{
public:
    static ScreenSpliter *createScreenSpliter(AppItem *appItem, DockEntryInter *entryInter);
};

#endif // SCREENSPLITER_H
