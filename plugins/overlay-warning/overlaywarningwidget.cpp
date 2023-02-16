// Copyright (C) 2011 ~ 2018 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "overlaywarningwidget.h"
#include "../../interfaces/constants.h"

#include <QSvgRenderer>
#include <QPainter>
#include <QMouseEvent>
#include <QApplication>
#include <QIcon>

OverlayWarningWidget::OverlayWarningWidget(QWidget *parent)
    : QWidget(parent)
{
}

QSize OverlayWarningWidget::sizeHint() const
{
    return QSize(26, 26);
}

void OverlayWarningWidget::paintEvent(QPaintEvent *e)
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
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(pixmap.rect());
    painter.drawPixmap(rf.center() - rfp.center() / devicePixelRatioF(), pixmap);
}

const QPixmap OverlayWarningWidget::loadSvg(const QString &fileName, const QSize &size) const
{
    const auto ratio = devicePixelRatioF();

    QSize pixmapSize = QCoreApplication::testAttribute(Qt::AA_UseHighDpiPixmaps) ? size : (size * ratio);
    QPixmap pixmap = QIcon::fromTheme(fileName).pixmap(pixmapSize);
    pixmap.setDevicePixelRatio(ratio);

    return pixmap;
}
