// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#include "quicksettingitem.h"
#include "pluginsiteminterface.h"
#include "imageutil.h"
#include "largerquickitem.h"
#include "standardquickitem.h"
#include "linequickitem.h"

#include <DGuiApplicationHelper>
#include <DFontSizeManager>
#include <DPaletteHelper>

#include <QIcon>
#include <QPainterPath>
#include <QPushButton>
#include <QFontMetrics>
#include <QPainter>

#define ICONWIDTH 24
#define ICONHEIGHT 24
#define ICONSPACE 10
#define RADIUS 8
#define FONTSIZE 10

#define BGWIDTH 128
#define BGSIZE 36
#define MARGINLEFTSPACE 10
#define OPENICONSIZE 12
#define MARGINRIGHTSPACE 9

DWIDGET_USE_NAMESPACE

QuickSettingItem::QuickSettingItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent)
    : QWidget(parent)
    , m_pluginInter(pluginInter)
    , m_itemKey(itemKey)
{
    setAcceptDrops(true);
    this->installEventFilter(this);
}

QuickSettingItem::~QuickSettingItem()
{
}

PluginsItemInterface *QuickSettingItem::pluginItem() const
{
    return m_pluginInter;
}

const QPixmap QuickSettingItem::dragPixmap()
{
    return grab();
}

const QString QuickSettingItem::itemKey() const
{
    return m_itemKey;
}

void QuickSettingItem::paintEvent(QPaintEvent *e)
{
    QWidget::paintEvent(e);
    QPainter painter(this);
    painter.setRenderHint(QPainter::RenderHint::Antialiasing);
    painter.setPen(foregroundColor());
    QPainterPath path;
    path.addRoundedRect(rect(), RADIUS, RADIUS);
    painter.setClipPath(path);
    // 绘制背景色
    QColor backColor(Qt::white);
    if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::ColorType::DarkType) {
        backColor = Qt::black;
    }
    backColor.setAlphaF(0.5);
    DPalette dpa = DPaletteHelper::instance()->palette(this);
    painter.fillRect(rect(), backColor);
}

QColor QuickSettingItem::foregroundColor() const
{
    DPalette dpa = DPaletteHelper::instance()->palette(this);
    // TODO: 此处的颜色是临时获取的，后期需要和设计师确认，改成正规的颜色
    if (m_pluginInter->status() == PluginsItemInterface::PluginMode::Active)
        return dpa.color(DPalette::ColorGroup::Active, DPalette::ColorRole::Text);

    if (m_pluginInter->status() == PluginsItemInterface::PluginMode::Deactive)
        return dpa.color(DPalette::ColorGroup::Disabled, DPalette::ColorRole::Text);

    return dpa.color(DPalette::ColorGroup::Normal, DPalette::ColorRole::Text);
}

QuickSettingItem *QuickSettingFactory::createQuickWidget(PluginsItemInterface * const pluginInter, const QString &itemKey)
{
    // 如果显示在面板的图标或者Widget为空，则不让显示(例如电池插件)
    if (!(pluginInter->flags() & PluginFlag::Type_Common))
        return nullptr;

    if (pluginInter->flags() & PluginFlag::Quick_Multi)
        return new LargerQuickItem(pluginInter, itemKey);

    if (pluginInter->flags() & PluginFlag::Quick_Full)
        return new LineQuickItem(pluginInter, itemKey);

    if (pluginInter->flags() & PluginFlag::Quick_Single)
        return new StandardQuickItem(pluginInter, itemKey);

    return nullptr;
}
