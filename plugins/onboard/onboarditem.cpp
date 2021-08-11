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

#include <QPainter>
#include <QPainterPath>
#include <QMouseEvent>
#include <QApplication>
#include <QIcon>

#include <DStyle>
#include <DGuiApplicationHelper>

DWIDGET_USE_NAMESPACE;

OnboardItem::OnboardItem(QWidget *parent)
    : QWidget(parent)
    , m_hover(false)
    , m_pressed(false)
{
    setMouseTracking(true);
    setMinimumSize(PLUGIN_BACKGROUND_MIN_SIZE, PLUGIN_BACKGROUND_MIN_SIZE);

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, [ = ] {
        update();
    });
    m_icon = QIcon::fromTheme(":/icons/icon/deepin-virtualkeyboard.svg");
}

void OnboardItem::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    QPixmap pixmap;
    QString iconName = "deepin-virtualkeyboard";
    int iconSize = PLUGIN_ICON_MAX_SIZE;

    QPainter painter(this);
    if (std::min(width(), height()) > PLUGIN_BACKGROUND_MIN_SIZE) {

        QColor color;
        if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
            color = Qt::black;
            painter.setOpacity(0.5);

            if (m_hover) {
                painter.setOpacity(0.6);
            }

            if (m_pressed) {
                painter.setOpacity(0.3);
            }
        } else {
            color = Qt::white;
            painter.setOpacity(0.1);

            if (m_hover) {
                painter.setOpacity(0.2);
            }

            if (m_pressed) {
                painter.setOpacity(0.05);
            }
        }

        painter.setRenderHint(QPainter::Antialiasing, true);

        DStyleHelper dstyle(style());
        const int radius = dstyle.pixelMetric(DStyle::PM_FrameRadius);

        QPainterPath path;

        int minSize = std::min(width(), height());
        QRect rc(0, 0, minSize, minSize);
        rc.moveTo(rect().center() - rc.center());

        path.addRoundedRect(rc, radius, radius);
        painter.fillPath(path, color);
    } else if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
        iconName.append(PLUGIN_MIN_ICON_NAME);
    }

    pixmap = loadSvg(iconName, QSize(iconSize, iconSize));

    painter.setOpacity(1);
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(pixmap.rect());
    painter.drawPixmap(rf.center() - rfp.center() / devicePixelRatioF(), pixmap);
}

const QPixmap OnboardItem::loadSvg(const QString &fileName, const QSize &size) const
{
    const auto ratio = devicePixelRatioF();

    QPixmap pixmap;
    pixmap = QIcon::fromTheme(fileName, m_icon).pixmap(size * ratio);
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
    m_hover = false;
    m_pressed = false;
    update();

    QWidget::leaveEvent(event);
}

void OnboardItem::resizeEvent(QResizeEvent *event)
{
    QWidget::resizeEvent(event);
}
