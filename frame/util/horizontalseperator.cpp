// SPDX-FileCopyrightText: 2011 - 2022 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "horizontalseperator.h"

#include <DApplicationHelper>

#include <QPainter>

/**
 * @brief HorizontalSeperator::HorizontalSeperator 分割线控件,高度值初始化为2个像素
 * @param parent
 */
HorizontalSeperator::HorizontalSeperator(QWidget *parent)
    : QWidget(parent)
{
    setFixedHeight(2);
    setSizePolicy(QSizePolicy::Expanding, QSizePolicy::Fixed);
}

QSize HorizontalSeperator::sizeHint() const
{
    return QSize(QWidget::sizeHint().width(), 2);
}

void HorizontalSeperator::paintEvent(QPaintEvent *e)
{
    Q_UNUSED(e)

    QPainter painter(this);
    QColor c = palette().color(QPalette::BrightText);
    c.setAlpha(int(0.1 * 255));

    painter.fillRect(rect(), c);
}
