// Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef TOUCHSIGNALMANAGER_H
#define TOUCHSIGNALMANAGER_H

#include "org_deepin_dde_gesture1.h"

#include <QObject>

using Gesture = org::deepin::dde::Gesture1;

class TouchSignalManager : public QObject
{
    Q_OBJECT

public:
    static TouchSignalManager *instance();
    bool isDragIconPress() const;

signals:
    // 转发后端拖拽图标触控按压信号，当前设计200ms
    void shortTouchPress(int time, double scaleX, double scaleY);
    void touchRelease(double scaleX, double scaleY);
    // 转发后端拖拽任务栏高度单指触控按压信号，当前设计1000ms
    void middleTouchPress(double scaleX, double scaleY);
    void touchMove(double scaleX, double scaleY);

private slots:
    void dealShortTouchPress(int time, double scaleX, double scaleY);
    void dealTouchRelease(double scaleX, double scaleY);
    void dealMiddleTouchPress(double scaleX, double scaleY);
    void dealTouchPress(int figerNum, int time, double scaleX, double scaleY);

private:
    explicit TouchSignalManager(QObject *parent = nullptr);

private:
    static TouchSignalManager *m_touchManager;
    Gesture *m_gestureInter;
    // 保存触控屏图标拖动长按状态，当前长按200ms
    bool m_dragIconPressed;
};

#endif // TOUCHSIGNALMANAGER_H
