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
#include "multiquickitem.h"
#include "singlequickitem.h"
#include "fullquickitem.h"
#include "quicksettingcontroller.h"

#include <DGuiApplicationHelper>
#include <DFontSizeManager>
#include <DPaletteHelper>

#include <QIcon>
#include <QPainterPath>
#include <QPushButton>
#include <QFontMetrics>

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

QuickSettingItem::QuickSettingItem(PluginsItemInterface *const pluginInter, QWidget *parent)
    : DockItem(parent)
    , m_pluginInter(pluginInter)
    , m_itemKey(QuickSettingController::instance()->itemKey(pluginInter))
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

DockItem::ItemType QuickSettingItem::itemType() const
{
    return DockItem::QuickSettingPlugin;
}

const QPixmap QuickSettingItem::dragPixmap()
{
    QPixmap pm = m_pluginInter->icon(DockPart::QuickPanel).pixmap(ICONWIDTH, ICONHEIGHT);

    QPainter pa(&pm);
    pa.setPen(foregroundColor());
    pa.setCompositionMode(QPainter::CompositionMode_SourceIn);
    pa.fillRect(pm.rect(), foregroundColor());

    QPixmap pmRet(ICONWIDTH + ICONSPACE + FONTSIZE * 2, ICONHEIGHT + ICONSPACE + FONTSIZE * 2);
    pmRet.fill(Qt::transparent);
    QPainter paRet(&pmRet);
    paRet.drawPixmap(QPoint((ICONSPACE + FONTSIZE * 2) / 2, 0), pm);
    paRet.setPen(pa.pen());

    QFont ft;
    ft.setPixelSize(FONTSIZE);
    paRet.setFont(ft);
    QTextOption option;
    option.setAlignment(Qt::AlignTop | Qt::AlignHCenter);
    paRet.drawText(QRect(QPoint(0, ICONHEIGHT + ICONSPACE),
                           QPoint(pmRet.width(), pmRet.height())), m_pluginInter->pluginDisplayName(), option);
    return pmRet;
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
    DPalette dpa = DPaletteHelper::instance()->palette(this);
    painter.fillRect(rect(), Qt::white);
}

QColor QuickSettingItem::foregroundColor() const
{
    DPalette dpa = DPaletteHelper::instance()->palette(this);
    // 此处的颜色是临时获取的，后期需要和设计师确认，改成正规的颜色
    if (m_pluginInter->status() == PluginsItemInterface::PluginStatus::Active)
        return dpa.color(DPalette::ColorGroup::Active, DPalette::ColorRole::Text);

    if (m_pluginInter->status() == PluginsItemInterface::PluginStatus::Deactive)
        return dpa.color(DPalette::ColorGroup::Disabled, DPalette::ColorRole::Text);

    return dpa.color(DPalette::ColorGroup::Normal, DPalette::ColorRole::Text);
}

QuickSettingItem *QuickSettingFactory::createQuickWidget(PluginsItemInterface * const pluginInter)
{
    QuickSettingController *quickController = QuickSettingController::instance();
    // 如果显示在面板的图标或者Widget为空，则不让显示(例如电池插件)
    if (!quickController->hasFlag(pluginInter, PluginFlag::Type_Common))
        return nullptr;

    if (quickController->hasFlag(pluginInter, PluginFlag::Quick_Multi))
        return new MultiQuickItem(pluginInter);

    if (quickController->hasFlag(pluginInter, PluginFlag::Quick_Full))
        return new FullQuickItem(pluginInter);

    if (quickController->hasFlag(pluginInter, PluginFlag::Quick_Single))
        return new SingleQuickItem(pluginInter);

    return nullptr;
}
