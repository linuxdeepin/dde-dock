/*
 * Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
 *
 * Author:     wangshaojun <wangshaojun_cm@deepin.com>
 *
 * Maintainer: wangshaojun <wangshaojun_cm@deepin.com>
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

#include "showdesktopwidget.h"
#include "showdesktopplugin.h"

#include <QPainter>
#include <QIcon>
#include <QMouseEvent>

ShowDesktopWidget::ShowDesktopWidget(QWidget *parent)
    : QWidget(parent)
{
}

void ShowDesktopWidget::refreshIcon()
{
    update();
}

void ShowDesktopWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    const auto ratio = devicePixelRatioF();
    QPixmap icon;

    if (Dock::Fashion == qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>()) {
        icon = QIcon::fromTheme("deepin-toggle-desktop").pixmap(size() * 0.8 * ratio);
    } else {
        icon = QIcon::fromTheme("deepin-toggle-desktop").pixmap(size() * 0.7 * ratio);
    }

    icon.setDevicePixelRatio(ratio);

    QPainter painter(this);
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(icon.rect());
    painter.drawPixmap(rf.center() - rfp.center() / ratio, icon);
}
