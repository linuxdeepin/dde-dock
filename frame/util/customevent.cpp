// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "customevent.h"

// 注册事件类型
static QEvent::Type pluginEventType = (QEvent::Type)QEvent::registerEventType(QEvent::User + 1001);

// 事件处理，当收到该事件的时候，加载插件
PluginLoadEvent::PluginLoadEvent()
    : QEvent(pluginEventType)
{
}

PluginLoadEvent::~PluginLoadEvent()
{
}

QEvent::Type PluginLoadEvent::eventType()
{
    return pluginEventType;
}
