// Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
// SPDX-FileCopyrightText: 2018 - 2023 UnionTech Software Technology Co., Ltd.
//
// SPDX-License-Identifier: LGPL-3.0-or-later

#ifndef SETTINGDELEGATE_H
#define SETTINGDELEGATE_H

#include <DStyledItemDelegate>

DWIDGET_USE_NAMESPACE

static const int itemCheckRole = Dtk::UserRole + 1;
static const int itemDataRole = Dtk::UserRole + 2;
static const int itemFlagRole = Dtk::UserRole + 3;

class SettingDelegate : public DStyledItemDelegate
{
    Q_OBJECT

public:
    explicit SettingDelegate(QAbstractItemView *parent = nullptr);
    ~SettingDelegate() override;

Q_SIGNALS:
    void selectIndexChanged(const QModelIndex &);

protected:
    void paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const override;
    bool editorEvent(QEvent *event, QAbstractItemModel *model, const QStyleOptionViewItem &option, const QModelIndex &index) override;
};

#endif // SETTINGDELEGATE_H
