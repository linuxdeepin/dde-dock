// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "linequickitem.h"
#include "pluginsiteminterface.h"

#include <DBlurEffectWidget>

#include <QHBoxLayout>

DWIDGET_USE_NAMESPACE

LineQuickItem::LineQuickItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent)
    : QuickSettingItem(pluginInter, itemKey, parent)
    , m_centerWidget(pluginInter->itemWidget(QUICK_ITEM_KEY))
    , m_centerParentWidget(nullptr)
{
    initUi();
    QMetaObject::invokeMethod(this, &LineQuickItem::resizeSelf, Qt::QueuedConnection);
}

LineQuickItem::~LineQuickItem()
{
    if (m_centerWidget)
        m_centerWidget->setParent(nullptr);
}

void LineQuickItem::doUpdate()
{
    if (m_centerWidget)
        m_centerWidget->update();
}

void LineQuickItem::detachPlugin()
{
    if (m_centerWidget)
        m_centerWidget->setParent(m_centerParentWidget);
}

QuickSettingItem::QuickItemStyle LineQuickItem::type() const
{
    return QuickSettingItem::QuickItemStyle::Line;
}

bool LineQuickItem::eventFilter(QObject *obj, QEvent *event)
{
    if (obj == m_centerWidget && event->type() == QEvent::Resize)
        resizeSelf();

    return QuickSettingItem::eventFilter(obj, event);
}

void LineQuickItem::initUi()
{
    // 如果图标不为空
    if (!m_centerWidget)
        return;

    m_centerWidget->setVisible(true);
    m_centerParentWidget = m_centerWidget->parentWidget();

    QHBoxLayout *layout = new QHBoxLayout;
    layout->setContentsMargins(0, 0, 0, 0);
    layout->setAlignment(Qt::AlignHCenter);
    layout->addWidget(m_centerWidget);

    QHBoxLayout *mainLayout = new QHBoxLayout(this);
    mainLayout->setContentsMargins(0, 0, 0, 0);
    mainLayout->addWidget(m_centerWidget);

    m_centerWidget->installEventFilter(this);
}

void LineQuickItem::resizeSelf()
{
    if (!m_centerWidget)
        return;

    setFixedHeight(m_centerWidget->height());
}
