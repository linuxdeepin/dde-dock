/*
 * Copyright (C) 2018 ~ 2020 Uniontech Technology Co., Ltd.
 *
 * Author:     liuxing <liuxing@uniontech.com>
 *
 * Maintainer: liuxing <liuxing@uniontech.com>
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

#ifndef TOUCHSIGNALMANAGER_H
#define TOUCHSIGNALMANAGER_H

#include <QObject>

#include <com_deepin_daemon_gesture.h>

using Gesture = com::deepin::daemon::Gesture;

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
