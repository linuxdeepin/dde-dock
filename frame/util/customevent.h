// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef CUSTOMEVENT_H
#define CUSTOMEVENT_H

#include <QEvent>

// 该插件用于处理插件的延迟加载，当退出安全模式后，会收到该事件并加载插件
class PluginLoadEvent : public QEvent
{
public:
    PluginLoadEvent();
    ~PluginLoadEvent() override;

    static Type eventType();
};

#endif // CUSTOMEVENT_H
