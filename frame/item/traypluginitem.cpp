// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "traypluginitem.h"
#include "utils.h"

#include <QEvent>

TrayPluginItem::TrayPluginItem(PluginsItemInterface * const pluginInter, const QString &itemKey, const QJsonObject &metaData, QWidget *parent)
    : PluginsItem(pluginInter, itemKey, metaData, parent)
{
    centralWidget()->installEventFilter(this);
}

void TrayPluginItem::setSuggestIconSize(QSize size)
{
    // invoke the method "setSuggestIconSize" of FashionTrayItem class
    QMetaObject::invokeMethod(centralWidget(), "setSuggestIconSize", Qt::QueuedConnection, Q_ARG(QSize, size));
}

void TrayPluginItem::setRightSplitVisible(const bool visible)
{
    // invoke the method "setRightSplitVisible" of FashionTrayItem class
    QMetaObject::invokeMethod(centralWidget(), "setRightSplitVisible", Qt::QueuedConnection, Q_ARG(bool, visible));
}

int TrayPluginItem::trayVisibleItemCount()
{
    return m_trayVisableItemCount;
}

bool TrayPluginItem::eventFilter(QObject *watched, QEvent *e)
{
    // 时尚模式下
    // 监听插件Widget的"FashionTraySize"属性
    // 当接收到这个属性变化的事件后，重新计算和设置dock的大小

    if (watched == centralWidget()) {
        if (e->type() == QEvent::MouseButtonPress || e->type() == QEvent::MouseButtonRelease) {
            const QGSettings *settings = Utils::ModuleSettingsPtr("systemtray", QByteArray(), this);
            if (settings && settings->keys().contains("control") && settings->get("control").toBool()) {
                return true;
            }
        }
    }

    if (watched == centralWidget() && e->type() == QEvent::DynamicPropertyChange) {
        const QString &propertyName = static_cast<QDynamicPropertyChangeEvent *>(e)->propertyName();
        if (propertyName == "TrayVisableItemCount") {
            m_trayVisableItemCount = watched->property("TrayVisableItemCount").toInt();
            Q_EMIT trayVisableCountChanged(m_trayVisableItemCount);
        }
    }

    return PluginsItem::eventFilter(watched, e);
}
