/*
 * Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     listenerri <listenerri@gmail.com>
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

#include "onboarditem.h"

#include <QSvgRenderer>
#include <QPainter>
#include <QMouseEvent>
#include <QApplication>
#include <QIcon>

#include <DStyle>

DWIDGET_USE_NAMESPACE;

OnboardItem::OnboardItem(QWidget *parent)
    : QWidget(parent)
    , m_hover(false)
    , m_pressed(false)
{
    setMouseTracking(true);
    setMinimumSize(PLUGIN_BACKGROUND_MIN_SIZE, PLUGIN_BACKGROUND_MIN_SIZE);
    setMaximumSize(PLUGIN_BACKGROUND_MAX_SIZE, PLUGIN_BACKGROUND_MAX_SIZE);
}

QSize OnboardItem::sizeHint() const
{
    return QSize(PLUGIN_BACKGROUND_MAX_SIZE, PLUGIN_BACKGROUND_MAX_SIZE);
}

void OnboardItem::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    QPixmap pixmap;
    QString iconName = "deepin-virtualkeyboard";
    int iconSize = PLUGIN_ICON_MAX_SIZE;
    const Dock::DisplayMode displayMode = qApp->property(PROP_DISPLAY_MODE).value<Dock::DisplayMode>();

    if (displayMode == Dock::Efficient) {
        iconName = iconName + "-symbolic";
    }

    QPainter painter(this);
    if (std::min(width(), height()) > PLUGIN_BACKGROUND_MIN_SIZE) {
        painter.setRenderHint(QPainter::Antialiasing, true);
        painter.setOpacity(0.5);

        DStyleHelper dstyle(style());
        const int radius = dstyle.pixelMetric(DStyle::PM_FrameRadius);

        QPainterPath path;
        path.addRoundedRect(rect(), radius, radius);

        QColor color = QColor::fromRgb(40, 40, 40);;

        if (m_hover) {
            color = QColor::fromRgb(60, 60, 60);
        }

        if (m_pressed) {
            color = QColor::fromRgb(20, 20, 20);
        }

        painter.fillPath(path, color);
    } else {
        iconName.append(PLUGIN_MIN_ICON_NAME);
    }

    pixmap = loadSvg(iconName, QSize(iconSize, iconSize));

    painter.drawPixmap(rect().center() - pixmap.rect().center() / devicePixelRatioF(), pixmap);
}

const QPixmap OnboardItem::loadSvg(const QString &fileName, const QSize &size) const
{
    const auto ratio = devicePixelRatioF();

    QPixmap pixmap;
    pixmap = QIcon::fromTheme(fileName).pixmap(size * ratio);
    pixmap.setDevicePixelRatio(ratio);

    return pixmap;
}

void OnboardItem::mousePressEvent(QMouseEvent *event)
{
    m_pressed = true;
    update();

    QWidget::mousePressEvent(event);
}

void OnboardItem::mouseReleaseEvent(QMouseEvent *event)
{
    m_pressed = false;
    m_hover = false;
    update();

    QWidget::mouseReleaseEvent(event);
}

void OnboardItem::mouseMoveEvent(QMouseEvent *event)
{
    m_hover = true;

    QWidget::mouseMoveEvent(event);
}

void OnboardItem::leaveEvent(QEvent *event)
{
    if (!rect().contains(mapFromGlobal(QCursor::pos()))) {
        m_hover = false;
        m_pressed = false;
        update();
    }

    QWidget::leaveEvent(event);
}

void OnboardItem::resizeEvent(QResizeEvent *event)
{
    const Dock::Position position = qApp->property(PROP_POSITION).value<Dock::Position>();
    // 保持横纵比
    if (position == Dock::Bottom || position == Dock::Top) {
        setMaximumWidth(height());
        setMaximumHeight(PLUGIN_BACKGROUND_MAX_SIZE);
    } else {
        setMaximumHeight(width());
        setMaximumWidth(PLUGIN_BACKGROUND_MAX_SIZE);
    }

    QWidget::resizeEvent(event);
}
