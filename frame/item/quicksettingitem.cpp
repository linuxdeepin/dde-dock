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

#include <DGuiApplicationHelper>
#include <DFontSizeManager>

#include <QIcon>
#include <QPainterPath>

#define ICONWIDTH 24
#define ICONHEIGHT 24
#define ICONSPACE 10
#define RADIUS 8
#define FONTSIZE 10

#define BGWIDTH 128
#define BGSIZE 36
#define MARGINLEFTSPACE 10
#define OPENICONSIZE 12
#define MARGINRIGHTSPACE 12

static QSize expandSize = QSize(6, 10);

QuickSettingItem::QuickSettingItem(PluginsItemInterface *const pluginInter, const QString &itemKey, QWidget *parent)
    : DockItem(parent)
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
    painter.fillRect(rect(), backgroundColor());
    // 让图标填上前景色
    QPixmap pm = m_pluginInter->icon(DockPart::QuickPanel).pixmap(ICONWIDTH, ICONHEIGHT);
    QPainter pa(&pm);
    pa.setCompositionMode(QPainter::CompositionMode_SourceIn);
    pa.fillRect(pm.rect(), painter.pen().brush());
    if (m_pluginInter->isPrimary()) {
        // 如果是主图标，则显示阴影背景
        int marginYSpace = yMarginSpace();
        QRect iconBg(MARGINLEFTSPACE, marginYSpace, BGSIZE, BGSIZE);
        QPixmap bgPixmap = ImageUtil::getShadowPixmap(pm, shadowColor(), QSize(BGSIZE, BGSIZE));
        painter.drawPixmap(iconBg, bgPixmap);
        // 绘制文字
        painter.setPen(QColor(0, 0, 0));

        QRect rctPluginName(iconBg.right() + 10, iconBg.top(), BGWIDTH - BGSIZE - OPENICONSIZE - 10 * 2, BGSIZE / 2);
        QFont font = DFontSizeManager::instance()->t6();
        font.setBold(true);
        painter.setFont(font);
        QTextOption textOption;
        textOption.setAlignment(Qt::AlignLeft | Qt::AlignVCenter);
        QString displayName = QFontMetrics(font).elidedText(m_pluginInter->pluginDisplayName(), Qt::TextElideMode::ElideRight, rctPluginName.width());
        QFontMetrics fm(font);
        painter.drawText(rctPluginName, displayName, textOption);
        // 绘制下方啊的状态文字
        QRect rctPluginStatus(rctPluginName.x(), rctPluginName.bottom() + 1,
                              rctPluginName.width(), BGSIZE / 2);
        font = DFontSizeManager::instance()->t10();
        painter.setFont(font);
        QString description = QFontMetrics(font).elidedText(m_pluginInter->description(), Qt::TextElideMode::ElideRight, rctPluginStatus.width());
        painter.drawText(rctPluginStatus, description, textOption);
        // 绘制右侧的展开按钮
        QPen pen;
        pen.setColor(QColor(0, 0, 0));
        pen.setWidth(2);
        painter.setPen(pen);
        int iconLeft = rect().width() - MARGINRIGHTSPACE - expandSize.width();
        int iconRight = rect().width() - MARGINRIGHTSPACE;
        painter.drawLine(QPoint(iconLeft, (iconBg.y() + (iconBg.height() - expandSize.height()) / 2)),
                         QPoint(iconRight, (iconBg.y() + iconBg.height() / 2)));
        painter.drawLine(QPoint(iconRight, (iconBg.y() + iconBg.height() / 2)),
                         QPoint(iconLeft, (iconBg.y() + (iconBg.height() + expandSize.height()) / 2)));
    } else {
        // 绘制图标
        QRect rctIcon = iconRect();
        painter.drawPixmap(rctIcon, pm);
        // 绘制文字
        QFont ft;
        ft.setPixelSize(FONTSIZE);
        painter.setFont(ft);
        QTextOption option;
        option.setAlignment(Qt::AlignTop | Qt::AlignHCenter);
        painter.drawText(QRect(QPoint(0, rctIcon.top() + ICONHEIGHT + ICONSPACE),
                               QPoint(width(), height())), m_pluginInter->pluginDisplayName(), option);
    }
}

QRect QuickSettingItem::iconRect()
{
    int left = (width() - ICONWIDTH) / 2;
    int top = (height() - ICONHEIGHT - ICONSPACE - 10) / 2;
    return QRect(left, top, ICONWIDTH, ICONHEIGHT);
}

QColor QuickSettingItem::foregroundColor() const
{
    // 此处的颜色是临时获取的，后期需要和设计师确认，改成正规的颜色
    if (m_pluginInter->status() == PluginsItemInterface::PluginStatus::Active)
        return QColor(0, 129, 255);

    if (m_pluginInter->status() == PluginsItemInterface::PluginStatus::Deactive)
        return QColor(51, 51, 51);

    return QColor(181, 181, 181);
}

QColor QuickSettingItem::backgroundColor() const
{
    // 此处的颜色是临时获取的，后期需要和设计师确认，改成正规的颜色
    if (m_pluginInter->status() == PluginsItemInterface::PluginStatus::Active)
        return QColor(250, 250, 252);

    return QColor(241, 241, 246);
}

QColor QuickSettingItem::shadowColor() const
{
    // 此处的颜色是临时获取的，后期需要和设计师确认，改成正规的颜色
    if (m_pluginInter->status() == PluginsItemInterface::PluginStatus::Active)
        return QColor(217, 219, 226);

    return QColor(199, 203, 222);
}

void QuickSettingItem::mouseReleaseEvent(QMouseEvent *event)
{
    // 如果是鼠标的按下事件
    if (m_pluginInter->isPrimary()) {
        QMouseEvent *mouseEvent = static_cast<QMouseEvent *>(event);
        QRect rctExpand(rect().width() - MARGINRIGHTSPACE - expandSize.width(),
                        (rect().height() - expandSize.height()) / 2,
                        expandSize.width(), expandSize.height());
        if (rctExpand.contains(mapFromGlobal(mouseEvent->globalPos())))
            Q_EMIT detailClicked(m_pluginInter);
    } else {
        const QString command = m_pluginInter->itemCommand(m_itemKey);
        if (!command.isEmpty())
            QProcess::startDetached(command);

        if (QWidget *w = m_pluginInter->itemPopupApplet(m_itemKey))
            showPopupApplet(w);
    }
}

int QuickSettingItem::yMarginSpace()
{
    return (rect().height() - BGSIZE) / 2;
}
