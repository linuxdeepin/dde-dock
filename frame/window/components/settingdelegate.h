/*
 * Copyright (C) 2022 ~ 2022 Deepin Technology Co., Ltd.
 *
 * Author:     donghualin <donghualin@uniontech.com>
 *
 * Maintainer:  donghualin <donghualin@uniontech.com>
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
#ifndef SETTINGDELEGATE_H
#define SETTINGDELEGATE_H

#include <DStyledItemDelegate>

DWIDGET_USE_NAMESPACE

static const int itemCheckRole = Dtk::UserRole + 1;
static const int itemDataRole = Dtk::UserRole + 2;

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
