/*
 * Copyright (C) 2023 ~ 2023 Deepin Technology Co., Ltd.
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
#include "iconmanager.h"
#include "dockplugincontroller.h"
#include "pluginsiteminterface.h"

#include <DDciIcon>
#include <DWindowManagerHelper>
#include <DSysInfo>
#include <DPlatformTheme>

#include <QPainter>
#include <QPainterPath>

#define ITEMSPACE 6
#define ITEMHEIGHT 16
#define ITEMWIDTH 18
#define MINISIZE 1
#define STARTPOS 2

static QStringList pluginNames = {"power", "sound", "network"};

DGUI_USE_NAMESPACE

IconManager::IconManager(DockPluginController *pluginController, QObject *parent)
    : QObject{parent}
    , m_pluginController(pluginController)
    , m_size(QSize(ITEMWIDTH, ITEMHEIGHT))
    , m_position(Dock::Position::Bottom)
    , m_displayMode(Dock::DisplayMode::Efficient)
{
}

void IconManager::updateSize(QSize size)
{
    m_size = size;
}

void IconManager::setPosition(Dock::Position position)
{
    m_position = position;
}

void IconManager::setDisplayMode(Dock::DisplayMode displayMode)
{
    m_displayMode = displayMode;
}

QPixmap IconManager::pixmap(DGuiApplicationHelper::ColorType colorType) const
{
    QList<PluginsItemInterface *> plugins;
    for (const QString &pluginName : pluginNames) {
        PluginsItemInterface *plugin = findPlugin(pluginName);
        if (plugin)
            plugins << plugin;
    }

    if (plugins.size() < 2) {
        // 缺省图标
        DDciIcon::Theme theme = DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::DarkType ? DDciIcon::Light : DDciIcon::Dark;
        DDciIcon dciIcon(QString(":/resources/dock_control.dci"));
        QPixmap pixmap = dciIcon.pixmap(QCoreApplication::testAttribute(Qt::AA_UseHighDpiPixmaps) ? 1 : qApp->devicePixelRatio(), ITEMHEIGHT, theme, DDciIcon::Normal);
        QColor foreColor = (colorType == DGuiApplicationHelper::ColorType::DarkType ? Qt::white : Qt::black);
        foreColor.setAlphaF(0.8);
        QPainter pa(&pixmap);
        pa.setCompositionMode(QPainter::CompositionMode_SourceIn);
        pa.fillRect(pixmap.rect(), foreColor);
        return pixmap;
    }

    int itemSpace = 0;
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom)
        itemSpace = (m_displayMode == Dock::DisplayMode::Efficient ? 8 : 10);
    else
        itemSpace = 2;
    // 组合图标
    QPixmap pixmap;
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        // 高效模式下，高度固定为30, 时尚模式下，高度随着任务栏的大小变化而变化
        int iconHeight = (m_displayMode == Dock::DisplayMode::Efficient ? 30 : m_size.height() - 8);
        if (iconHeight <= 0)
            iconHeight = MINISIZE;
        int iconWidth = STARTPOS;
        for (PluginsItemInterface *plugin : plugins) {
            QIcon icon = plugin->icon(DockPart::QuickShow);
            QSize iconSize = QSize(ITEMWIDTH, ITEMHEIGHT) * qApp->devicePixelRatio();
            QList<QSize> iconSizes = icon.availableSizes();
            if (iconSizes.size() > 0)
                iconSize = iconSizes.first();
            iconWidth += iconSize.width();
        }
        iconWidth += itemSpace * (plugins.size() - 1);
        pixmap = QPixmap(iconWidth, iconHeight);
    } else {
        // 左右方向，高效模式下，宽度固定为30，时尚模式下，宽度随任务栏的大小变化而变化
        int iconWidth = m_displayMode == Dock::DisplayMode::Efficient ? 30 : m_size.width() - 8;
        if (iconWidth <= 0)
            iconWidth = MINISIZE;
        int iconHeight = STARTPOS;
        for (PluginsItemInterface *plugin : plugins) {
            QIcon icon = plugi