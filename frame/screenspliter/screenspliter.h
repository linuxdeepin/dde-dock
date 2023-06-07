// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

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
    explicit ScreenSpliter(AppItem *appItem, QObject *parent = nullptr);
    virtual ~ScreenSpliter();
    AppItem *appItem() const;

private:
    AppItem *m_appItem;
};

class ScreenSpliterFactory
{
public:
    static ScreenSpliter *createScreenSpliter(AppItem *appItem);
};

#endif // SCREENSPLITER_H
