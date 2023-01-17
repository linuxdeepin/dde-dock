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

#include <QPainter>
#include <QPainterPath>

#define ITEMSPACE 6
#define ITEMHEIGHT 16
#define ITEMWIDTH 18
static QStringList pluginNames = {"power", "sound", "network"};

DGUI_USE_NAMESPACE

IconManager::IconManager(DockPluginController *pluginController, QObject *parent)
    : QObject{parent}
    , m_pluginController(pluginController)
    , m_size(QSize(ITEMWIDTH, ITEMHEIGHT))
    , m_position(Dock::Position::Bottom)
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

QPixmap IconManager::pixmap() const
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
        return dciIcon.pixmap(QCoreApplication::testAttribute(Qt::AA_UseHighDpiPixmaps) ? 1 : qApp->devicePixelRatio(), ITEMHEIGHT, theme, DDciIcon::Normal);
    }

    // 组合图标
    QPixmap pixmap;
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        pixmap = QPixmap(ITEMWIDTH * plugins.size() + ITEMSPACE * (plugins.size() + 1), m_size.height() - 8);
    } else {
        pixmap = QPixmap(m_size.width() - 8, ITEMWIDTH * plugins.size() + ITEMSPACE * (plugins.size() + 1));
    }

    pixmap.fill(Qt::transparent);
    QPainter painter(&pixmap);
    painter.setRenderHints(QPainter::SmoothPixmapTransform | QPainter::Antialiasing);
    QColor backColor = DGuiApplicationHelper::ColorType::DarkType == DGuiApplicationHelper::instance()->themeType() ? QColor(20, 20, 20) : Qt::white;
    backColor.setAlphaF(0.2);
    QPainterPath path;
    path.addRoundedRect(pixmap.rect(), 5, 5);
    painter.fillPath(path, backColor);
    if (m_position == Dock::Position::Top || m_position == Dock::Position::Bottom) {
        QPoint pointPixmap(ITEMSPACE, (pixmap.height() - ITEMHEIGHT) / 2);
        for (PluginsItemInterface *plugin : plugins) {
            QIcon icon = plugin->icon(DockPart::QuickShow);
             QPixmap pixmapDraw = icon.pixmap(ITEMHEIGHT, ITEMHEIGHT);
             painter.drawPixmap(pointPixmap, pixmapDraw);
             pointPixmap.setX(pointPixmap.x() + ITEMWIDTH + ITEMSPACE);
         }
     } else {
         QPoint pointPixmap((pixmap.width() - ITEMWIDTH) / 2, ITEMSPACE);
         for (PluginsItemInterface *plugin : plugins) {
             QIcon icon = plugin->icon(DockPart::QuickShow);
             QPixmap pixmapDraw = icon.pixmap(ITEMHEIGHT, ITEMHEIGHT);
             painter.drawPixmap(pointPixmap, pixmapDraw);
             pointPixmap.setY(pointPixmap.y() + ITEMWIDTH + ITEMSPACE);
         }
     }
     painter.end();
     return pixmap;
}

bool IconManager::isFixedPlugin(PluginsItemInterface *plugin) const
{
    return pluginNames.contains(plugin->pluginName());
}

PluginsItemInterface *IconManager::findPlugin(const QString &pluginName) const
{
    QList<PluginsItemInterface *> plugins = m_pluginController->currentPlugins();
    for (PluginsItemInterface *plugin : plugins) {
        if (plugin->pluginName() == pluginName)
            return plugin;
    }

    return nullptr;
}
