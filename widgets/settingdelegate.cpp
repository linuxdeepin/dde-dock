// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#include "settingdelegate.h"

#include <DListView>
#include <QMouseEvent>
#include <QPainter>
#include <DGuiApplicationHelper>
#include <QPainterPath>

DWIDGET_USE_NAMESPACE

SettingDelegate::SettingDelegate(QAbstractItemView *parent)
    : DStyledItemDelegate(parent)
{
    parent->installEventFilter(this);
}

SettingDelegate::~SettingDelegate()
{
}

void SettingDelegate::paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const
{
    painter->save();

    QRect indexRect = option.rect;
    // 绘制背景色
    bool isOver = option.state & QStyle::State_MouseOver;
    bool isDefault = index.data(itemCheckRole).toBool();
    if (isDefault) {
        QPainterPath path, path1;
        path.addRoundedRect(indexRect, 8, 8);

        DPalette palette = DGuiApplicationHelper::instance()->applicationPalette();
        painter->fillPath(path, palette.color(QPalette::ColorRole::Highlight));
    } else {
        QPainterPath path;
        path.addRoundedRect(indexRect, 8, 8);
        painter->fillPath(path, isOver ? QColor(0, 0, 0, 100) : QColor(0, 0, 0, 64));
    }
    // 绘制图标
    QRect rectIcon = indexRect;
    rectIcon.setX(20);
    QIcon icon = index.data(Qt::DecorationRole).value<QIcon>();
    QPixmap pixmap(icon.pixmap(16, 16));
    rectIcon.setY(indexRect.y() + (rectIcon.height() - pixmap.height()) / 2);
    rectIcon.setWidth(pixmap.width());
    rectIcon.setHeight(pixmap.height());
    painter->drawPixmap(rectIcon, pixmap);
#define RIGHTSPACE 11
#define SELECTICONSIZE 10
    // 绘制文本
    QRect rectText;
    rectText.setX(rectIcon.left() + rectIcon.width() + 8);
    rectText.setWidth(indexRect.width() - rectText.x() - RIGHTSPACE - SELECTICONSIZE - 5);
    QPen pen(isDefault ? QColor(255, 255, 255) : QColor(0, 0, 0));
    pen.setWidth(2);
    painter->setPen(pen);
    QFont ft(DFontSizeManager::instance()->t6());
    QFontMetrics ftm(ft);
    QString text = QFontMetrics(ft).elidedText(index.data(Qt::DisplayRole).toString(), Qt::TextElideMode::ElideRight,
                                               rectText.width());
    painter->setFont(ft);
    rectText.setY(indexRect.y() + (indexRect.height() - QFontMetrics(ft).height()) / 2);
    rectText.setHeight(QFontMetrics(ft).height());
    painter->drawText(rectText, text);
    // 如果当前是默认的输出设备，则绘制右侧的对钩
    if (isDefault) {
        QPointF points[3] = {
            QPointF(indexRect.width() - RIGHTSPACE - SELECTICONSIZE, indexRect.center().y()),
            QPointF(indexRect.width() - RIGHTSPACE - SELECTICONSIZE / 2, rectIcon.bottom() + 2),
            QPointF(indexRect.width() - RIGHTSPACE, rectIcon.top() - 2)
        };
        painter->drawPolyline(points, 3);
    }

    painter->restore();
}

bool SettingDelegate::editorEvent(QEvent *event, QAbstractItemModel *model, const QStyleOptionViewItem &option, const QModelIndex &index)
{
    if (event->type() == QEvent::MouseButtonRelease) {
        QRect rctIndex = option.rect;
        rctIndex.setHeight(rctIndex.height() - spacing());
        QMouseEvent *mouseEvent = static_cast<QMouseEvent *>(event);
        if (rctIndex.contains(mouseEvent->pos()))
            Q_EMIT selectIndexChanged(index);
    }

    return DStyledItemDelegate::editorEvent(event, model, option, index);
}
