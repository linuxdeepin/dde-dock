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
#include "linequickitem.h"
#include "pluginsiteminterface.h"

#include <DBlurEffectWidget>

#include <QHBoxLayout>

DWIDGET_USE_NAMESPACE

LineQuickItem::LineQuickItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent)
    : QuickSettingItem(pluginInter, itemKey, parent)
    , m_centerWidget(pluginInter->itemWidget(QUICK_ITEM_KEY))
    , m_centerParentWidget(nullptr)
    , m_effectWidget(new DBlurEffectWidget(this))
{
    initUi();
    initConnection();
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
    m_effectWidget->setBlurRectXRadius(8);
    m_effectWidget->setBlurRectYRadius(8);
    onThemeTypeChanged(DGuiApplicationHelper::instance()->themeType());

    // 如果图标不为空
    if (!m_centerWidget)
        return;

    m_centerWidget->setVisible(true);
    m_centerParentWidget = m_centerWidget->parentWidget();

    QHBoxLayout *layout = new QHBoxLayout(m_effectWidget);
    layout->setContentsMargins(0, 0, 0, 0);
    layout->setAlignment(Qt::AlignHCenter);
    layout->addWidget(m_centerWidget);

    QHBoxLayout *mainLayout = new QHBoxLayout(this);
    mainLayout->setContentsMargins(0, 0, 0, 0);
    mainLayout->addWidget(m_effectWidget);

    m_centerWidget->installEventFilter(this);
}

void LineQuickItem::initConnection()
{
    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, &LineQuickItem::onThemeTypeChanged);
}

void LineQuickItem::resizeS