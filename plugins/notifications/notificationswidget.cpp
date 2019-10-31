/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     sbw <sbw@sbw.so>
 *
 * Maintainer: sbw <sbw@sbw.so>
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

#include "notificationswidget.h"

#include <QSvgRenderer>
#include <QPainter>
#include <QMouseEvent>
#include <QApplication>
#include <QIcon>

#include <DStyle>
#include <DGuiApplicationHelper>

DWIDGET_USE_NAMESPACE;

NotificationsWidget::NotificationsWidget(QWidget *parent)
    : QWidget(parent)
{
    setMouseTracking(true);
    setMinimumSize(PLUGIN_BACKGROUND_MIN_SIZE, PLUGIN_BACKGROUND_MIN_SIZE);

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, [ = ] {
        update();
    });
}

void NotificationsWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    QString iconName = "notification";
    int iconSize = PLUGIN_ICON_MAX_SIZE;
    if (height() <= PLUGIN_BACKGROUND_MIN_SIZE && DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType)
        iconName.append(PLUGIN_MIN_ICON_NAME);

    QPainter painter(this);
    const auto ratio = devicePixelRatioF();

    QPixmap pixmap = QIcon::fromTheme(iconName, QIcon::fromTheme(QString(":/icons/resources/icons/%1").arg(iconName))).pixmap(QSize(iconSize, iconSize) * ratio);
    pixmap.setDevicePixelRatio(ratio);

    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(pixmap.rect());

    painter.drawPixmap(rf.center() - rfp.center() / ratio, pixmap);
}

void NotificationsWidget::resizeEvent(QResizeEvent *event)
{
    QWidget::resizeEvent(event);
}
