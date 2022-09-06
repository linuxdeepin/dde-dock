// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
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
    m_icon = QIcon::fromTheme(":/icons/resources/icons/system-shutdown.svg");
}

void ShutdownWidget::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e);

    QPixmap pixmap;
    QString iconName = "system-shutdown";
    int iconSize = PLUGIN_ICON_MAX_SIZE;

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
    } else if (DGuiApplicationHelper::instance()->themeType() == DGuiApplicationHelper::LightType) {
        // 最小尺寸时，不画背景，采用深色图标
        iconName.append(PLUGIN_MIN_ICON_NAME);
    }

    painter.setOpacity(1);

    pixmap = loadSvg(iconName, QSize(iconSize, iconSize));

    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(pixmap.rect());
    painter.drawPixmap(rf.center() - rfp.center() / pixmap.devicePixelRatioF(), pixmap);
}

const QPixmap ShutdownWidget::loadSvg(const QString &fileName, const QSize &size) const
{
    const auto ratio = devicePixelRatioF();

    QPixmap pixmap;
    pixmap = QIcon::fromTheme(fileName, m_icon).pixmap(size * ratio);
    pixmap.setDevicePixelRatio(ratio);

    return pixmap;
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
