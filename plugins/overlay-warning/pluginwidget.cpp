/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *             listenerri <listenerri@gmail.com>
 *
 * Maintainer: listenerri <listenerri@gmail.com>
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

#include "pluginwidget.h"

#include <QSvgRenderer>
#include <QPainter>
#include <QMouseEvent>
#include <QApplication>
#include <QIcon>

PluginWidget::PluginWidget(QWidget *parent)
    : QWidget(parent)
{
}

QSize PluginWidget::sizeHint() const
{
    return QSize(26, 26);
}

void PluginWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    QPixmap pixmap;
    QString iconName = ":/icons/resources/icons/overlay-warning.svg";
    int iconSize;
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

    if (displayMode == Dock::Efficient) {
//        iconName = iconName + "-symbolic";
        iconSize = 16;
    } else {
        iconSize = std::min(width(), height()) * 0.8;
    }

    pixmap = loadSvg(iconName, QSize(iconSize, iconSize));

    QPainter painter(this);
    painter.drawPixmap(rect().center() - pixmap.rect().center() / qApp->devicePixelRatio(), pixmap);
}

const QPixmap PluginWidget::loadSvg(const QString &fileName, const QSize &size) const
{
    const auto ratio = qApp->devicePixelRatio();

    QPixmap pixmap;
    pixmap = QIcon::fromTheme(fileName).pixmap(size * ratio);
    pixmap.setDevicePixelRatio(ratio);

    return pixmap;
}
