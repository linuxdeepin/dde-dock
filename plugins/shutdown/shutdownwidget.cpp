// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "shutdownwidget.h"

#include <QSvgRenderer>
#include <QPainter>
#include <QPainterPath>
#include <QMouseEvent>
#include <QApplication>
#include <QIcon>

#include <DStyle>
#include <DGuiApplicationHelper>

DWIDGET_USE_NAMESPACE;

ShutdownWidget::ShutdownWidget(QWidget *parent)
    : QWidget(parent)
    , m_hover(false)
    , m_pressed(false)
{
    setMouseTracking(true);
    setMinimumSize(PLUGIN_BACKGROUND_MIN_SIZE, PLUGIN_BACKGROUND_MIN_SIZE);

    connect(DGuiApplicationHelper::instance(), &DGuiApplicationHelper::themeTypeChanged, this, [ = ] {
        update();
    });
}

QPixmap ShutdownWidget::loadPixmap() const
{
    const QString iconName = "system-shutdown";
    const auto ratio = devicePixelRatioF();
    return QIcon::fromTheme(iconName).pixmap(PLUGIN_ICON_MAX_SIZE * ratio);
}

void ShutdownWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    QPainter painter(this);

    if (rect().height() > PLUGIN_BACKGROUND_MIN_SIZE) {

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
    }

    painter.setOpacity(1);

    QPixmap pixmap = loadPixmap();

    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(pixmap.rect());
    painter.drawPixmap(rf.center() - rfp.center() / pixmap.devicePixelRatioF(), pixmap);
}

void ShutdownWidget::mousePressEvent(QMouseEvent *event)
{
    m_pressed = true;
    update();

    QWidget::mousePressEvent(event);
}

void ShutdownWidget::mouseReleaseEvent(QMouseEvent *event)
{
    m_pressed = false;
    m_hover = false;
    update();

    QWidget::mouseReleaseEvent(event);
}

void ShutdownWidget::mouseMoveEvent(QMouseEvent *event)
{
    m_hover = true;
    QWidget::mouseMoveEvent(event);
}

void ShutdownWidget::leaveEvent(QEvent *event)
{
    m_hover = false;
    m_pressed = false;
    update();

    QWidget::leaveEvent(event);
}
