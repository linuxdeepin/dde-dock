// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "testplugin.h"

#include <QWidget>

TestPlugin::TestPlugin()
    : m_sortKey(0)
    , m_type(Normal)
    , m_widget(new QWidget)
{
}

TestPlugin::~TestPlugin()
{
    if (m_widget) {
        delete m_widget;
        m_widget = nullptr;
    }
}

const QString TestPlugin::pluginName() const
{
    return QString(Name);
}

const QString TestPlugin::pluginDisplayName() const
{
    return QString(Name);
}

void TestPlugin::init(PluginProxyInterface *)
{
}

QWidget *TestPlugin::itemWidget(const QString &)
{
    return m_widget;
}

const QString TestPlugin::itemContextMenu(const QString &itemKey)
{
    return QString("test plugin menu: " + itemKey);
}

int TestPlugin::itemSortKey(const QString &)
{
    return m_sortKey;
}

void TestPlugin::setSortKey(const QString &, const int order)
{
    m_sortKey = order;
}

PluginsItemInterface::PluginSizePolicy TestPlugin::pluginSizePolicy() const
{
    return PluginsItemInterface::Custom;
}

PluginsItemInterface::PluginType TestPlugin::type()
{
    return m_type;
}

void TestPlugin::setType(const PluginType type)
{
    m_type = type;
}
