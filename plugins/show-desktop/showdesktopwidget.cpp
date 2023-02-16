// Copyright (C) 2019 ~ 2019 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "showdesktopwidget.h"
#include "showdesktopplugin.h"
#include "imageutil.h"

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
        icon = ImageUtil::loadSvg("deepin-toggle-desktop", QString(), int(size().width() * 0.8), ratio);
    } else {
        icon = ImageUtil::loadSvg("deepin-toggle-desktop", QString(), int(size().width() * 0.7), ratio);
    }

    icon.setDevicePixelRatio(ratio);

    QPainter painter(this);
    const QRectF &rf = QRectF(rect());
    const QRectF &rfp = QRectF(icon.rect());
    painter.drawPixmap(rf.center() - rfp.center() / ratio, icon);
}
